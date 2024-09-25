// cmd/function/function_watch.go
package function

import (
	"fmt"
	"time"

	"github.com/brevdev/nvcf/api"
	"github.com/brevdev/nvcf/config"
	"github.com/brevdev/nvcf/output"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
	"github.com/tmc/nvcf-go"
)

func functionWatchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "watch",
		Short: "Watch functions status in real-time",
		Long:  `Display a real-time view of NVCF functions.`,
		Args:  cobra.MaximumNArgs(1),
		RunE:  runFunctionWatch,
	}
	cmd.Flags().StringSlice("status", []string{"ACTIVE", "DEPLOYING", "ERROR", "INACTIVE", "DELETED"}, "Filter by status (ACTIVE, DEPLOYING, ERROR, INACTIVE, DELETED). Defaults to all.")
	cmd.Flags().StringSlice("visibility", []string{"private"}, "Filter by visibility (authorized, private, public). Defaults to private.")
	return cmd
}

func runFunctionWatch(cmd *cobra.Command, args []string) error {
	client := api.NewClient(config.GetAPIKey())
	app := tview.NewApplication()
	statusParams, err := parseStatusFlags(cmd)
	if err != nil {
		return err
	}
	visibilityParams, err := parseVisibilityFlags(cmd)
	if err != nil {
		return err
	}

	mainTable := tview.NewTable().
		SetBorders(true).
		SetSelectable(true, false)
	mainTable.SetSelectedStyle(
		tcell.StyleDefault.
			Background(tcell.ColorGreen).
			Foreground(tcell.ColorWhite),
	)

	detailsView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetScrollable(true)

	pages := tview.NewPages()

	updateTable := func() {
		var functions *nvcf.ListFunctionsResponse
		var err error
		if len(args) == 1 {
			functionID := args[0]
			versions, err := client.Functions.Versions.List(cmd.Context(), functionID)
			if err != nil {
				_ = output.Error(cmd, "Error getting function", err)
				return
			}
			functions = &nvcf.ListFunctionsResponse{
				Functions: versions.Functions,
			}
		} else {
			functions, err = client.Functions.List(cmd.Context(), nvcf.FunctionListParams{
				Visibility: nvcf.F(visibilityParams),
			})
			if err != nil {
				_ = output.Error(cmd, "Error listing functions", err)
				return
			}
		}
		mainTable.Clear()
		mainTable.SetCell(0, 0, tview.NewTableCell("Name").SetTextColor(tcell.ColorYellow).SetSelectable(false))
		mainTable.SetCell(0, 1, tview.NewTableCell("ID").SetTextColor(tcell.ColorYellow).SetSelectable(false))
		mainTable.SetCell(0, 2, tview.NewTableCell("Version ID").SetTextColor(tcell.ColorYellow).SetSelectable(false))
		mainTable.SetCell(0, 3, tview.NewTableCell("Status").SetTextColor(tcell.ColorYellow).SetSelectable(false))
		row := 1
		for _, fn := range functions.Functions {
			if containsStatus(statusParams, fn.Status) {
				mainTable.SetCell(row, 0, tview.NewTableCell(fn.Name))
				mainTable.SetCell(row, 1, tview.NewTableCell(fn.ID))
				mainTable.SetCell(row, 2, tview.NewTableCell(fn.VersionID))
				mainTable.SetCell(row, 3, tview.NewTableCell(string(fn.Status)).SetTextColor(getStatusColor(fn.Status)))
				row++
			}
		}
	}

	updateTable() // Initial update

	mainTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			row, _ := mainTable.GetSelection()
			if row > 0 {
				functionID := mainTable.GetCell(row, 1).Text
				versionID := mainTable.GetCell(row, 2).Text
				showFunctionDetails(client, cmd, functionID, versionID, detailsView, pages)
			}
		}
		return event
	})

	pages.AddPage("main", mainTable, true, true)
	pages.AddPage("details", detailsView, true, false)

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		for {
			select {
			case <-ticker.C:
				app.QueueUpdateDraw(func() {
					updateTable()
				})
			case <-cmd.Context().Done():
				return
			}
		}
	}()

	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		return output.Error(cmd, "error running UI", err)
	}
	return nil
}

func showFunctionDetails(client *api.Client, cmd *cobra.Command, functionID, versionID string, detailsView *tview.TextView, pages *tview.Pages) {
	function, err := client.Functions.Versions.Get(cmd.Context(), functionID, versionID, nvcf.FunctionVersionGetParams{
		IncludeSecrets: nvcf.Bool(false),
	})
	if err != nil {
		detailsView.SetText(fmt.Sprintf("Error getting function details: %v", err))
		pages.SwitchToPage("details")
		return
	}

	fn := function.Function
	details := fmt.Sprintf(`[yellow]Function Details for %s (Version: %s)[white]

[green]General Information:[white]
Name: %s
Version ID: %s
Status: %s
Function Type: %s
Created At: %s

[green]Configuration:[white]
Inference URL: %s
Inference Port: %d
Container Image: %s
Container Args: %s
API Body Format: %s

[green]Health Check:[white]
Protocol: %s
Port: %d
Timeout: %s
Expected Status Code: %d
URI: %s

[green]Active Instances:[white]
`,
		fn.Name, fn.VersionID,
		fn.Name, fn.VersionID, fn.Status, fn.FunctionType, fn.CreatedAt.Format(time.RFC3339),
		fn.InferenceURL, fn.InferencePort, fn.ContainerImage, fn.ContainerArgs, fn.APIBodyFormat,
		fn.Health.Protocol, fn.Health.Port, fn.Health.Timeout, fn.Health.ExpectedStatusCode, fn.Health.Uri)

	for _, instance := range fn.ActiveInstances {
		details += fmt.Sprintf("- ID: %s, Status: %s, GPU: %s, Instance Type: %s\n",
			instance.InstanceID, instance.InstanceStatus, instance.GPU, instance.InstanceType)
	}

	details += "\n[yellow]Press Esc to return to the main view[white]"

	detailsView.SetText(details)
	detailsView.ScrollToBeginning()

	detailsView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			pages.SwitchToPage("main")
		}
		return event
	})

	pages.SwitchToPage("details")
}

func getStatusColor(status nvcf.ListFunctionsResponseFunctionsStatus) tcell.Color {
	switch status {
	case nvcf.ListFunctionsResponseFunctionsStatusActive:
		return tcell.ColorGreen
	case nvcf.ListFunctionsResponseFunctionsStatusDeploying:
		return tcell.ColorYellow
	case nvcf.ListFunctionsResponseFunctionsStatusError:
		return tcell.ColorRed
	default:
		return tcell.ColorWhite
	}
}

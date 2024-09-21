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
		Short: "Watch functions in real-time",
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

	table := tview.NewTable().
		SetBorders(true).
		SetSelectable(true, false)

	table.SetSelectedStyle(
		tcell.StyleDefault.
			Background(tcell.ColorGreen).
			Foreground(tcell.ColorWhite),
	)

	updateTable := func() {
		var functions *nvcf.ListFunctionsResponse
		var err error

		if len(args) == 1 {
			functionID := args[0]
			versions, err := client.Functions.Versions.List(cmd.Context(), functionID)
			if err != nil {
				output.Error(cmd, "Error getting function", err)
				return
			}
			for _, v := range versions.Functions {
				functions = &nvcf.ListFunctionsResponse{
					Functions: []nvcf.ListFunctionsResponseFunction{v},
				}
			}
		} else {
			functions, err = client.Functions.List(cmd.Context(), nvcf.FunctionListParams{
				Visibility: nvcf.F(visibilityParams),
			})
			if err != nil {
				output.Error(cmd, "Error listing functions", err)
				return
			}
		}

		table.Clear()
		table.SetCell(0, 0, tview.NewTableCell("Name").SetTextColor(tcell.ColorYellow).SetSelectable(false))
		table.SetCell(0, 1, tview.NewTableCell("Version ID").SetTextColor(tcell.ColorYellow).SetSelectable(false))
		table.SetCell(0, 2, tview.NewTableCell("Status").SetTextColor(tcell.ColorYellow).SetSelectable(false))

		row := 1
		for _, fn := range functions.Functions {
			if containsStatus(statusParams, fn.Status) {
				table.SetCell(row, 0, tview.NewTableCell(fn.Name))
				table.SetCell(row, 1, tview.NewTableCell(fn.VersionID))
				table.SetCell(row, 2, tview.NewTableCell(string(fn.Status)).SetTextColor(getStatusColor(fn.Status)))
				row++
			}
		}
	}

	updateTable() // Initial update

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

	if err := app.SetRoot(table, true).EnableMouse(true).Run(); err != nil {
		return fmt.Errorf("error running UI: %w", err)
	}

	return nil
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

// this logic is based on the ui call
package function

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/brevdev/nvcf/api"
	"github.com/brevdev/nvcf/collections"
	"github.com/brevdev/nvcf/config"
	"github.com/brevdev/nvcf/output"
	"github.com/spf13/cobra"
)

func functionLogsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "logs <function-id>",
		Args:    cobra.ExactArgs(1),
		Short:   "Get logs for a function",
		Long:    "Get logs for a function. If a version-id is not provided and there are multiple versions associated with a function, we will look for all versions and prompt for a version-id.",
		Example: "nvcf function logs <function-id> --version-id <version-id> --start \"2023-04-01 00:00:00\" --end \"2023-04-02 00:00:00\"",
		RunE:    runFunctionLogs,
	}

	const timeFormat = "2006-01-02 15:04:05"

	cmd.Flags().String("version-id", "", "The ID of the version")
	cmd.Flags().String("start", time.Now().Add(-24*time.Hour).Format(timeFormat), "Start time for logs (format: YYYY-MM-DD HH:MM:SS)")
	cmd.Flags().String("end", time.Now().Format(timeFormat), "End time for logs (format: YYYY-MM-DD HH:MM:SS)")

	return cmd
}

func runFunctionLogs(cmd *cobra.Command, args []string) error {
	client := api.NewClient(config.GetAPIKey())
	functionID := args[0]
	versionID, _ := cmd.Flags().GetString("version-id")

	if versionID == "" {
		versions, err := client.Functions.Versions.List(cmd.Context(), functionID)
		if err != nil {
			return output.Error(cmd, "Error listing function versions", err)
		}

		if len(versions.Functions) == 1 {
			versionID = versions.Functions[0].VersionID
		} else {
			output.Info(cmd, "Multiple versions found. Please specify a version-id")
			for _, version := range versions.Functions {
				output.Info(cmd, fmt.Sprintf("Version ID: %s || Status: %s", version.VersionID, version.Status))
			}
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Enter version-id: ")
			versionID, _ = reader.ReadString('\n')
			versionID = strings.TrimSpace(versionID)
		}
	}

	// Parse start and end times
	startTime, _ := cmd.Flags().GetString("start")
	endTime, _ := cmd.Flags().GetString("end")

	// Call getDeploymentLogs with the parsed arguments
	logs, err := getDeploymentLogs(cmd.Context(), provider.GetDeploymentLogArgs{
		ID:        fmt.Sprintf("%s:%s", functionID, versionID),
		StartTime: startTime,
		EndTime:   endTime,
	})
	if err != nil {
		return output.Error(cmd, "Error getting function logs", err)
	}

	// Display logs
	output.Info(cmd, "Function Logs:")
	for _, log := range logs {
		output.Info(cmd, log.Message)
	}

	return nil
}

// note: this is not a spec method - im building this based on the ui call
func getDeploymentLogs(ctx context.Context, client *api.Client, functionID, versionID string, startTime, endTime time.Time) ([]DeploymentLog, error) {
	url := buildLogsURL(functionID, versionID)
	payload := buildLogsPayload(startTime, endTime)
	var logsResponse NVCFLogsResponse
	err := client.Post(ctx, url, payload, &logsResponse)
	if err != nil {
		return nil, err
	}

	logs, err := formatLogs(logsResponse.Data)
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func buildLogsURL(funcID string, versionID string) string {
	return fmt.Sprintf("/v2/orgs/%s/nvcf/logs/functions/%s/versions/%s", n.OrgID, funcID, versionID)
}

func buildLogsPayload(startTime time.Time, endTime time.Time) LogsPayload {
	const timeFormat = "2006-01-02 15:04:05"
	return LogsPayload{
		Parameters: []LogParameter{
			{Name: "start", Value: startTime.Format(timeFormat)},
			{Name: "end", Value: endTime.Format(timeFormat)},
			{Name: "sort", Value: "desc"},
		},
	}
}

func formatLogs(logs []NVCFLog) ([]DeploymentLog, error) {
	deploymentLogs, err := collections.MapE(logs, func(log NVCFLog) (DeploymentLog, error) {
		deploymentLog, err := collections.TryCopyToNew[NVCFLog, DeploymentLog](log)
		if err != nil {
			return DeploymentLog{}, err
		}
		return deploymentLog, nil
	})
	if err != nil {
		return nil, err
	}
	return deploymentLogs, nil
}

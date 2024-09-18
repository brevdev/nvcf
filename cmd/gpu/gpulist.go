package gpu

import (
	"github.com/brevdev/nvcf/api"
	"github.com/brevdev/nvcf/config"
	"github.com/brevdev/nvcf/output"
	"github.com/spf13/cobra"
	"github.com/tmc/nvcf-go"
)

func gpuListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List available GPUs",
		Long:  "List available GPUs in NVCF. Note this is still in beta and lists all available GPUs ",
		RunE:  runGpuList,
	}

	cmd.Flags().String("backend", "", "Filter by backend (e.g., GFN)")
	cmd.Flags().String("gpu", "", "Filter by GPU type (e.g., L40G)")

	return cmd
}

func runGpuList(cmd *cobra.Command, args []string) error {
	client := api.NewClient(config.GetAPIKey())

	gpus, err := client.ClusterGroups.List(cmd.Context())
	if err != nil {
		return output.Error(cmd, "Error listing cluster groups", err)
	}

	backend, _ := cmd.Flags().GetString("backend")
	gpuType, _ := cmd.Flags().GetString("gpu")

	filteredGroups := filterClusterGroups(gpus.ClusterGroups, backend, gpuType)

	output.GPUs(cmd, filteredGroups)

	return nil
}

func filterClusterGroups(groups []nvcf.ClusterGroupsResponseClusterGroup, backend, gpuType string) []nvcf.ClusterGroupsResponseClusterGroup {
	var filtered []nvcf.ClusterGroupsResponseClusterGroup

	for _, group := range groups {
		if backend != "" && group.Name != backend {
			continue
		}

		var filteredGPUs []nvcf.ClusterGroupsResponseClusterGroupsGPU
		for _, gpu := range group.GPUs {
			if gpuType != "" && gpu.Name != gpuType {
				continue
			}
			filteredGPUs = append(filteredGPUs, gpu)
		}

		if len(filteredGPUs) > 0 {
			filteredGroup := group
			filteredGroup.GPUs = filteredGPUs
			filtered = append(filtered, filteredGroup)
		}
	}

	return filtered
}

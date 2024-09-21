// the logic used to filter available clusters comes from the NGC frontend for NVCF
package gpu

import (
	"context"
	"fmt"

	"github.com/brevdev/nvcf/api"
	"github.com/brevdev/nvcf/collections"
	"github.com/brevdev/nvcf/config"
	"github.com/brevdev/nvcf/output"
	"github.com/spf13/cobra"
	"github.com/tmc/nvcf-go"
)

func gpuListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List available GPUs",
		Long:  "List available GPUs in NVCF. Note this is still in beta and lists all available GPUs",
		RunE:  runGpuList,
	}

	cmd.Flags().String("backend", "", "Filter by backend (e.g., GFN)")
	cmd.Flags().String("gpu", "", "Filter by GPU type (e.g., L40G)")

	return cmd
}

func runGpuList(cmd *cobra.Command, args []string) error {
	backend, _ := cmd.Flags().GetString("backend")
	gpuType, _ := cmd.Flags().GetString("gpu")

	availableInstanceTypes, err := GetAvailableInstanceTypes(cmd.Context(), backend, gpuType)
	if err != nil {
		return output.Error(cmd, "Error getting available instance types", err)
	}

	filteredClusterGroups := filterClusterGroups(availableInstanceTypes, backend, gpuType)

	output.GPUs(cmd, filteredClusterGroups)

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

func buildNVCFOrgInformationURL(orgID string) string {
	return fmt.Sprintf("/v3/orgs/%s/nvcf", orgID)
}

func getNVCFOrgInformation(orgID string) (OrgClusterGroupsResponse, error) {
	client := api.NewClient(config.GetAPIKey())
	url := buildNVCFOrgInformationURL(orgID)
	var res OrgClusterGroupsResponse
	err := client.Get(context.Background(), url, nil, &res)
	if err != nil {
		return OrgClusterGroupsResponse{}, err
	}
	return res, nil
}

func GetNVCFClusterGroups(ctx context.Context) (*nvcf.ClusterGroupsResponse, error) {
	client := api.NewClient(config.GetAPIKey())
	availableCluster, err := client.ClusterGroups.List(ctx)
	if err != nil {
		return nil, err
	}
	return availableCluster, nil
}

type OrgClusterGroupsResponse struct {
	BillingAccountID string        `json:"billingAccountId"`
	Clusters         []Cluster     `json:"clusters"`
	RequestStatus    RequestStatus `json:"requestStatus"`
}

type Cluster struct {
	Cluster          string `json:"cluster"`
	GpuType          string `json:"gpuType"`
	InstanceType     string `json:"instanceType"`
	MaxInstances     int    `json:"maxInstances"`
	CurrentInstances int    `json:"currentInstances"`
}

type RequestStatus struct {
	StatusCode string `json:"statusCode"`
}

func GetAvailableInstanceTypes(ctx context.Context, backend, gpuType string) ([]nvcf.ClusterGroupsResponseClusterGroup, error) {
	orgInfo, err := getNVCFOrgInformation(config.GetOrgID())
	if err != nil {
		return nil, err
	}
	clusterGroups, err := GetNVCFClusterGroups(ctx)
	if err != nil {
		return nil, err
	}

	gpuConfigs := make(map[string]bool)
	for _, cluster := range orgInfo.Clusters {
		gpuConfigs[cluster.Cluster] = true
	}

	filteredGroups := collections.Filter(clusterGroups.ClusterGroups, func(cluster nvcf.ClusterGroupsResponseClusterGroup) bool {
		// Check if the cluster is owned by the organization (BYOC)
		if cluster.NcaID == orgInfo.BillingAccountID {
			return true
		}

		// Check if the cluster is shared from another account but authorized for use
		for _, authorizedNcaID := range cluster.AuthorizedNcaIDs {
			if authorizedNcaID == orgInfo.BillingAccountID {
				return true
			}
		}

		// Check if the cluster is a shared gated cluster and has GPU resources available
		if collections.ListContains(cluster.AuthorizedNcaIDs, "*") && gpuConfigs[cluster.Name] {
			return true
		}

		return false
	})

	return filteredGroups, nil
}

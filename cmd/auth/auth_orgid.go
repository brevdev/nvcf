package auth

import (
	"context"
	"fmt"

	"github.com/brevdev/nvcf/api"
	"github.com/brevdev/nvcf/config"
	"github.com/brevdev/nvcf/output"
	"github.com/spf13/cobra"
)

func authOrgIDCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "org-id",
		Short: "Display the name of the first organization",
		RunE: func(cmd *cobra.Command, args []string) error {
			orgId, err := GetOrgId()
			if err != nil {
				return output.Error(cmd, "Failed to fetch organization information", err)
			}
			fmt.Println(orgId)
			return nil
		},
	}
}

func GetOrgId() (string, error) {
	client := api.NewClient(config.GetAPIKey())
	orgsInfo := map[string]interface{}{}
	err := client.Get(context.Background(), "/v2/orgs", nil, &orgsInfo)
	if err != nil {
		return "", err
	}

	organizations, ok := orgsInfo["organizations"].([]interface{})
	if !ok || len(organizations) == 0 {
		return "", fmt.Errorf("no organizations found")
	}

	firstOrg, ok := organizations[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("failed to parse organization information")
	}

	name, ok := firstOrg["name"].(string)
	if !ok {
		return "", fmt.Errorf("organization name not found")
	}

	return name, nil
}

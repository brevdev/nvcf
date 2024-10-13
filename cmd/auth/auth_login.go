package auth

import (
	"fmt"

	"github.com/brevdev/nvcf/api"
	"github.com/brevdev/nvcf/config"
	"github.com/brevdev/nvcf/output"
	"github.com/spf13/cobra"
)

func authLoginCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "login",
		Short: "Authenticate with NVIDIA Cloud",
		RunE: func(cmd *cobra.Command, args []string) error {
			apiKey := output.Prompt("Enter your NVIDIA Cloud API key: ", true)

			err := config.SetAPIKey(apiKey)
			if err != nil {
				return output.Error(cmd, "Error saving API key", err)
			}

			// Use the API key to get the first org
			client := api.NewClient(apiKey)
			orgsInfo := map[string]interface{}{}
			err = client.Get(cmd.Context(), "/v2/orgs", nil, &orgsInfo)
			if err != nil {
				return output.Error(cmd, "Failed to fetch organization information", err)
			}

			organizations, ok := orgsInfo["organizations"].([]interface{})
			if !ok || len(organizations) == 0 {
				return output.Error(cmd, "No organizations found", nil)
			}

			firstOrg, ok := organizations[0].(map[string]interface{})
			if !ok {
				return output.Error(cmd, "Failed to parse organization information", nil)
			}

			orgID, ok := firstOrg["name"].(string)
			if !ok {
				return output.Error(cmd, "Organization ID not found", nil)
			}

			err = config.SetOrgID(orgID)
			if err != nil {
				return output.Error(cmd, "Error saving Org ID", err)
			}

			output.PrintASCIIArt(cmd)
			output.Success(cmd, fmt.Sprintf("Authentication successful. You are now authenticated with organization ID: %s", orgID))
			return nil
		},
	}
}

package auth

import (
	"encoding/json"

	"github.com/brevdev/nvcf/api"
	"github.com/brevdev/nvcf/config"
	"github.com/brevdev/nvcf/output"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var whoamiURL = "/v2/users/me"

func authWhoAmICmd() *cobra.Command {
	return &cobra.Command{
		Use:   "whoami",
		Short: "Display information about the authenticated user",
		RunE: func(cmd *cobra.Command, args []string) error {
			client := api.NewClient(config.GetAPIKey())
			whoamiInfo := map[string]any{}
			err := client.Get(cmd.Context(), whoamiURL, nil, &whoamiInfo)
			if err != nil {
				return output.Error(cmd, "Failed to fetch user information", err)
			}

			jsonMode, _ := cmd.Flags().GetBool("json")
			if jsonMode {
				err = json.NewEncoder(cmd.OutOrStdout()).Encode(whoamiInfo)
				if err != nil {
					return output.Error(cmd, "Failed to encode user information", err)
				}
				return nil
			}
			userInfo, _ := whoamiInfo["user"].(map[string]any)
			table := tablewriter.NewWriter(cmd.OutOrStdout())
			table.SetHeader([]string{"Email", "Name"})
			table.SetBorder(false)
			table.Append([]string{
				userInfo["email"].(string),
				userInfo["name"].(string),
			})
			table.Render()
			return nil
		},
	}
}

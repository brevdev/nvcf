package auth

import (
	"errors"
	"fmt"

	"github.com/brevdev/nvcf/api"
	"github.com/brevdev/nvcf/config"
	"github.com/brevdev/nvcf/output"
	"github.com/spf13/cobra"
)

func authStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Check the authentication status",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !config.IsAuthenticated() {
				return output.Error(cmd, "Not authenticated", errors.New("no API key found"))
			}

			client := api.NewClient(config.GetAPIKey())

			userInfo := map[string]interface{}{}
			err := client.Get(cmd.Context(), "/v2/users/me", nil, &userInfo)
			if err != nil {
				return output.Error(cmd, "Failed to fetch user information", err)
			}

			orgsInfo := map[string]interface{}{}
			err = client.Get(cmd.Context(), "/v2/orgs", nil, &orgsInfo)
			if err != nil {
				return output.Error(cmd, "Failed to fetch organization information", err)
			}

			user, _ := userInfo["user"].(map[string]interface{})
			email, _ := user["email"].(string)
			name, _ := user["name"].(string)
			currentOrgID := config.GetOrgID()

			output.Success(cmd, "Authenticated")
			fmt.Printf("User: %s (%s)\n", name, email)
			fmt.Printf("Current Organization ID: %s\n", currentOrgID)
			return nil
		},
	}
}

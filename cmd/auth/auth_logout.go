package auth

import (
	"github.com/brevdev/nvcf/config"
	"github.com/brevdev/nvcf/output"
	"github.com/spf13/cobra"
)

func authLogoutCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Logout from NVIDIA Cloud",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !config.IsAuthenticated() {
				output.Info(cmd, "You are currently not logged in")
				return nil
			}
			err := config.ClearAPIKey()
			if err != nil {
				return output.Error(cmd, "Failed to clear API key", err)
			}
			err = config.ClearOrgID()
			if err != nil {

				return output.Error(cmd, "Failed to clear Org ID", err)
			}
			output.Success(cmd, "Logged out successfully")
			return nil
		},
	}
}

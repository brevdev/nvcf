package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/brevdev/nvcf/api"
	"github.com/brevdev/nvcf/config"
	"github.com/brevdev/nvcf/output"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func AuthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Manage authentication for the CLI",
		Long:  `Authenticate with NVIDIA Cloud and configure the CLI to use your API key.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			config.Init()
			if cmd.Parent().Name() != "auth" && !config.IsAuthenticated() {
				fmt.Println("You are not authenticated. Please run 'nvcf auth login' first.")
				os.Exit(1)
			}
		},
	}

	cmd.AddCommand(authLoginCmd())
	cmd.AddCommand(authLogoutCmd())
	cmd.AddCommand(authStatusCmd())
	cmd.AddCommand(authConfigureDockerCmd())

	cmd.AddCommand(authWhoAmICmd())
	cmd.AddCommand(authOrgsCmd())
	cmd.AddCommand(authOrgIDCmd())

	return cmd
}

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

func authConfigureDockerCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "configure-docker",
		Short: "Configure Docker to use NGC API key for nvcr.io",
		RunE: func(cmd *cobra.Command, args []string) error {
			apiKey := config.GetAPIKey()
			if apiKey == "" {
				return output.Error(cmd, "NGC API key not found. Please run 'nvcf auth login' first.", nil)
			}
			// Check if Docker is installed
			_, err := exec.LookPath("docker")
			if err != nil {
				return output.Error(cmd, "Docker is not installed or not in the system PATH", err)
			}
			// TODO: check for existing nvcr.io config?
			dockerCmd := exec.Command("docker", "login", "nvcr.io", "-u", "$oauthtoken", "--password-stdin")
			dockerCmd.Stdin = strings.NewReader(apiKey)
			out, err := dockerCmd.CombinedOutput()
			if err != nil {
				cmd.Println(string(out))
				return output.Error(cmd, "Failed to configure Docker", err)
			}
			output.Success(cmd, "Docker configured successfully for nvcr.io")
			cmd.Println(string(out))
			return nil
		},
	}
}

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

func authOrgsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "orgs",
		Short: "Display organization and team information for the authenticated user",
		RunE: func(cmd *cobra.Command, args []string) error {
			client := api.NewClient(config.GetAPIKey())
			userInfo := map[string]interface{}{}
			err := client.Get(cmd.Context(), "/v2/users/me", nil, &userInfo)
			if err != nil {
				return output.Error(cmd, "Failed to fetch user information", err)
			}
			jsonMode, _ := cmd.Flags().GetBool("json")
			if jsonMode {
				err = json.NewEncoder(cmd.OutOrStdout()).Encode(userInfo)
				if err != nil {
					return output.Error(cmd, "Failed to encode user information", err)
				}
				return nil
			}
			userRoles, ok := userInfo["userRoles"].([]interface{})
			if !ok {
				return output.Error(cmd, "Failed to parse user roles information", nil)
			}
			type OrgTeamInfo struct {
				OrgName        string
				OrgDisplayName string
				OrgType        string
				TeamName       string
				OrgRoles       string
			}
			var orgTeamList []OrgTeamInfo
			for _, role := range userRoles {
				roleMap, ok := role.(map[string]interface{})
				if !ok {
					continue
				}
				org, ok := roleMap["org"].(map[string]interface{})
				if !ok {
					continue
				}
				team, _ := roleMap["team"].(map[string]interface{})
				orgName, _ := org["name"].(string)
				orgDisplayName, _ := org["displayName"].(string)
				orgType, _ := org["type"].(string)
				teamName, _ := team["name"].(string)
				orgRoles, _ := roleMap["orgRoles"].([]interface{})
				orgRolesStr := strings.Join(convertToStringSlice(orgRoles), ",")
				orgTeamList = append(orgTeamList, OrgTeamInfo{
					OrgName:        orgName,
					OrgDisplayName: orgDisplayName,
					OrgType:        orgType,
					TeamName:       teamName,
					OrgRoles:       orgRolesStr,
				})
			}
			// Sort the list by org name, then team name
			sort.Slice(orgTeamList, func(i, j int) bool {
				if orgTeamList[i].OrgName == orgTeamList[j].OrgName {
					return orgTeamList[i].TeamName < orgTeamList[j].TeamName
				}
				return orgTeamList[i].OrgName < orgTeamList[j].OrgName
			})

			wideMode, _ := cmd.Flags().GetBool("wide")
			table := tablewriter.NewWriter(cmd.OutOrStdout())
			if wideMode {
				table.SetHeader([]string{"Org Name", "Team Name", "Org Roles"})
			} else {
				table.SetHeader([]string{"Org Name", "Org Display Name", "Org Type", "Team Name"})
			}
			table.SetBorder(false)
			for _, info := range orgTeamList {
				if wideMode {
					table.Append([]string{info.OrgName, info.TeamName, info.OrgRoles})
				} else {
					table.Append([]string{info.OrgName, info.OrgDisplayName, info.OrgType, info.TeamName})
				}
			}
			table.Render()
			return nil
		},
	}

	cmd.Flags().BoolP("wide", "o", false, "Display wide output including org roles")
	return cmd
}

func authOrgIDCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "org-id",
		Short: "Display the name of the first organization",
		RunE: func(cmd *cobra.Command, args []string) error {
			client := api.NewClient(config.GetAPIKey())
			orgsInfo := map[string]interface{}{}
			err := client.Get(cmd.Context(), "/v2/orgs", nil, &orgsInfo)
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

			name, ok := firstOrg["name"].(string)
			if !ok {
				return output.Error(cmd, "Organization name not found", nil)
			}

			fmt.Println(name)
			return nil
		},
	}
}

func convertToStringSlice(slice []interface{}) []string {
	result := make([]string, len(slice))
	for i, v := range slice {
		result[i] = fmt.Sprint(v)
	}
	return result
}

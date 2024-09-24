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
		Run: func(cmd *cobra.Command, args []string) {
			apiKey := output.Prompt("Enter your NVIDIA Cloud API key: ", true)

			// Save the API key first
			err := config.SetAPIKey(apiKey)
			if err != nil {
				output.Error(cmd, "Error saving API key", err)
				return
			}

			// Use the API key to get the first org
			client := api.NewClient(apiKey)
			orgsInfo := map[string]interface{}{}
			err = client.Get(cmd.Context(), "/v2/orgs", nil, &orgsInfo)
			if err != nil {
				output.Error(cmd, "Failed to fetch organization information", err)
				return
			}

			organizations, ok := orgsInfo["organizations"].([]interface{})
			if !ok || len(organizations) == 0 {
				output.Error(cmd, "No organizations found", nil)
				return
			}

			firstOrg, ok := organizations[0].(map[string]interface{})
			if !ok {
				output.Error(cmd, "Failed to parse organization information", nil)
				return
			}

			orgID, ok := firstOrg["name"].(string)
			if !ok {
				output.Error(cmd, "Organization ID not found", nil)
				return
			}

			// Save the org ID
			err = config.SetOrgID(orgID)
			if err != nil {
				output.Error(cmd, "Error saving Org ID", err)
				return
			}

			output.PrintASCIIArt(cmd)
			output.Success(cmd, fmt.Sprintf("Authentication successful. You are now authenticated with organization ID: %s", orgID))
		},
	}
}

func authConfigureDockerCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "configure-docker",
		Short: "Configure Docker to use NGC API key for nvcr.io",
		Run: func(cmd *cobra.Command, args []string) {
			apiKey := config.GetAPIKey()
			if apiKey == "" {
				output.Error(cmd, "NGC API key not found. Please run 'nvcf auth login' first.", nil)
				return
			}
			// Check if Docker is installed
			_, err := exec.LookPath("docker")
			if err != nil {
				output.Error(cmd, "Docker is not installed or not in the system PATH", err)
				return
			}
			// TODO: check for existing nvcr.io config?
			dockerCmd := exec.Command("docker", "login", "nvcr.io", "-u", "$oauthtoken", "--password-stdin")
			dockerCmd.Stdin = strings.NewReader(apiKey)
			out, err := dockerCmd.CombinedOutput()
			if err != nil {
				output.Error(cmd, "Failed to configure Docker", err)
				cmd.Println(string(out))
				return
			}
			output.Success(cmd, "Docker configured successfully for nvcr.io")
			cmd.Println(string(out))
		},
	}
}

func authStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Check the authentication status",
		Run: func(cmd *cobra.Command, args []string) {
			if !config.IsAuthenticated() {
				output.Error(cmd, "Not authenticated", errors.New("No API key found"))
				return
			}

			client := api.NewClient(config.GetAPIKey())

			// Fetch user information
			userInfo := map[string]interface{}{}
			err := client.Get(cmd.Context(), "/v2/users/me", nil, &userInfo)
			if err != nil {
				output.Error(cmd, "Failed to fetch user information", err)
				return
			}

			// Fetch organization information
			orgsInfo := map[string]interface{}{}
			err = client.Get(cmd.Context(), "/v2/orgs", nil, &orgsInfo)
			if err != nil {
				output.Error(cmd, "Failed to fetch organization information", err)
				return
			}

			// Extract relevant information
			user, _ := userInfo["user"].(map[string]interface{})
			email, _ := user["email"].(string)
			name, _ := user["name"].(string)
			currentOrgID := config.GetOrgID()

			// Print status information
			output.Success(cmd, "Authenticated")
			fmt.Printf("User: %s (%s)\n", name, email)
			fmt.Printf("Current Organization ID: %s\n", currentOrgID)
		},
	}
}

func authLogoutCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Logout from NVIDIA Cloud",
		Run: func(cmd *cobra.Command, args []string) {
			if !config.IsAuthenticated() {
				output.Info(cmd, "You are currently not logged in")
				return
			}
			config.ClearAPIKey()
			output.Success(cmd, "Logged out successfully")
		},
	}
}

var whoamiURL = "/v2/users/me"

func authWhoAmICmd() *cobra.Command {
	return &cobra.Command{
		Use:   "whoami",
		Short: "Display information about the authenticated user",
		Run: func(cmd *cobra.Command, args []string) {
			client := api.NewClient(config.GetAPIKey())
			whoamiInfo := map[string]any{}
			err := client.Get(cmd.Context(), whoamiURL, nil, &whoamiInfo)
			if err != nil {
				output.Error(cmd, "Failed to fetch user information", err)
				return
			}

			jsonMode, _ := cmd.Flags().GetBool("json")
			if jsonMode {
				json.NewEncoder(cmd.OutOrStdout()).Encode(whoamiInfo)
				return
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
		},
	}
}

func authOrgsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "orgs",
		Short: "Display organization and team information for the authenticated user",
		Run: func(cmd *cobra.Command, args []string) {
			client := api.NewClient(config.GetAPIKey())
			userInfo := map[string]interface{}{}
			err := client.Get(cmd.Context(), "/v2/users/me", nil, &userInfo)
			if err != nil {
				output.Error(cmd, "Failed to fetch user information", err)
				return
			}
			jsonMode, _ := cmd.Flags().GetBool("json")
			if jsonMode {
				json.NewEncoder(cmd.OutOrStdout()).Encode(userInfo)
				return
			}
			userRoles, ok := userInfo["userRoles"].([]interface{})
			if !ok {
				output.Error(cmd, "Failed to parse user roles information", nil)
				return
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
		},
	}

	cmd.Flags().BoolP("wide", "o", false, "Display wide output including org roles")
	return cmd
}

func authOrgIDCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "org-id",
		Short: "Display the name of the first organization",
		Run: func(cmd *cobra.Command, args []string) {
			client := api.NewClient(config.GetAPIKey())
			orgsInfo := map[string]interface{}{}
			err := client.Get(cmd.Context(), "/v2/orgs", nil, &orgsInfo)
			if err != nil {
				output.Error(cmd, "Failed to fetch organization information", err)
				return
			}

			organizations, ok := orgsInfo["organizations"].([]interface{})
			if !ok || len(organizations) == 0 {
				output.Error(cmd, "No organizations found", nil)
				return
			}

			firstOrg, ok := organizations[0].(map[string]interface{})
			if !ok {
				output.Error(cmd, "Failed to parse organization information", nil)
				return
			}

			name, ok := firstOrg["name"].(string)
			if !ok {
				output.Error(cmd, "Organization name not found", nil)
				return
			}

			fmt.Println(name)
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

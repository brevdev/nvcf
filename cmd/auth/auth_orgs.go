package auth

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/brevdev/nvcf/api"
	"github.com/brevdev/nvcf/config"
	"github.com/brevdev/nvcf/output"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

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

func convertToStringSlice(slice []interface{}) []string {
	result := make([]string, len(slice))
	for i, v := range slice {
		result[i] = fmt.Sprint(v)
	}
	return result
}

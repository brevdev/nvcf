package debug

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/brevdev/nvcf/api"
	"github.com/brevdev/nvcf/cmd/preflight/brev"
	"github.com/brevdev/nvcf/config"
	"github.com/brevdev/nvcf/output"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/tmc/nvcf-go"
)

func debugStartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start <function-id>",
		Short: "Start a debug environment",
		Long: `Create and start a debug environment for an NVCF function.

This command sets up a GPU-powered VM using Brev CLI, pulls the function's
container image, and runs it in a debugging session. It allows you to:
- Debug functions with ERROR status
- Specify a particular version for debugging
- Access the debug environment via SSH for hands-on troubleshooting

Use this command when you need to investigate issues with your NVCF function in a live environment.`,
		Args: cobra.ExactArgs(1),
		RunE: runDebugStart,
	}

	cmd.Flags().String("version-id", "", "The ID of the version")
	return cmd
}

func runDebugStart(cmd *cobra.Command, args []string) error {
	nvcfClient := api.NewClient(config.GetAPIKey())
	brevClient := brev.NewBrevClient()

	functionId := args[0]
	versionId, _ := cmd.Flags().GetString("version-id")

	if versionId == "" {
		versions, err := nvcfClient.Functions.Versions.List(cmd.Context(), functionId)
		if err != nil {
			return output.Error(cmd, "Error listing function versions", err)
		}

		errorVersions := []string{}
		for _, version := range versions.Functions {
			if version.Status == nvcf.ListFunctionsResponseFunctionsStatusError {
				errorVersions = append(errorVersions, version.VersionID)
			}
		}

		if len(errorVersions) == 0 {
			return output.Error(cmd, "No versions with ERROR status found", nil)
		}

		if len(errorVersions) == 1 {
			versionId = errorVersions[0]
		} else {
			output.Info(cmd, "Multiple versions with ERROR status found. Please select a version to debug:")
			for _, version := range errorVersions {
				output.Info(cmd, fmt.Sprintf("Version ID: %s", version))
			}
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Enter version-id to debug: ")
			versionId, _ = reader.ReadString('\n')
			versionId = strings.TrimSpace(versionId)
		}
	} else {
		targetVersion, err := nvcfClient.FunctionDeployment.Functions.Versions.GetDeployment(cmd.Context(), functionId, versionId)
		if err != nil {
			return output.Error(cmd, "Error getting deployment", err)
		}
		if targetVersion.Deployment.FunctionStatus != nvcf.DeploymentResponseDeploymentFunctionStatusError {
			return output.Error(cmd, "Deployment is not active", nil)
		}
	}

	// get function deployment information
	deployment, err := nvcfClient.Functions.Versions.Get(cmd.Context(), functionId, versionId, nvcf.FunctionVersionGetParams{
		IncludeSecrets: nvcf.Bool(false),
	})
	if err != nil {
		return output.Error(cmd, "Error getting deployment", err)
	}

	image := deployment.Function.ContainerImage
	imageArgs := deployment.Function.ContainerArgs

	if !brevClient.IsBrevCLIInstalled() {
		return fmt.Errorf("brev CLI is not installed. Please install it first")
	}

	loggedIn, err := brevClient.IsLoggedIn()
	if err != nil {
		return err
	}

	if !loggedIn {
		fmt.Println("You are not logged in. Starting Brev login process...")
		err = brevClient.Login()
		if err != nil {
			return err
		}
		fmt.Println("Successfully logged in with Brev CLI")
	}
	fmt.Println("Setting up a GPU powered VM for debugging")

	// i want to append a random 5 character string to the end of the instance name
	randomString := uuid.New().String()[:5]
	instanceName := fmt.Sprintf("nvcf-%s-debug-%s", functionId, randomString)

	// hit the brev api to create an instance using
	if err := brevClient.CreateInstance(functionId, instanceName); err != nil {
		return output.Error(cmd, "Error creating Brev instance", err)
	}

	// run the debugging script on the instance
	err = brevClient.RunDebuggingScript(instanceName, image, imageArgs)
	if err != nil {
		return output.Error(cmd, "Error running debugging script", err)
	}

	fmt.Printf("You can ssh into this instance using ssh %s\n", instanceName)

	return nil
}

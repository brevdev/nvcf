package localdeploymenttest

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/brevdev/nvcf/containerutil"
	"github.com/spf13/cobra"
)

var (
	imageName               string
	containerPort           string
	protocol                string
	healthEndpoint          string
	secondsToWaitForHealthy int
	forceCleanup            bool
	grpcServiceName         string
	grpcMethodName          string
	grpcInputData           string
)

// LocalDeploymentTestCmd represents the local-deployment-test command
var LocalDeploymentTestCmd = &cobra.Command{
	Use:   "local-deployment-test",
	Short: "Run local deployment test",
	Long:  `Run local deployment test to verify container compatibility with NVCF`,
	Run:   runLocalDeploymentTest,
}

func init() {
	LocalDeploymentTestCmd.Flags().StringVar(&imageName, "image-name", "", "Name of the Docker image")
	LocalDeploymentTestCmd.Flags().StringVar(&containerPort, "container-port", "8000", "Port that the server is listening on")
	LocalDeploymentTestCmd.Flags().StringVar(&protocol, "protocol", "http", "Protocol that the server is running (http or grpc)")
	LocalDeploymentTestCmd.Flags().StringVar(&healthEndpoint, "health-endpoint", "/v2/health/ready", "Health endpoint exposed by the server (for HTTP)")
	LocalDeploymentTestCmd.Flags().IntVar(&secondsToWaitForHealthy, "seconds-to-wait-for-healthy", 600, "How long to wait for the health endpoint to be ready")
	LocalDeploymentTestCmd.Flags().BoolVar(&forceCleanup, "force-cleanup", false, "Force cleanup of existing containers before starting")
	LocalDeploymentTestCmd.Flags().StringVar(&grpcServiceName, "grpc-service-name", "", "gRPC service name (required for gRPC)")
	LocalDeploymentTestCmd.Flags().StringVar(&grpcMethodName, "grpc-method-name", "", "gRPC method name (required for gRPC)")
	LocalDeploymentTestCmd.Flags().StringVar(&grpcInputData, "grpc-input-data", "{}", "JSON string representing input data for gRPC method")

	LocalDeploymentTestCmd.MarkFlagRequired("image-name")
}

func runLocalDeploymentTest(cmd *cobra.Command, args []string) {
	if protocol == "grpc" && (grpcServiceName == "" || grpcMethodName == "") {
		fmt.Println("Error: gRPC service name and method name are required for gRPC protocol")
		os.Exit(1)
	}

	cst, err := containerutil.NewContainerSmokeTest()
	if err != nil {
		fmt.Printf("Error creating ContainerSmokeTest: %v\n", err)
		os.Exit(1)
	}

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\nReceived interrupt signal. Cleaning up...")
		cst.Cleanup()
		os.Exit(1)
	}()

	// Force cleanup if requested
	if forceCleanup {
		fmt.Println("Forcing cleanup of existing containers...")
		if err := cst.ForceCleanup(imageName); err != nil {
			fmt.Printf("Error during force cleanup: %v\n", err)
			os.Exit(1)
		}
	}

	defer cst.Cleanup()

	err = cst.LaunchContainer(imageName, containerPort)
	if err != nil {
		fmt.Printf("Error launching container: %v\n", err)
		os.Exit(1)
	}

	if protocol == "http" {
		err = cst.CheckHTTPHealthEndpoint(healthEndpoint, secondsToWaitForHealthy)
	} else {
		err = cst.CheckGRPCHealthEndpoint(secondsToWaitForHealthy)
	}

	if err != nil {
		fmt.Printf("Error checking health endpoint: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Health Check succeeded!")

	if protocol == "grpc" {
		var inputData map[string]interface{}
		err := json.Unmarshal([]byte(grpcInputData), &inputData)
		if err != nil {
			fmt.Printf("Error parsing gRPC input data: %v\n", err)
			os.Exit(1)
		}

		err = cst.TestGRPCInference(grpcServiceName, grpcMethodName, inputData, false)
		if err != nil {
			fmt.Printf("Error testing gRPC inference: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Println("HTTP inference testing not implemented in this version.")
	}
}
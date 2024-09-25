package check

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/brevdev/nvcf/preflight/containerutil"
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
	httpInferenceEndpoint   string
	httpPayload             string
)

// This command runs a local deployment test to ensure the specified Docker image
// can be successfully deployed and accessed.
func CheckCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check <image-name>",
		Short: "Run local deployment test for NVCF compatibility",
		Long: `Run a local deployment test to verify container compatibility with NVCF.
This command deploys the specified Docker image locally, checks its health,
and performs basic connectivity tests to ensure it meets NVCF requirements.`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			imageName = args[0]
			runLocalDeploymentTest()
		},
	}
	cmd.Flags().StringVar(&containerPort, "container-port", "8000", "Port that the server is listening on")
	cmd.Flags().StringVar(&protocol, "protocol", "http", "Protocol that the server is running (http or grpc)")
	cmd.Flags().StringVar(&healthEndpoint, "health-endpoint", "/v2/health/ready", "Health endpoint exposed by the server (for HTTP)")
	cmd.Flags().IntVar(&secondsToWaitForHealthy, "seconds-to-wait-for-healthy", 600, "Maximum time to wait for the health endpoint to be ready (in seconds)")
	cmd.Flags().BoolVar(&forceCleanup, "force-cleanup", false, "Force cleanup of existing containers before starting the test")
	cmd.Flags().StringVar(&grpcServiceName, "grpc-service-name", "", "gRPC service name (required for gRPC)")
	cmd.Flags().StringVar(&grpcMethodName, "grpc-method-name", "", "gRPC method name (required for gRPC)")
	cmd.Flags().StringVar(&grpcInputData, "grpc-input-data", "{}", "JSON string representing input data for gRPC method")
	cmd.Flags().StringVar(&httpInferenceEndpoint, "http-inference-endpoint", "/", "HTTP inference endpoint")
	cmd.Flags().StringVar(&httpPayload, "http-payload", "{}", "JSON string representing input data for HTTP inference")
	return cmd
}

func runLocalDeploymentTest() {
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

	if protocol == "http" {
		var payload interface{}
		err := json.Unmarshal([]byte(httpPayload), &payload)
		if err != nil {
			fmt.Printf("Error parsing HTTP payload: %v\n", err)
			os.Exit(1)
		}

		err = cst.TestHTTPInference(httpInferenceEndpoint, payload)
		if err != nil {
			fmt.Printf("Error testing HTTP inference: %v\n", err)
			fmt.Println("HTTP inference test failed. Please check your application's endpoints and try again.")
			os.Exit(1)
			fmt.Println("HTTP inference test succeeded!")
		} else if protocol == "grpc" {
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
		}
	}
}

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
and performs basic connectivity tests to ensure it meets NVCF requirements.

Key features:
- Supports both HTTP and gRPC protocols
- Customizable health check and inference endpoints
- Configurable wait times for container readiness
- Option to force cleanup of existing containers

Use this command to validate your NVCF function before deployment.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			imageName = args[0]
			return runLocalDeploymentTest()
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

func runLocalDeploymentTest() error {
	if protocol == "grpc" && (grpcServiceName == "" || grpcMethodName == "") {
		return fmt.Errorf("gRPC service name and method name are required for gRPC protocol")
	}
	cst, err := containerutil.NewContainerSmokeTest()
	if err != nil {
		return fmt.Errorf("error creating ContainerSmokeTest: %w", err)
	}
	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\nReceived interrupt signal. Cleaning up...")
		_ = cst.Cleanup()
	}()
	// Force cleanup if requested
	if forceCleanup {
		fmt.Println("Forcing cleanup of existing containers...")
		if err := cst.ForceCleanup(imageName); err != nil {
			return fmt.Errorf("error during force cleanup: %w", err)
		}
	}
	defer func() {
		_ = cst.Cleanup()
	}()
	if err := cst.LaunchContainer(imageName, containerPort); err != nil {
		return fmt.Errorf("error launching container: %w", err)
	}
	var healthErr error
	if protocol == "http" {
		healthErr = cst.CheckHTTPHealthEndpoint(healthEndpoint, secondsToWaitForHealthy)
	} else {
		healthErr = cst.CheckGRPCHealthEndpoint(secondsToWaitForHealthy)
	}
	if healthErr != nil {
		return fmt.Errorf("error checking health endpoint: %w", healthErr)
	}
	fmt.Println("Health Check succeeded!")
	if protocol == "http" {
		var payload interface{}
		if err := json.Unmarshal([]byte(httpPayload), &payload); err != nil {
			return fmt.Errorf("error parsing HTTP payload: %w", err)
		}
		if err := cst.TestHTTPInference(httpInferenceEndpoint, payload); err != nil {
			return fmt.Errorf("error testing HTTP inference: %w", err)
		}
		fmt.Println("HTTP inference test succeeded!")
	} else if protocol == "grpc" {
		var inputData map[string]interface{}
		if err := json.Unmarshal([]byte(grpcInputData), &inputData); err != nil {
			return fmt.Errorf("error parsing gRPC input data: %w", err)
		}
		if err := cst.TestGRPCInference(grpcServiceName, grpcMethodName, inputData, false); err != nil {
			return fmt.Errorf("error testing gRPC inference: %w", err)
		}
		fmt.Println("gRPC inference test succeeded!")
	}
	return nil
}

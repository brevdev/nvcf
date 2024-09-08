// TODO
// - add support for asset/model mounting
// - add support for env vars
package function

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/brevdev/nvcf/api"
	"github.com/brevdev/nvcf/config"
	"github.com/brevdev/nvcf/output"
	"github.com/spf13/cobra"
	"github.com/tmc/nvcf-go"
)

// functionCreateCmd creates a new Cobra command for function creation and deployment
func functionCreateCmd() *cobra.Command {
	// Define command-line flags
	var (
		// Function creation parameters
		name           string
		inferenceURL   string
		inferencePort  int64
		healthUri      string
		containerImage string
		containerArgs  string
		description    string
		tags           []string
		apiBodyFormat  string
		functionType   string

		// Health check parameters
		healthProtocol   string
		healthPort       int64
		healthTimeout    string
		healthStatusCode int64

		// Deployment parameters
		minInstances          int64
		maxInstances          int64
		gpu                   string
		instanceType          string
		backend               string
		maxRequestConcurrency int64
		deploy                bool
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new function",
		Long:  `Create a new NVIDIA Cloud Function with the specified parameters.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Create a new API client
			client := api.NewClient(config.GetAPIKey())

			// Prepare function creation parameters
			params := prepareFunctionParams(name, inferenceURL, inferencePort, healthUri, containerImage, apiBodyFormat, description, tags, functionType, healthProtocol, healthPort, healthTimeout, healthStatusCode, containerArgs)

			output.Info(cmd, "Creating function")

			// Create the function
			resp, err := client.Functions.New(cmd.Context(), params)
			if err != nil {
				return fmt.Errorf("error creating function: %w", err)
			}

			if !deploy {
				output.Success(cmd, fmt.Sprintf("Function with FunctionID %s and VersionID %s created successfully", resp.Function.ID, resp.Function.VersionID))
				return nil
			}

			// Deploy the function if the deploy flag is set
			return deployFunction(cmd, client, resp, gpu, instanceType, backend, maxInstances, minInstances, maxRequestConcurrency)
		},
	}

	// function create flags
	cmd.Flags().StringVar(&name, "name", "", "Name of the function (required)")
	cmd.Flags().StringVar(&inferenceURL, "inference-url", "", "URL for function invocation (required)")
	cmd.Flags().Int64Var(&inferencePort, "inference-port", 80, "Port for function invocation")
	cmd.Flags().StringVar(&healthUri, "health-uri", "/health", "Health check URI")
	cmd.Flags().StringVar(&containerImage, "container-image", "", "Container image for the function")
	cmd.Flags().StringVar(&containerArgs, "container-args", "", "Container arguments")
	cmd.Flags().StringVar(&apiBodyFormat, "api-body-format", "CUSTOM", "API body format (PREDICT_V2 or CUSTOM)")
	cmd.Flags().StringVar(&description, "description", "", "Description of the function")
	cmd.Flags().StringSliceVar(&tags, "tag", nil, "Tags for the function (can be used multiple times)")
	cmd.Flags().StringVar(&functionType, "function-type", "DEFAULT", "Function type (DEFAULT or STREAMING)")
	// optional health specification flags
	cmd.Flags().StringVar(&healthProtocol, "health-protocol", "HTTP", "Health check protocol (HTTP or GRPC)")
	cmd.Flags().Int64Var(&healthPort, "health-port", 80, "Health check port")
	cmd.Flags().StringVar(&healthTimeout, "health-timeout", "5s", "Health check timeout")
	cmd.Flags().Int64Var(&healthStatusCode, "health-status-code", 200, "Expected health check status code")

	// deployment flags
	cmd.Flags().Int64Var(&minInstances, "min-instances", 0, "Minimum number of instances")
	cmd.Flags().Int64Var(&maxInstances, "max-instances", 0, "Maximum number of instances")
	cmd.Flags().StringVar(&gpu, "gpu", "H100", "GPU type to use")
	cmd.Flags().StringVar(&instanceType, "instance-type", "GCP.GPU.H100_1x", "Instance type to use")
	cmd.Flags().StringVar(&backend, "backend", "gcp-asia-se-1a", "Backend to deploy the function to (see NGC for available backends)")
	cmd.Flags().Int64Var(&maxRequestConcurrency, "max-request-concurrency", 0, "Maximum number of concurrent requests")
	cmd.Flags().BoolVar(&deploy, "deploy", false, "Create and deploy the function in one step")

	// these follow current nvcf-ci specs
	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("inference-url")
	cmd.MarkFlagRequired("inference-port")
	cmd.MarkFlagRequired("health-uri")
	cmd.MarkFlagRequired("container-image")

	return cmd
}

// prepareFunctionParams creates and returns the FunctionNewParams struct
func prepareFunctionParams(name, inferenceURL string, inferencePort int64, healthUri, containerImage, apiBodyFormat, description string,
	tags []string, functionType, healthProtocol string, healthPort int64, healthTimeout string, healthStatusCode int64, containerArgs string) nvcf.FunctionNewParams {
	return nvcf.FunctionNewParams{
		Name:           nvcf.String(name),
		InferenceURL:   nvcf.String(inferenceURL),
		InferencePort:  nvcf.Int(inferencePort),
		HealthUri:      nvcf.String(healthUri),
		ContainerImage: nvcf.String(containerImage),
		ContainerArgs:  nvcf.String(containerArgs),
		APIBodyFormat:  nvcf.F(nvcf.FunctionNewParamsAPIBodyFormat(apiBodyFormat)),
		Description:    nvcf.F(description),
		Tags:           nvcf.F(tags),
		FunctionType:   nvcf.F(nvcf.FunctionNewParamsFunctionType(functionType)),
		Health: nvcf.F(nvcf.FunctionNewParamsHealth{
			Protocol:           nvcf.F(nvcf.FunctionNewParamsHealthProtocol(healthProtocol)),
			Port:               nvcf.F(healthPort),
			Timeout:            nvcf.F(healthTimeout),
			ExpectedStatusCode: nvcf.F(healthStatusCode),
		}),
	}
}

// deployFunction handles the deployment of the created function
func deployFunction(cmd *cobra.Command, client *api.Client, resp *nvcf.CreateFunctionResponse, gpu, instanceType, backend string,
	maxInstances, minInstances, maxRequestConcurrency int64) error {
	output.Info(cmd, "Deployment flag was provided. Deploying function...")

	deploymentParams := nvcf.FunctionDeploymentFunctionVersionNewParams{
		DeploymentSpecifications: nvcf.F([]nvcf.FunctionDeploymentFunctionVersionNewParamsDeploymentSpecification{{
			GPU:                   nvcf.String(gpu),
			InstanceType:          nvcf.String(instanceType),
			Backend:               nvcf.String(backend),
			MaxInstances:          nvcf.Int(maxInstances),
			MinInstances:          nvcf.Int(minInstances),
			MaxRequestConcurrency: nvcf.Int(maxRequestConcurrency),
		}}),
	}

	_, err := client.FunctionDeployment.Functions.Versions.New(
		cmd.Context(),
		resp.Function.ID,
		resp.Function.VersionID,
		deploymentParams,
	)
	if err != nil {
		return fmt.Errorf("error deploying function: %w", err)
	}

	output.Success(cmd, fmt.Sprintf("Function with FunctionID %s and VersionID %s deployed successfully", resp.Function.ID, resp.Function.VersionID))

	var fn nvcf.ListFunctionsResponseFunction
	if err := jsonMarshalUnmarshal(&fn, resp.Function); err != nil {
		return fmt.Errorf("issue marshaling+unmarshaling: %w", err)
	}
	output.Function(cmd, fn)
	return nil
}

// jsonMarshalUnmarshal marshals the src to JSON and then unmarshals it into dest
func jsonMarshalUnmarshal(dest any, src any) error {
	// Validate dest is a pointer
	if reflect.ValueOf(dest).Kind() != reflect.Ptr {
		return fmt.Errorf("destination must be a pointer")
	}
	jsonData, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonData, dest)
}

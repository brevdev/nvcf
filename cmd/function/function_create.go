// TODO
// - add support for asset/model mounting
// - add support for env vars
// - Implement HelmChart and HelmChartServiceName functionality
// - Add support for Models array
// - Implement Resources array functionality
// - Add support for Secrets array
package function

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/brevdev/nvcf/api"
	"github.com/brevdev/nvcf/config"
	"github.com/brevdev/nvcf/output"
	"github.com/spf13/cobra"
	"github.com/tmc/nvcf-go"
	"gopkg.in/yaml.v3"
)

func functionCreateCmd() *cobra.Command {
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
		envVars        []string

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

		// Optional function specification file
		fileSpec string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new function",
		Long:  `Create a new NVIDIA Cloud Function with the specified parameters.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client := api.NewClient(config.GetAPIKey())

			if fileSpec != "" {
				return createFunctionsFromFile(cmd, client, fileSpec)
			}

			containerEnv, err := parseEnvVars(envVars)
			if err != nil {
				return fmt.Errorf("error parsing environment variables: %w", err)
			}

			params := prepareFunctionParams(name, inferenceURL, inferencePort, healthUri, containerImage, apiBodyFormat, description, tags, functionType, healthProtocol, healthPort, healthTimeout, healthStatusCode, containerArgs, containerEnv)
			output.Info(cmd, fmt.Sprintf("Creating function %s...", name))

			// create function
			resp, err := client.Functions.New(cmd.Context(), params)
			if err != nil {
				return fmt.Errorf("error creating function: %w", err)
			}
			output.Success(cmd, fmt.Sprintf("Function %s with id %s and version %s created successfully", name, resp.Function.ID, resp.Function.VersionID))

			// deploy function if the deploy flag is set
			if !deploy {
				return deployFunction(cmd, client, resp, gpu, instanceType, backend, maxInstances, minInstances, maxRequestConcurrency)
			}

			return nil
		},
	}

	// function create flags
	cmd.Flags().StringVar(&name, "name", "", "Name of the function (required)")
	cmd.Flags().StringVar(&inferenceURL, "inference-url", "", "URL for function invocation (required)")
	cmd.Flags().Int64Var(&inferencePort, "inference-port", 80, "Port for function invocation. Default is 80")
	cmd.Flags().StringVar(&healthUri, "health-uri", "/health", "Health check URI. Default is /health")
	cmd.Flags().StringVar(&containerImage, "container-image", "", "Container image for the function")
	cmd.Flags().StringVar(&containerArgs, "container-args", "", "Container arguments. Put these in quotes if you are passing flags")
	cmd.Flags().StringVar(&apiBodyFormat, "api-body-format", defaultAPIBodyFormat, "API body format (PREDICT_V2 or CUSTOM). Default is CUSTOM")
	cmd.Flags().StringVar(&description, "description", "", "Description of the function")
	cmd.Flags().StringSliceVar(&tags, "tag", nil, "Tags for the function (can be used multiple times)")
	cmd.Flags().StringVar(&functionType, "function-type", defaultFunctionType, "Function type (DEFAULT or STREAMING). Default is DEFAULT")
	cmd.Flags().StringSliceVar(&envVars, "env", []string{}, "Environment variables for the function (can be used multiple times, format: key:value)")

	// optional health specification flags
	cmd.Flags().StringVar(&healthProtocol, "health-protocol", "HTTP", "Health check protocol (HTTP or GRPC). Default is HTTP")
	cmd.Flags().Int64Var(&healthPort, "health-port", 80, "Health check port. Default is 80")
	cmd.Flags().StringVar(&healthTimeout, "health-timeout", "PT20S", "Health check timeout. Default is PT20S") //TODO: allow user to specify in s and we convert
	cmd.Flags().Int64Var(&healthStatusCode, "health-status-code", 200, "Expected health check status code. Default is 200")

	// deployment flags
	cmd.Flags().Int64Var(&minInstances, "min-instances", 0, "Minimum number of instances. Default is 0")
	cmd.Flags().Int64Var(&maxInstances, "max-instances", 1, "Maximum number of instances. Default is 1")
	cmd.Flags().StringVar(&gpu, "gpu", "H100", "GPU type to use. Default is H100")
	cmd.Flags().StringVar(&instanceType, "instance-type", "GCP.GPU.H100_1x", "Instance type to use. Default is GCP.GPU.H100_1x")
	cmd.Flags().StringVar(&backend, "backend", "gcp-asia-se-1a", "Backend to deploy the function to (see NGC for available backends). Default is gcp-asia-se-1a")
	cmd.Flags().Int64Var(&maxRequestConcurrency, "max-request-concurrency", 1, "Maximum number of concurrent requests. Default is 1")
	cmd.Flags().BoolVar(&deploy, "deploy", false, "Create and deploy the function in one step. Default is false")
	cmd.Flags().StringVarP(&fileSpec, "file", "f", "", "Path to a YAML file containing function specifications")

	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("inference-url")
	cmd.MarkFlagRequired("inference-port")
	cmd.MarkFlagRequired("health-uri")
	cmd.MarkFlagRequired("container-image")

	return cmd
}

func parseEnvVars(envVars []string) ([]nvcf.FunctionNewParamsContainerEnvironment, error) {
	var containerEnv []nvcf.FunctionNewParamsContainerEnvironment

	for _, env := range envVars {
		parts := strings.SplitN(env, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid environment variable format: %s", env)
		}
		containerEnv = append(containerEnv, nvcf.FunctionNewParamsContainerEnvironment{
			Key:   nvcf.F(parts[0]),
			Value: nvcf.F(parts[1]),
		})
	}

	return containerEnv, nil
}

func prepareFunctionParams(name, inferenceURL string, inferencePort int64, healthUri, containerImage, apiBodyFormat, description string,
	tags []string, functionType, healthProtocol string, healthPort int64, healthTimeout string, healthStatusCode int64, containerArgs string, containerEnv []nvcf.FunctionNewParamsContainerEnvironment) nvcf.FunctionNewParams {
	params := nvcf.FunctionNewParams{
		Name:                 nvcf.String(name),
		InferenceURL:         nvcf.String(inferenceURL),
		InferencePort:        nvcf.Int(inferencePort),
		ContainerImage:       nvcf.String(containerImage),
		ContainerArgs:        nvcf.String(containerArgs),
		ContainerEnvironment: nvcf.F(containerEnv),
		APIBodyFormat:        nvcf.F(nvcf.FunctionNewParamsAPIBodyFormat(apiBodyFormat)),
		Description:          nvcf.F(description),
		Tags:                 nvcf.F(tags),
		FunctionType:         nvcf.F(nvcf.FunctionNewParamsFunctionType(functionType)),
		Health: nvcf.F(nvcf.FunctionNewParamsHealth{
			Protocol:           nvcf.F(nvcf.FunctionNewParamsHealthProtocol(healthProtocol)),
			Port:               nvcf.F(healthPort),
			Timeout:            nvcf.F(healthTimeout),
			ExpectedStatusCode: nvcf.F(healthStatusCode),
			Uri:                nvcf.String(healthUri),
		}),
	}
	return params
}

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

func createFunctionsFromFile(cmd *cobra.Command, client *api.Client, yamlFile string) error {
	data, err := os.ReadFile(yamlFile)
	if err != nil {
		return fmt.Errorf("error reading YAML file: %w", err)
	}

	var spec FunctionSpec
	if err := yaml.Unmarshal(data, &spec); err != nil {
		return fmt.Errorf("error parsing YAML file: %w", err)
	}

	for _, fn := range spec.Functions {
		params := prepareFunctionParamsFromFile(spec.FnImage, fn)
		if err := createAndDeployFunctionFromFile(cmd, client, params, true, fn.InstGPUType, fn.InstType, fn.InstBackend, fn.InstMax, fn.InstMin, fn.InstMaxRequestConcurrency); err != nil {
			return err
		}
	}

	return nil
}

func prepareFunctionParamsFromFile(fnImage string, fn FunctionDef) nvcf.FunctionNewParams {
	// Use provided values if present, otherwise use defaults
	apiBodyFormat := fn.APIBodyFormat
	if apiBodyFormat == "" {
		apiBodyFormat = defaultAPIBodyFormat
	}

	functionType := fn.FunctionType
	if functionType == "" {
		functionType = defaultFunctionType
	}

	return nvcf.FunctionNewParams{
		Name:           nvcf.String(fn.FnName),
		InferenceURL:   nvcf.String(fn.InferenceURL),
		InferencePort:  nvcf.Int(fn.InferencePort),
		ContainerImage: nvcf.String(fnImage),
		ContainerArgs:  nvcf.String(fn.ContainerArgs),
		APIBodyFormat:  nvcf.F(nvcf.FunctionNewParamsAPIBodyFormat(apiBodyFormat)),
		Description:    nvcf.F(fn.Description),
		Tags:           nvcf.F(fn.Tags),
		FunctionType:   nvcf.F(nvcf.FunctionNewParamsFunctionType(functionType)),
		Health: nvcf.F(nvcf.FunctionNewParamsHealth{
			Protocol:           nvcf.F(nvcf.FunctionNewParamsHealthProtocol(fn.Health.Protocol)),
			Port:               nvcf.F(fn.Health.Port),
			Timeout:            nvcf.F(fn.Health.Timeout),
			ExpectedStatusCode: nvcf.F(fn.Health.StatusCode),
			Uri:                nvcf.String(fn.Health.URI),
		}),
	}
}

func createAndDeployFunctionFromFile(cmd *cobra.Command, client *api.Client, params nvcf.FunctionNewParams, deploy bool, gpu, instanceType, backend string, maxInstances, minInstances, maxRequestConcurrency int64) error {
	resp, err := client.Functions.New(cmd.Context(), params)
	if err != nil {
		return fmt.Errorf("error creating function: %w", err)
	}

	if !deploy {
		output.Success(cmd, fmt.Sprintf("Function %s with id %s and version %s created successfully", params.Name, resp.Function.ID, resp.Function.VersionID))
		return nil
	}

	return deployFunction(cmd, client, resp, gpu, instanceType, backend, maxInstances, minInstances, maxRequestConcurrency)
}

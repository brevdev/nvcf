// TODO
// - Implement HelmChart and HelmChartServiceName functionality
// - Implement Resources array functionality
// - Add support for Secrets array
// - We cannot programatically get backends just yet. A user has to go to NGC to find their backend and gpu allocation
package function

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/brevdev/nvcf/api"
	"github.com/brevdev/nvcf/config"
	"github.com/brevdev/nvcf/flagutil"
	"github.com/brevdev/nvcf/output"
	"github.com/spf13/cobra"
	"github.com/tmc/nvcf-go"
	"gopkg.in/yaml.v3"
)

// default values for function creation
const (
	defaultAPIBodyFormat = "CUSTOM"
	defaultFunctionType  = "STREAMING"
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
		custom         bool //if false this sets apiBodyFormat to PREDICT_V2
		streaming      bool //if false this sets functionType to DEFAULT
		functionType   string
		envVars        []string
		modelVars      []string

		// Health check parameters
		healthProtocol   string
		healthPort       int64
		healthTimeout    time.Duration
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

		// Optional new version flag
		existingFunctionID string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new function",
		Long:  `Create a new NVIDIA Cloud Function with the specified parameters. If you specify `,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			fileSpec, _ := cmd.Flags().GetString("file")
			if fileSpec == "" {
				requiredFlags := []string{"name", "inference-url", "inference-port", "health-uri", "container-image"}
				for _, flag := range requiredFlags {
					if err := cmd.MarkFlagRequired(flag); err != nil {
						return err
					}
				}
			}
			return nil
		}, RunE: func(cmd *cobra.Command, args []string) error {
			client := api.NewClient(config.GetAPIKey())

			if existingFunctionID != "" {
				_, err := client.Functions.Versions.List(cmd.Context(), existingFunctionID)
				if err != nil {
					output.Error(cmd, "Error listing function versions", err)
					return fmt.Errorf("error getting function: %w", err)
				}
			}

			if fileSpec != "" {
				return createFunctionsFromFile(cmd, client, fileSpec, deploy)
			}

			apiBodyFormat := defaultAPIBodyFormat
			if !custom {
				apiBodyFormat = "PREDICT_V2"
			}

			functionType := defaultFunctionType
			if streaming {
				functionType = "STREAMING"
			}

			if existingFunctionID != "" {
				containerEnv, err := parseEnvVarsNewVersion(envVars)
				if err != nil {
					return fmt.Errorf("error parsing environment variables: %w", err)
				}
				models, err := parseModelsNewVersion(modelVars)
				if err != nil {
					return fmt.Errorf("error parsing models: %w", err)
				}
				params := nvcf.FunctionVersionNewParams{
					Name:                 nvcf.String(name),
					InferenceURL:         nvcf.String(inferenceURL),
					InferencePort:        nvcf.Int(inferencePort),
					ContainerImage:       nvcf.String(containerImage),
					ContainerArgs:        nvcf.String(containerArgs),
					ContainerEnvironment: nvcf.F(containerEnv),
					APIBodyFormat:        nvcf.F(nvcf.FunctionVersionNewParamsAPIBodyFormat(apiBodyFormat)),
					Description:          nvcf.F(description),
					Tags:                 nvcf.F(tags),
					FunctionType:         nvcf.F(nvcf.FunctionVersionNewParamsFunctionType(functionType)),
					Models:               nvcf.F(models),
					Health: nvcf.F(nvcf.FunctionVersionNewParamsHealth{
						Protocol:           nvcf.F(nvcf.FunctionVersionNewParamsHealthProtocol(healthProtocol)),
						Port:               nvcf.F(healthPort),
						Timeout:            nvcf.F(flagutil.DurationToISO8601(healthTimeout)),
						ExpectedStatusCode: nvcf.F(healthStatusCode),
						Uri:                nvcf.String(healthUri),
					}),
				}
				output.Info(cmd, fmt.Sprintf("Creating new version for function %s...", name))
				// create function
				resp, err := client.Functions.Versions.New(cmd.Context(), existingFunctionID, params)
				if err != nil {
					return fmt.Errorf("error creating function: %w", err)
				}
				output.Success(cmd, fmt.Sprintf("Function version %s created successfully", resp.Function.VersionID))
				// deploy function if the deploy flag is set
				if deploy {
					return deployFunction(cmd, client, resp, gpu, instanceType, backend, maxInstances, minInstances, maxRequestConcurrency)
				}
			} else {
				containerEnv, err := parseEnvVars(envVars)
				if err != nil {
					return fmt.Errorf("error parsing environment variables: %w", err)
				}
				models, err := parseModels(modelVars)
				if err != nil {
					return fmt.Errorf("error parsing models: %w", err)
				}
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
					Models:               nvcf.F(models),
					Health: nvcf.F(nvcf.FunctionNewParamsHealth{
						Protocol:           nvcf.F(nvcf.FunctionNewParamsHealthProtocol(healthProtocol)),
						Port:               nvcf.F(healthPort),
						Timeout:            nvcf.F(flagutil.DurationToISO8601(healthTimeout)),
						ExpectedStatusCode: nvcf.F(healthStatusCode),
						Uri:                nvcf.String(healthUri),
					}),
				}
				output.Info(cmd, fmt.Sprintf("Creating new function %s...", name))
				// create function
				resp, err := client.Functions.New(cmd.Context(), params)
				if err != nil {
					return fmt.Errorf("error creating function: %w", err)
				}
				output.Success(cmd, fmt.Sprintf("Function %s with id %s and version %s created successfully", name, resp.Function.ID, resp.Function.VersionID))
				// deploy function if the deploy flag is set
				if deploy {
					return deployFunction(cmd, client, resp, gpu, instanceType, backend, maxInstances, minInstances, maxRequestConcurrency)
				}
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
	cmd.Flags().BoolVar(&custom, "custom", true, "Set API body format to CUSTOM. If false - set API body format to PREDICT_V2.")
	cmd.Flags().StringVar(&description, "description", "", "Description of the function")
	cmd.Flags().StringSliceVar(&tags, "tag", nil, "Tags for the function (can be used multiple times)")
	cmd.Flags().BoolVar(&streaming, "streaming", true, "Set function type to STREAMING. Default is true")
	cmd.Flags().StringVar(&functionType, "function-type", defaultFunctionType, "Function type (DEFAULT or STREAMING). Default is DEFAULT")
	cmd.Flags().StringSliceVar(&envVars, "env", []string{}, "Environment variables for the function (can be used multiple times, format: key:value)")
	cmd.Flags().StringSliceVar(&modelVars, "model", []string{}, "Models for the function (can be used multiple times, format: name:uri:version)")

	//optional new version flag
	cmd.Flags().StringVar(&existingFunctionID, "new-version", "", "Create a new version of an existing function. Requires a valid function id")

	// optional health specification flags
	cmd.Flags().StringVar(&healthProtocol, "health-protocol", "HTTP", "Health check protocol (HTTP or GRPC). Default is HTTP")
	cmd.Flags().Int64Var(&healthPort, "health-port", 80, "Health check port. Default is 80")
	cmd.Flags().DurationVar(&healthTimeout, "health-timeout", 20*time.Second, "Health check timeout.") //TODO: allow user to specify in s and we convert
	cmd.Flags().Int64Var(&healthStatusCode, "health-status-code", 200, "Expected health check status code. Default is 200")

	// deployment flags
	cmd.Flags().Int64Var(&minInstances, "min-instances", 0, "Minimum number of instances. Default is 0")
	cmd.Flags().Int64Var(&maxInstances, "max-instances", 1, "Maximum number of instances. Default is 1")
	cmd.Flags().StringVar(&gpu, "gpu", "", "GPU type to use")
	cmd.Flags().StringVar(&instanceType, "instance-type", "", "Instance type to use. Default is GCP.GPU.H100_1x")
	cmd.Flags().StringVar(&backend, "backend", "", "Backend to deploy the function to (see your NGC org available backends)")
	cmd.Flags().Int64Var(&maxRequestConcurrency, "max-request-concurrency", 1, "Maximum number of concurrent requests. Default is 1")
	cmd.Flags().BoolVar(&deploy, "deploy", false, "Create and deploy the function in one step. Default is false")
	cmd.Flags().StringVarP(&fileSpec, "file", "f", "", "Path to a YAML file containing function specifications")

	return cmd
}

func parseEnvVars(envVars []string) ([]nvcf.FunctionNewParamsContainerEnvironment, error) {
	var containerEnv []nvcf.FunctionNewParamsContainerEnvironment

	for _, env := range envVars {
		parts := strings.SplitN(env, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid environment variable format: %s. ensure that you are using the format key:value", env)
		}
		containerEnv = append(containerEnv, nvcf.FunctionNewParamsContainerEnvironment{
			Key:   nvcf.F(parts[0]),
			Value: nvcf.F(parts[1]),
		})
	}

	return containerEnv, nil
}

func parseEnvVarsNewVersion(envVars []string) ([]nvcf.FunctionVersionNewParamsContainerEnvironment, error) {
	var containerEnv []nvcf.FunctionVersionNewParamsContainerEnvironment

	for _, env := range envVars {
		parts := strings.SplitN(env, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid environment variable format: %s. ensure that you are using the format key:value", env)
		}
		containerEnv = append(containerEnv, nvcf.FunctionVersionNewParamsContainerEnvironment{
			Key:   nvcf.F(parts[0]),
			Value: nvcf.F(parts[1]),
		})
	}

	return containerEnv, nil
}

func parseModels(modelVars []string) ([]nvcf.FunctionNewParamsModel, error) {
	var models []nvcf.FunctionNewParamsModel

	for _, model := range modelVars {
		parts := strings.Split(model, ":")
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid model format: %s", model)
		}
		models = append(models, nvcf.FunctionNewParamsModel{
			Name:    nvcf.F(parts[0]),
			Uri:     nvcf.F(parts[1]),
			Version: nvcf.F(parts[2]),
		})
	}

	return models, nil
}

func parseModelsNewVersion(modelVars []string) ([]nvcf.FunctionVersionNewParamsModel, error) {
	var models []nvcf.FunctionVersionNewParamsModel

	for _, model := range modelVars {
		parts := strings.Split(model, ":")
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid model format: %s", model)
		}
		models = append(models, nvcf.FunctionVersionNewParamsModel{
			Name:    nvcf.F(parts[0]),
			Uri:     nvcf.F(parts[1]),
			Version: nvcf.F(parts[2]),
		})
	}
	return models, nil
}

func deployFunction(cmd *cobra.Command, client *api.Client, resp *nvcf.CreateFunctionResponse, gpu, instanceType, backend string,
	maxInstances, minInstances, maxRequestConcurrency int64) error {
	output.Info(cmd, "Deployment flag was provided. Deploying function...")

	deploymentParams := nvcf.FunctionDeploymentFunctionVersionInitiateDeploymentParams{
		DeploymentSpecifications: nvcf.F([]nvcf.FunctionDeploymentFunctionVersionInitiateDeploymentParamsDeploymentSpecification{{
			GPU:                   nvcf.String(gpu),
			InstanceType:          nvcf.String(instanceType),
			Backend:               nvcf.String(backend),
			MaxInstances:          nvcf.Int(maxInstances),
			MinInstances:          nvcf.Int(minInstances),
			MaxRequestConcurrency: nvcf.Int(maxRequestConcurrency),
			// missing attributes, availabilityZones, clusters, configuration, preferredOrder, regions
		}}),
	}

	_, err := client.FunctionDeployment.Functions.Versions.InitiateDeployment(
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
	output.MultiFunction(cmd, fn)

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

func createFunctionsFromFile(cmd *cobra.Command, client *api.Client, yamlFile string, deploy bool) error {
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
		if err := createAndDeployFunctionFromFile(cmd, client, params, deploy, fn.InstGPUType, fn.InstType, fn.InstBackend, fn.InstMax, fn.InstMin, fn.InstMaxRequestConcurrency); err != nil {
			return err
		}
	}

	return nil
}

func prepareFunctionParamsFromFile(fnImage string, fn FunctionDef) nvcf.FunctionNewParams {
	// Use provided values if present, otherwise use defaults
	var apiBodyFormat string
	if fn.Custom {
		apiBodyFormat = defaultAPIBodyFormat
	}

	var functionType string
	if fn.Streaming {
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
			ExpectedStatusCode: nvcf.F(fn.Health.ExpectedStatusCode),
			Uri:                nvcf.String(fn.Health.Uri),
		}),
	}
}

func createAndDeployFunctionFromFile(cmd *cobra.Command, client *api.Client, params nvcf.FunctionNewParams, deploy bool, gpu, instanceType, backend string, maxInstances, minInstances, maxRequestConcurrency int64) error {
	resp, err := client.Functions.New(cmd.Context(), params)
	if err != nil {
		return fmt.Errorf("error creating function: %w", err)
	}

	if deploy {
		return deployFunction(cmd, client, resp, gpu, instanceType, backend, maxInstances, minInstances, maxRequestConcurrency)
	}

	output.Success(cmd, fmt.Sprintf("Function %s with id %s and version %s created successfully", params.Name, resp.Function.ID, resp.Function.VersionID))
	return nil
}

// TODO
// - Implement HelmChart and HelmChartServiceName functionality
// - Implement Resources array functionality
// - Add support for Secrets array
package function

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/brevdev/nvcf/api"
	"github.com/brevdev/nvcf/config"
	"github.com/brevdev/nvcf/flagutil"
	"github.com/brevdev/nvcf/output"
	"github.com/brevdev/nvcf/timeout"
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
		name               string
		inferenceURL       string
		inferencePort      int64
		healthUri          string
		containerImage     string
		containerArgs      string
		description        string
		tags               []string
		custom             bool //if false this sets apiBodyFormat to PREDICT_V2
		streaming          bool //if false this sets functionType to DEFAULT
		functionType       string
		envVars            []string
		modelVars          []string
		existingFunctionID string

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
		detatched             bool

		// Optional function specification file
		fileSpec string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new function",
		Long:  `Create a new NVCF Function with the specified parameters. If you specify --from-version, we will create a new version of an existing function. You can also create and deploy a function in one step using the --deploy flag.`,
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

			if fileSpec != "" {
				return createFunctionsFromFile(cmd, client, fileSpec, deploy)
			}

			if existingFunctionID != "" {
				_, err := client.Functions.Versions.List(cmd.Context(), existingFunctionID)
				if err != nil {
					return output.Error(cmd, "Error listing function versions", err)
				}
			}

			apiBodyFormat := defaultAPIBodyFormat
			if !custom {
				apiBodyFormat = "PREDICT_V2"
			}

			functionType := defaultFunctionType
			if !streaming {
				functionType = "DEFAULT"
			}

			if existingFunctionID != "" {
				containerEnv, err := parseEnvVarsNewVersion(cmd, envVars)
				if err != nil {
					return output.Error(cmd, "error parsing environment variables", err)
				}
				models, err := parseModelsNewVersion(cmd, modelVars)
				if err != nil {
					return output.Error(cmd, "error parsing models", err)
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
						Timeout:            nvcf.F(flagutil.DurationToISO8601(time.Duration(healthTimeout) * time.Second)),
						ExpectedStatusCode: nvcf.F(healthStatusCode),
						Uri:                nvcf.String(healthUri),
					}),
				}
				output.Info(cmd, fmt.Sprintf("Creating new version for function %s...", name))
				// create function
				resp, err := client.Functions.Versions.New(cmd.Context(), existingFunctionID, params)
				if err != nil {
					return output.Error(cmd, "error creating function", err)
				}
				output.Success(cmd, fmt.Sprintf("Function version %s created successfully", resp.Function.VersionID))
				// deploy function if the deploy flag is set
				if deploy {
					return deployFunction(cmd, client, resp, gpu, instanceType, backend, maxInstances, minInstances, maxRequestConcurrency)
				}
			} else {
				containerEnv, err := parseEnvVars(cmd, envVars)
				if err != nil {
					return output.Error(cmd, "error parsing environment variables", err)
				}
				models, err := parseModels(cmd, modelVars)
				if err != nil {
					return output.Error(cmd, "error parsing models", err)
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
						Timeout:            nvcf.F(flagutil.DurationToISO8601(time.Duration(healthTimeout) * time.Second)),
						ExpectedStatusCode: nvcf.F(healthStatusCode),
						Uri:                nvcf.String(healthUri),
					}),
				}
				output.Info(cmd, fmt.Sprintf("Creating new function %s...", name))
				resp, err := client.Functions.New(cmd.Context(), params)
				if err != nil {
					return output.Error(cmd, "error creating function", err)
				}
				output.Success(cmd, fmt.Sprintf("Function %s with id %s and version %s created successfully", name, resp.Function.ID, resp.Function.VersionID))
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
	cmd.Flags().StringVar(&existingFunctionID, "from-version", "", "Create a new version of an existing function. Requires a valid function id")

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
	cmd.Flags().BoolVarP(&detatched, "detatched", "d", false, "Deploy the function in the background. Default is false")
	return cmd
}

func parseEnvVars(cmd *cobra.Command, envVars []string) ([]nvcf.FunctionNewParamsContainerEnvironment, error) {
	var containerEnv []nvcf.FunctionNewParamsContainerEnvironment

	for _, env := range envVars {
		parts := strings.SplitN(env, ":", 2)
		if len(parts) != 2 {
			return nil, output.Error(cmd, fmt.Sprintf("invalid environment variable format: %s. ensure that you are using the format key:value", env), nil)
		}
		containerEnv = append(containerEnv, nvcf.FunctionNewParamsContainerEnvironment{
			Key:   nvcf.F(parts[0]),
			Value: nvcf.F(parts[1]),
		})
	}

	return containerEnv, nil
}

func parseEnvVarsNewVersion(cmd *cobra.Command, envVars []string) ([]nvcf.FunctionVersionNewParamsContainerEnvironment, error) {
	var containerEnv []nvcf.FunctionVersionNewParamsContainerEnvironment

	for _, env := range envVars {
		parts := strings.SplitN(env, ":", 2)
		if len(parts) != 2 {
			return nil, output.Error(cmd, fmt.Sprintf("invalid environment variable format: %s. ensure that you are using the format key:value", env), nil)
		}
		containerEnv = append(containerEnv, nvcf.FunctionVersionNewParamsContainerEnvironment{
			Key:   nvcf.F(parts[0]),
			Value: nvcf.F(parts[1]),
		})
	}

	return containerEnv, nil
}

func parseModels(cmd *cobra.Command, modelVars []string) ([]nvcf.FunctionNewParamsModel, error) {
	var models []nvcf.FunctionNewParamsModel

	for _, model := range modelVars {
		parts := strings.Split(model, ":")
		if len(parts) != 3 {
			return nil, output.Error(cmd, fmt.Sprintf("invalid model format: %s", model), nil)
		}
		models = append(models, nvcf.FunctionNewParamsModel{
			Name:    nvcf.F(parts[0]),
			Uri:     nvcf.F(parts[1]),
			Version: nvcf.F(parts[2]),
		})
	}

	return models, nil
}

func parseModelsNewVersion(cmd *cobra.Command, modelVars []string) ([]nvcf.FunctionVersionNewParamsModel, error) {
	var models []nvcf.FunctionVersionNewParamsModel

	for _, model := range modelVars {
		parts := strings.Split(model, ":")
		if len(parts) != 3 {
			return nil, output.Error(cmd, fmt.Sprintf("invalid model format: %s", model), nil)
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

	deployResp, err := client.FunctionDeployment.Functions.Versions.InitiateDeployment(
		cmd.Context(),
		resp.Function.ID,
		resp.Function.VersionID,
		deploymentParams,
	)
	if err != nil {
		return output.Error(cmd, "error deploying function", err)
	}

	output.Success(cmd, fmt.Sprintf("Function with FunctionID %s and VersionID %s deployed successfully", resp.Function.ID, resp.Function.VersionID))

	detatched, _ := cmd.Flags().GetBool("detatched")
	if !detatched {
		return WaitForDeployment(cmd, client, deployResp.Deployment.FunctionID, deployResp.Deployment.FunctionVersionID)
	}
	return nil
}

func createFunctionsFromFile(cmd *cobra.Command, client *api.Client, yamlFile string, deploy bool) error {
	data, err := os.ReadFile(yamlFile)
	if err != nil {
		return output.Error(cmd, "error reading YAML file", err)
	}

	var spec FunctionSpec
	if err := yaml.Unmarshal(data, &spec); err != nil {
		return output.Error(cmd, "error parsing YAML file", err)
	}

	for _, fn := range spec.Functions {
		if fn.ExistingFunctionID != "" {
			params := prepareFunctionVersionParamsFromFile(spec.FnImage, fn)
			if err := createAndDeployFunctionVersionFromFile(cmd, client, fn.ExistingFunctionID, params, deploy, fn.InstGPUType, fn.InstType, fn.InstBackend, fn.InstMax, fn.InstMin, fn.InstMaxRequestConcurrency); err != nil {
				return err
			}
		} else {
			params := prepareFunctionParamsFromFile(spec.FnImage, fn)
			if err := createAndDeployFunctionFromFile(cmd, client, params, deploy, fn.InstGPUType, fn.InstType, fn.InstBackend, fn.InstMax, fn.InstMin, fn.InstMaxRequestConcurrency); err != nil {
				return err
			}
		}
	}

	return nil
}

func prepareFunctionVersionParamsFromFile(fnImage string, fn FunctionDef) nvcf.FunctionVersionNewParams {
	apiBodyFormat := defaultAPIBodyFormat
	if !fn.Custom {
		apiBodyFormat = "PREDICT_V2"
	}

	functionType := defaultFunctionType
	if !fn.Streaming {
		functionType = "DEFAULT"
	}

	return nvcf.FunctionVersionNewParams{
		Name:           nvcf.String(fn.FnName),
		InferenceURL:   nvcf.String(fn.InferenceURL),
		InferencePort:  nvcf.Int(fn.InferencePort),
		ContainerImage: nvcf.String(fnImage),
		ContainerArgs:  nvcf.String(fn.ContainerArgs),
		APIBodyFormat:  nvcf.F(nvcf.FunctionVersionNewParamsAPIBodyFormat(apiBodyFormat)),
		Description:    nvcf.F(fn.Description),
		Tags:           nvcf.F(fn.Tags),
		FunctionType:   nvcf.F(nvcf.FunctionVersionNewParamsFunctionType(functionType)),
		Health: nvcf.F(nvcf.FunctionVersionNewParamsHealth{
			Protocol:           nvcf.F(nvcf.FunctionVersionNewParamsHealthProtocol(fn.Health.Protocol)),
			Port:               nvcf.F(fn.Health.Port),
			Timeout:            nvcf.F(flagutil.DurationToISO8601(fn.Health.Timeout)),
			ExpectedStatusCode: nvcf.F(fn.Health.ExpectedStatusCode),
			Uri:                nvcf.String(fn.Health.Uri),
		}),
	}
}

func createAndDeployFunctionVersionFromFile(cmd *cobra.Command, client *api.Client, existingFunctionID string, params nvcf.FunctionVersionNewParams, deploy bool, gpu, instanceType, backend string, maxInstances, minInstances, maxRequestConcurrency int64) error {
	resp, err := client.Functions.Versions.New(cmd.Context(), existingFunctionID, params)
	if err != nil {
		output.Error(cmd, "error creating function version", err)
		return nil
	}

	output.Success(cmd, fmt.Sprintf("Function version %s created successfully for function %s", resp.Function.VersionID, existingFunctionID))

	if deploy {
		return deployFunction(cmd, client, resp, gpu, instanceType, backend, maxInstances, minInstances, maxRequestConcurrency)
	}

	return nil
}

func prepareFunctionParamsFromFile(fnImage string, fn FunctionDef) nvcf.FunctionNewParams {
	apiBodyFormat := defaultAPIBodyFormat
	if !fn.Custom {
		apiBodyFormat = "PREDICT_V2"
	}

	functionType := defaultFunctionType
	if !fn.Streaming {
		functionType = "DEFAULT"
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
			Timeout:            nvcf.F(flagutil.DurationToISO8601(fn.Health.Timeout)),
			ExpectedStatusCode: nvcf.F(fn.Health.ExpectedStatusCode),
			Uri:                nvcf.String(fn.Health.Uri),
		}),
	}
}

func createAndDeployFunctionFromFile(cmd *cobra.Command, client *api.Client, params nvcf.FunctionNewParams, deploy bool, gpu, instanceType, backend string, maxInstances, minInstances, maxRequestConcurrency int64) error {
	resp, err := client.Functions.New(cmd.Context(), params)
	if err != nil {
		return output.Error(cmd, "error creating function", err)
	}

	if deploy {
		return deployFunction(cmd, client, resp, gpu, instanceType, backend, maxInstances, minInstances, maxRequestConcurrency)
	}

	output.Success(cmd, fmt.Sprintf("Function %s with id %s and version %s created successfully", params.Name, resp.Function.ID, resp.Function.VersionID))
	return nil
}

func WaitForDeployment(cmd *cobra.Command, client *api.Client, functionID, versionID string) error {
	spinner := output.NewSpinner(fmt.Sprintf("Waiting for deployment of function %s with version %s to complete...", functionID, versionID))
	output.StartSpinner(spinner)
	defer output.StopSpinner(spinner)

	err := timeout.DoWithTimeout(func(ctx context.Context) error {
		for ctx.Err() == nil {
			deploymentStatus, err := client.FunctionDeployment.Functions.Versions.GetDeployment(
				ctx,
				functionID,
				versionID,
			)
			if err != nil {
				return err
			}
			if deploymentStatus.Deployment.FunctionStatus == nvcf.DeploymentResponseDeploymentFunctionStatusActive {
				return nil
			}
			if deploymentStatus.Deployment.FunctionStatus == nvcf.DeploymentResponseDeploymentFunctionStatusError {
				return fmt.Errorf("deployment failed. please try again")
			}
			time.Sleep(time.Second * 5)
		}
		return nil
	}, 30*time.Minute) // GFN can take a while

	if err != nil {
		return output.Error(cmd, "Error waiting for deployment", err)
	}

	output.Success(cmd, fmt.Sprintf("Function with FunctionID %s and VersionID %s deployed successfully", functionID, versionID))
	return nil
}

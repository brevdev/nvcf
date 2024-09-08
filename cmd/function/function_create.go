package function

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/spf13/cobra"
	"github.com/tmc/nvcf-go"
	"github.com/tmc/nvcf/api"
	"github.com/tmc/nvcf/config"
	"github.com/tmc/nvcf/output"
)

func functionCreateCmd() *cobra.Command {
	var (
		name             string
		inferenceURL     string
		containerImage   string
		description      string
		tags             []string
		minInstances     int64
		maxInstances     int64
		gpu              string
		instanceType     string
		apiBodyFormat    string
		functionType     string
		healthProtocol   string
		healthPort       int64
		healthTimeout    string
		healthStatusCode int64
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new function",
		Long:  `Create a new NVIDIA Cloud Function with the specified parameters.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client := api.NewClient(config.GetAPIKey())

			params := nvcf.FunctionNewParams{
				Name:           nvcf.F(name),
				InferenceURL:   nvcf.F(inferenceURL),
				ContainerImage: nvcf.F(containerImage),
				Description:    nvcf.F(description),
				Tags:           nvcf.F(tags),
				APIBodyFormat:  nvcf.F(nvcf.FunctionNewParamsAPIBodyFormat(apiBodyFormat)),
				FunctionType:   nvcf.F(nvcf.FunctionNewParamsFunctionType(functionType)),
				Health: nvcf.F(nvcf.FunctionNewParamsHealth{
					Protocol:           nvcf.F(nvcf.FunctionNewParamsHealthProtocol(healthProtocol)),
					Port:               nvcf.F(healthPort),
					Timeout:            nvcf.F(healthTimeout),
					ExpectedStatusCode: nvcf.F(healthStatusCode),
				}),
			}

			// Remove the DeploymentSpecifications field
			// Instead, we'll use the minInstances, maxInstances, gpu, and instanceType
			// to configure the function after creation if needed

			resp, err := client.Functions.New(cmd.Context(), params)
			if err != nil {
				return err
			}

			var fn nvcf.ListFunctionsResponseFunction
			if err := jsonMarshalUnmarshal(&fn, resp.Function); err != nil {
				return fmt.Errorf("issue marshaling+unmarshaling: %w", err)
			}
			output.Function(cmd, fn)
			return nil
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Name of the function (required)")
	cmd.Flags().StringVar(&inferenceURL, "inference-url", "", "URL for function invocation (required)")
	cmd.Flags().StringVar(&containerImage, "container-image", "", "Container image for the function")
	cmd.Flags().StringVar(&description, "description", "", "Description of the function")
	cmd.Flags().StringSliceVar(&tags, "tag", nil, "Tags for the function (can be used multiple times)")
	cmd.Flags().Int64Var(&minInstances, "min-instances", 0, "Minimum number of instances")
	cmd.Flags().Int64Var(&maxInstances, "max-instances", 0, "Maximum number of instances")
	cmd.Flags().StringVar(&gpu, "gpu", "", "GPU type to use")
	cmd.Flags().StringVar(&instanceType, "instance-type", "", "Instance type to use")
	cmd.Flags().StringVar(&apiBodyFormat, "api-body-format", "PREDICT_V2", "API body format (PREDICT_V2 or CUSTOM)")
	cmd.Flags().StringVar(&functionType, "function-type", "DEFAULT", "Function type (DEFAULT or STREAMING)")
	cmd.Flags().StringVar(&healthProtocol, "health-protocol", "HTTP", "Health check protocol (HTTP or GRPC)")
	cmd.Flags().Int64Var(&healthPort, "health-port", 8080, "Health check port")
	cmd.Flags().StringVar(&healthTimeout, "health-timeout", "5s", "Health check timeout")
	cmd.Flags().Int64Var(&healthStatusCode, "health-status-code", 200, "Expected health check status code")
	// todo: handle this correctly
	cmd.Flags().Bool("deploy", false, "Create and deploy the function in one step")

	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("inference-url")

	return cmd
}

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

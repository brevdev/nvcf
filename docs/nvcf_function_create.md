## nvcf function create

Create a new function

### Synopsis

Create a new NVCF Function with the specified parameters. If you specify --from-version, we will create a new version of an existing function. You can also create and deploy a function in one step using the --deploy flag.

```
nvcf function create [flags]
```

### Examples

```
Create a new function:
nvcf function create --name myfunction --inference-url /v1/chat/completions --inference-port 80 --health-uri /health --container-image nvcr.io/nvidia/example-image:latest
   
Create and deploy a new function from a file:
nvcf function create --file deploy.yaml --deploy

Create a new version of an existing function:
nvcf function create --from-version existing-function-id --name newversion --inference-url /v2/chat/completions --inference-port 8080 --health-uri /healthcheck --container-image nvcr.io/nvidia/updated-image:v2

```

### Options

```
      --backend string                Backend to deploy the function to (see your NGC org available backends)
      --container-args string         Container arguments. Put these in quotes if you are passing flags
      --container-image string        Container image for the function
      --custom                        Set API body format to CUSTOM. If false - set API body format to PREDICT_V2. (default true)
      --deploy                        Create and deploy the function in one step. Default is false
      --description string            Description of the function
  -d, --detatched                     Deploy the function in the background. Default is false
      --env strings                   Environment variables for the function (can be used multiple times, format: key:value)
  -f, --file string                   Path to a YAML file containing function specifications
      --from-version string           Create a new version of an existing function. Requires a valid function id
      --function-type string          Function type (DEFAULT or STREAMING). Default is DEFAULT (default "STREAMING")
      --gpu string                    GPU type to use
      --health-port int               Health check port. Default is 80 (default 80)
      --health-protocol string        Health check protocol (HTTP or GRPC). Default is HTTP (default "HTTP")
      --health-status-code int        Expected health check status code. Default is 200 (default 200)
      --health-timeout duration       Health check timeout. (default 20s)
      --health-uri string             Health check URI. Default is /health (default "/health")
  -h, --help                          help for create
      --inference-port int            Port for function invocation. Default is 80 (default 80)
      --inference-url string          URL for function invocation (required)
      --instance-type string          Instance type to use. Default is GCP.GPU.H100_1x
      --max-instances int             Maximum number of instances. Default is 1 (default 1)
      --max-request-concurrency int   Maximum number of concurrent requests. Default is 1 (default 1)
      --min-instances int             Minimum number of instances. Default is 0
      --model strings                 Models for the function (can be used multiple times, format: name:uri:version)
      --name string                   Name of the function (required)
      --streaming                     Set function type to STREAMING. Default is true (default true)
      --tag strings                   Tags for the function (can be used multiple times)
```

### Options inherited from parent commands

```
      --json       Output results in JSON format
      --no-color   Disable color output
  -q, --quiet      Suppress non-error output
  -v, --verbose    Enable verbose output and show underlying API calls
```

### SEE ALSO

* [nvcf function](nvcf_function.md)	 - Manage NVIDIA Cloud Functions


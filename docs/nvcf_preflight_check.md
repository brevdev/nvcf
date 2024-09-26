## nvcf preflight check

Run local deployment test for NVCF compatibility

### Synopsis

Run a local deployment test to verify container compatibility with NVCF.

This command deploys the specified Docker image locally, checks its health,
and performs basic connectivity tests to ensure it meets NVCF requirements.

Key features:
- Supports both HTTP and gRPC protocols
- Customizable health check and inference endpoints
- Configurable wait times for container readiness
- Option to force cleanup of existing containers

Use this command to validate your NVCF function before deployment.

```
nvcf preflight check <image-name> [flags]
```

### Options

```
      --container-port string             Port that the server is listening on (default "8000")
      --force-cleanup                     Force cleanup of existing containers before starting the test
      --grpc-input-data string            JSON string representing input data for gRPC method (default "{}")
      --grpc-method-name string           gRPC method name (required for gRPC)
      --grpc-service-name string          gRPC service name (required for gRPC)
      --health-endpoint string            Health endpoint exposed by the server (for HTTP) (default "/v2/health/ready")
  -h, --help                              help for check
      --http-inference-endpoint string    HTTP inference endpoint (default "/")
      --http-payload string               JSON string representing input data for HTTP inference (default "{}")
      --protocol string                   Protocol that the server is running (http or grpc) (default "http")
      --seconds-to-wait-for-healthy int   Maximum time to wait for the health endpoint to be ready (in seconds) (default 600)
```

### Options inherited from parent commands

```
      --json       Output results in JSON format
      --no-color   Disable color output
  -q, --quiet      Suppress non-error output
  -v, --verbose    Enable verbose output and show underlying API calls
```

### SEE ALSO

* [nvcf preflight](nvcf_preflight.md)	 - Perform preflight checks for NVCF compatibility


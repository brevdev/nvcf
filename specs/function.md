# NVIDIA Cloud Functions (NVCF) Function Specification

This document outlines the structure and components of the YAML specification used to create and deploy functions via the NVCF CLI. Users can utilize this spec with the command `nvcf fn create -f <path-to-file>`. You can also auto deploy the function by using the `--deploy` flag.

## Overview

The specification defines the container image and one or more functions to be deployed. Each function can have its own configuration, including environment variables, instance types, associated models, and other parameters.

## Specification Structure

### Root Level

- `fn_image`: The container image URL from NVCR (NVIDIA Container Registry).
- `functions`: An array of function configurations based on the `fn_image`.

### Function Configuration

Each function in the `functions` array can have the following properties:

| Field | Description | Required | Default |
|-------|-------------|----------|---------|
| `fn_name` | The name of the function | Yes | N/A |
| `existingFunctionID` | ID of an existing function to create a new version | No | N/A |
| `inferenceUrl` | The entrypoint URL for invoking the container | Yes | N/A |
| `inferencePort` | The port number for the inference listener | No | 80 |
| `containerImage` | The container image for this function | No | Value of `fn_image` |
| `containerArgs` | Command-line arguments for the container | No | "" |
| `custom` | If true, sets API body format to CUSTOM; if false, sets to PREDICT_V2 | No | true |
| `description` | A description of the function | No | "" |
| `streaming` | If true, sets function type to STREAMING; if false, sets to DEFAULT | No | true |
| `tags` | An array of tags for the function | No | [] |
| `health` | Health check configuration | Yes | See below |
| `inst_backend` | The backend infrastructure | Yes | N/A |
| `inst_gpu_type` | The type of GPU required | Yes | N/A |
| `inst_type` | The instance type specification | Yes | N/A |
| `inst_min` | Minimum number of instances | No | 0 |
| `inst_max` | Maximum number of instances | No | 1 |
| `inst_max_request_concurrency` | Max concurrent requests per instance | No | 1 |
| `containerEnvironment` | Array of environment variables | No | [] |
| `models` | Array of models associated with the function | No | [] |

### Health Check Configuration

The `health` property must include:

| Field | Description | Required | Default |
|-------|-------------|----------|---------|
| `protocol` | The health check protocol (HTTP or gRPC) | No | "HTTP" |
| `port` | The port for health checks | No | 80 |
| `timeout` | The timeout for health checks | No | "PT20S" |
| `expectedStatusCode` | Expected status code for a successful check | No | 200 |
| `uri` | The health check endpoint URI | Yes | "/health" |

### Environment Variables

Each item in the `containerEnvironment` array should have:

| Field | Description | Required |
|-------|-------------|----------|
| `key` | The name of the environment variable | Yes |
| `value` | The value of the environment variable | Yes |

### Models

Each item in the `models` array should have:

| Field | Description | Required |
|-------|-------------|----------|
| `name` | The name of the model | Yes |
| `version` | The version of the model | Yes |
| `uri` | The URI of the model | Yes |

## Example Specification

```yaml
fn_image: nvcr.io/sklmhpjhptei/test-team/brev-tgi:2.2.0
functions:
  - fn_name: example-function
    existingFunctionID: "1234567890"
    inferenceUrl: "/v1/chat/completions"
    inferencePort: 80
    healthUri: "/health"
    containerArgs: "--model-id Qwen/Qwen2-7B-Instruct"
    custom: true
    description: "Example function for CLI fn testing"
    streaming: true
    tags:
      - "nlp"
      - "inference"
    health:
      protocol: HTTP
      port: 80
      timeout: "PT20S"
      expectedStatusCode: 200
      uri: "/health"
    inst_backend: gcp-asia-se-1a
    inst_gpu_type: H100
    inst_type: GCP.GPU.H100_1x
    inst_min: 1
    inst_max: 2
    inst_max_request_concurrency: 2
    containerEnvironment:
      - key: DEBUG
        value: "true"
      - key: LOG_LEVEL
        value: "info"
    models:
      - name: example-model-1
        version: v1.0
        uri: nvcr.io/nvidia/example-model-1:v1.0
      - name: example-model-2
        version: v2.0
        uri: nvcr.io/nvidia/example-model-2:v2.0
```

## Usage

To create a function using this specification:

```bash
nvcf fn create -f path/to/your/spec.yaml
```

To create and immediately deploy the function:

```bash
nvcf fn create -f path/to/your/spec.yaml --deploy
```

## Notes

1. The `existingFunctionID` field is used to create a new version of an existing function. If provided, the CLI will create a new version for the specified function instead of creating a new function.
2. The `custom` field replaces the previous `apiBodyFormat` field. If `custom` is true, the API body format is set to CUSTOM; if false, it's set to PREDICT_V2.
3. The `streaming` field replaces the previous `functionType` field. If `streaming` is true, the function type is set to STREAMING; if false, it's set to DEFAULT.
4. Resources and secrets are not currently supported in the implementation and thus not included in this specification.
5. HelmChart and HelmChartServiceName functionality is not currently implemented.

This specification provides a way to define and deploy NVIDIA Cloud Functions, allowing for configuration of multiple functions with varying parameters, including environment variables and associated models. It reflects the current implementation capabilities.

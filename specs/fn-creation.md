# NVIDIA Cloud Functions (NVCF) Specification

This document outlines the structure and components of the YAML specification used to create and deploy functions via the NVCF CLI. Users can utilize this spec with the command `nvcf fn create -f <path-to-file>`.

## Overview

The specification defines the container image and one or more functions to be deployed. Each function can have its own configuration, including environment variables, instance types, and associated models.

## Specification Structure

### Root Level

- `fn_image`: The container image URL from NVCR (NVIDIA Container Registry).
- `functions`: An array of function configurations.

### Function Configuration

Each function in the `functions` array can have the following properties:

| Field | Description | Required | Default |
|-------|-------------|----------|---------|
| `fn_name` | The name of the function | Yes | N/A |
| `inferenceUrl` | The entrypoint URL for invoking the container | Yes | N/A |
| `inferencePort` | The port number for the inference listener | Yes | 80 |
| `healthUri` | The health check endpoint URI | Yes | "/health" |
| `containerImage` | The container image for this function | Yes | N/A |
| `containerArgs` | Command-line arguments for the container | No | "" |
| `apiBodyFormat` | The format of the API request body (PREDICT_V2 or CUSTOM) | No | "CUSTOM" |
| `description` | A description of the function | No | "" |
| `functionType` | The type of function (DEFAULT or STREAMING) | No | "DEFAULT" |
| `tags` | An array of tags for the function | No | [] |
| `health` | Health check configuration | Yes | See below |
| `inst_backend` | The backend infrastructure | Yes | N/A |
| `inst_gpu_type` | The type of GPU required | Yes | N/A |
| `inst_type` | The instance type specification | Yes | N/A |
| `inst_min` | Minimum number of instances | Yes | 0 |
| `inst_max` | Maximum number of instances | Yes | 1 |
| `inst_max_request_concurrency` | Max concurrent requests per instance | Yes | 1 |

### Health Check Configuration

The `health` property must include:

| Field | Description | Required | Default |
|-------|-------------|----------|---------|
| `protocol` | The health check protocol (HTTP or gRPC) | Yes | "HTTP" |
| `port` | The port for health checks | Yes | 80 |
| `timeout` | The timeout for health checks | Yes | "PT20S" |
| `expectedStatusCode` | Expected status code for a successful check | Yes | 200 |
| `uri` | The health check endpoint URI | Yes | "/health" |

## Example Specification

```yaml
fn_image: nvcr.io/nvidia/example-image:v1.0
functions:
  - fn_name: example-function
    inferenceUrl: "/v1/chat/completions"
    inferencePort: 80
    healthUri: "/health"
    containerImage: nvcr.io/nvidia/example-image:v1.0
    containerArgs: "--model-id example-model"
    apiBodyFormat: CUSTOM
    description: "Example function for NLP tasks"
    functionType: DEFAULT
    tags:
      - "nlp"
      - "inference"
    health:
      protocol: HTTP
      port: 80
      timeout: "PT20S"
      expectedStatusCode: 200
      uri: "/health"
    inst_backend: GFN
    inst_gpu_type: A100
    inst_type: gn1.a100.2x.30gb
    inst_min: 1
    inst_max: 3
    inst_max_request_concurrency: 2
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

1. Environment variables, models, resources, and secrets are not currently supported in the implementation and thus not included in this specification.
2. HelmChart and HelmChartServiceName functionality is not currently implemented.

This specification provides a way to define and deploy NVIDIA Cloud Functions, allowing for configuration of multiple functions with varying parameters. It reflects the current implementation capabilities.

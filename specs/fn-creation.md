# Function Specification for NVCF CLI

This document outlines the structure and components of the YAML/JSON specification used to create and deploy functions via the NVCF CLI. Users can utilize this spec with the command `nvcf fn create -f ./<path-to-file>`. You can also pass in an optional --deploy flag to automatically deploy the function after creation.

## Overview

The specification defines the container image and one or more functions to be deployed. Each function can have its own configuration, including environment variables, instance types, and associated models.

## Specification Structure

### Root Level

- `fn_image`: The container image URL from NVCR (NVIDIA Container Registry).

- `functions`: An array of function configurations.

### Function Configuration

Each function in the `functions` array can have the following properties:

- `containerArgs`: Command-line arguments passed to the container.
- `env`: An array of key-value pairs for environment variables.
- `fn_name`: The name of the function (must be unique).
- `inst_backend`: The backend infrastructure (e.g., "GFN" for GPU Flexible Nodes).
- `inst_gpu_type`: The type of GPU required (e.g., "L40S").
- `inst_max`: Maximum number of instances.
- `inst_max_request_concurrency`: Maximum concurrent requests per instance.
- `inst_min`: Minimum number of instances.
- `inst_type`: The instance type specification.
- `models`: An array of model configurations associated with the function.

### Model Configuration

Each model in the `models` array can have:

- `name`: The name of the model.
- `uri`: The URI or identifier for the model.
- `version`: The version of the model.

If referencing any models (i.e. for volume mounts) ensure these models exist and have been uploaded to NGC via:

```bash
ngc registry model --org ${FN_NGC_ORG} create --application Other --framework Other --precision Other --format Other --short-desc ${FN_NGC_MODEL} ${FN_NGC_ORG}/${FN_NGC_TEAM}/${FN_NGC_MODEL}

cd your/model/dir

ngc registry model --org ${FN_NGC_ORG} --team ${FN_NGC_TEAM} upload-version ${FN_NGC_ORG}/${FN_NGC_TEAM}/${FN_NGC_MODEL}:${FN_NGC_MODEL_VERSION}
```

## Example Specification

```yaml
fn_image: nvcr.io/sklmhpjhptei/test-team/test-tgi:v1.0.1
functions:
  - containerArgs: --model-id meta-llama/Meta-Llama-3-8B
    env:
      - key: FN_SAMPLE_ENV_KEY
        value: FN_SAMPLE_ENV_VALUE
    fn_name: inference-l40sx1
    inst_backend: GFN
    inst_gpu_type: L40S
    inst_max: 1
    inst_max_request_concurrency: 1
    inst_min: 1
    inst_type: gl40s_1.br25_2xlarge
    models:
      - name: sample-model
        uri: sample-model
        version: 0.1
  - containerArgs: --model-id meta-llama/Meta-Llama-3-8B
    env:
      - key: FN_SAMPLE_ENV_KEY
        value: FN_SAMPLE_ENV_VALUE
    fn_name: inference-l40sx1-v2
    inst_backend: GFN
    inst_gpu_type: L40S
    inst_max: 1
    inst_max_request_concurrency: 1
    inst_min: 1
    inst_type: gl40s_1.br25_2xlarge
    models:
      - name: sample-model
        uri: sample-model
        version: 0.1
```


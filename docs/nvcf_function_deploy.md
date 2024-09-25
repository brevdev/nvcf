## nvcf function deploy

Deploy a function

### Synopsis

Deploy an existing NVCF function. If you want to deploy a specific version, use the --version-id flag.

```
nvcf function deploy <function-id> [flags]
```

### Examples

```
nvcf function deploy fid --version-id vid --gpu A100 --instance-type g5.4xlarge
```

### Options

```
      --backend string                Backend to deploy the function to
  -d, --detached                      Detach from the deployment and return to the prompt
      --gpu string                    GPU type to use
  -h, --help                          help for deploy
      --instance-type string          Instance type to use
      --max-instances int             Maximum number of instances (default 1)
      --max-request-concurrency int   Maximum number of concurrent requests (default 1)
      --min-instances int             Minimum number of instances
      --version-id string             The ID of the version to deploy
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


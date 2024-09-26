## nvcf function update

Update a deployed function

### Synopsis

If a version-id is not provided, we look for versions that are actively deployed. If a single function is deployed, we update that version. If multiple functions are deployed, we prompt for the version-id to update.

```
nvcf function update <function-id> [flags]
```

### Examples

```
nvcf function update fid --version-id vid --gpu A100 --instance-type g5.4xlarge --min-instances 1 --max-instances 5 --max-request-concurrency 100
```

### Options

```
      --gpu string                    GPU name from the cluster
  -h, --help                          help for update
      --instance-type string          Instance type, based on GPU, assigned to a Worker
      --max-instances int             Maximum number of spot instances for the deployment
      --max-request-concurrency int   Max request concurrency between 1 (default) and 1024
      --min-instances int             Minimum number of spot instances for the deployment
      --version-id string             The ID of the version
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


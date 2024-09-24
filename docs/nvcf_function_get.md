## nvcf function get

Get details about a single function and its versions

### Synopsis

Get details about a single function and its versions or deployments. If a version-id is not provided and there are multiple versions associated with a function, we will look for all versions and prompt for a version-id.

```
nvcf function get <function-id> [flags]
```

### Examples

```
nvcf function get fid --version-id vid --include-secrets
```

### Options

```
  -h, --help                help for get
      --include-secrets     Include secrets in the response
      --version-id string   The ID of the version
```

### Options inherited from parent commands

```
      --json       Output results in JSON format
      --no-color   Disable color output
  -q, --quiet      Suppress non-error output
  -v, --verbose    Enable verbose output
```

### SEE ALSO

* [nvcf function](nvcf_function.md)	 - Manage NVIDIA Cloud Functions


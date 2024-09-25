## nvcf function get

Get details about a single function and its versions

### Synopsis

Get details about a single function and its versions or deployments. The identifier can be a function name, function ID, or version ID. If no identifier is provided, all functions will be listed.

```
nvcf function get [identifier] [flags]
```

### Examples

```
nvcf function get myFunction
nvcf function get fid123
nvcf function get --name myFunction
nvcf function get --function-id fid123 --version-id vid456
```

### Options

```
      --function-id string   Filter by function ID
  -h, --help                 help for get
      --include-secrets      Include secrets in the response
      --name string          Filter by function name
      --version-id string    Filter by version ID
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


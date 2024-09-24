## nvcf function stop

Stop a deployed function

### Synopsis

Stop a deployed function. If a version-id is not provided, we look for versions that are actively deployed. If a single function is deployed, we stop that version. If multiple functions are deployed, we prompt for the version-id to stop.

```
nvcf function stop <function-id> [flags]
```

### Examples

```
nvcf function stop fid --version-id vid
```

### Options

```
      --all                 Stop all deployed versions of the function
      --force               Gracefully stop the function if it's already deployed. If not, we forcefully stop the function.
  -h, --help                help for stop
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


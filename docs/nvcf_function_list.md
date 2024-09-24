## nvcf function list

List all functions. Use flags to filter by visibility and status.

### Synopsis

List all functions. Use the --visibility flag to filter by visibility and --status flag to filter by status. Note that deployments are functions that are NOT inactive.

```
nvcf function list [flags]
```

### Examples

```
nvcf function list --visibility authorized --status ACTIVE,DEPLOYING
```

### Options

```
  -h, --help                 help for list
      --status strings       Filter by status (ACTIVE, DEPLOYING, ERROR, INACTIVE, DELETED). Defaults to all. (default [ACTIVE,DEPLOYING,ERROR,INACTIVE,DELETED])
      --visibility strings   Filter by visibility (authorized, private, public). Defaults to private. (default [private])
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


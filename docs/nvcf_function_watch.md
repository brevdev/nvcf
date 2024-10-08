## nvcf function watch

Watch functions status in real-time

### Synopsis

Display a real-time view of NVCF functions.

```
nvcf function watch [flags]
```

### Options

```
  -h, --help                 help for watch
      --status strings       Filter by status (ACTIVE, DEPLOYING, ERROR, INACTIVE, DELETED). Defaults to all. (default [ACTIVE,DEPLOYING,ERROR,INACTIVE,DELETED])
      --visibility strings   Filter by visibility (authorized, private, public). Defaults to private. (default [private])
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


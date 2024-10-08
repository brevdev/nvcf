## nvcf function delete

Delete a function. If you want to delete a specific version, use the --version-id flag.

### Synopsis

Delete a function. If there is only 1 version, we will delete the function. If there are multiple versions, we will prompt you to specify which version to delete. The --all flag will delete all versions of the function. Deleting a deployed function will change a function status to INACTIVE and using the --force flag will delete the function immediately.

```
nvcf function delete <function-id> [flags]
```

### Examples

```
nvcf function delete fid --version-id vid
```

### Options

```
      --all                 Delete all versions of the function
      --force               Forcefully delete a deployed function
  -h, --help                help for delete
      --version-id string   The ID of the version
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


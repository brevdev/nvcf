# NVIDIA Cloud Functions CLI

## Table of Contents

- [Overview](#overview)
- [Installation](#installation)
- [Authentication](#authentication)
- [Global Options](#global-options)
- [Commands](#commands)
  - [function](#function)
  - [deployment](#deployment)
  - [invoke](#invoke)
  - [asset](#asset)
  - [auth](#auth)
  - [queue](#queue)
  - [cluster-group](#cluster-group)
- [Examples](#examples)
- [Getting Help](#getting-help)

## Overview

The NVIDIA Cloud Functions CLI (`nvcf`) is a command-line interface for managing and interacting with NVIDIA Cloud Functions. It allows you to create, deploy, invoke, and manage serverless functions on NVIDIA's cloud infrastructure.

## Installation

To install the NVIDIA Cloud Functions CLI, run:

```
go install github.com/tmc/nvcf@latest
```

## Authentication

Before using the CLI, you need to authenticate. You can do this by running:

```
nvcf auth login
```

This will prompt you to enter your NVIDIA Cloud credentials or API key.

## Global Options

The following options can be used with any command:

- `--help, -h`: Show help for the command
- `--json`: Output results in JSON format
- `--no-color`: Disable color output
- `--quiet, -q`: Suppress non-error output

## Commands

### function

Manage NVIDIA Cloud Functions.

#### Usage

```
nvcf function <subcommand> [flags]
```

#### Subcommands

- `list`: List all functions
- `create`: Create a new function
- `get`: Get details of a function
- `update`: Update a function
- `delete`: Delete a function
- `version`: Manage function versions

##### function list

List all functions in your account.

```
nvcf function list [flags]
```

Flags:
- `--limit <number>`: Maximum number of functions to list
- `--visibility <string>`: Filter by visibility (authorized, private, public)

##### function create

Create a new function.

```
nvcf function create [flags]
```

Flags:
- `--name <string>`: Name of the function (required)
- `--inference-url <string>`: URL for function invocation (required)
- `--container-image <string>`: Container image for the function
- `--description <string>`: Description of the function
- `--tag <string>`: Tags for the function (can be used multiple times)

##### function get

Get details of a specific function.

```
nvcf function get <function-id> [flags]
```

##### function update

Update an existing function.

```
nvcf function update <function-id> [flags]
```

Flags:
- `--name <string>`: New name for the function
- `--inference-url <string>`: New URL for function invocation
- `--description <string>`: New description for the function
- `--tag <string>`: New tags for the function (can be used multiple times)

##### function delete

Delete a function.

```
nvcf function delete <function-id> [flags]
```

Flags:
- `--force`: Force deletion without confirmation

##### function version

Manage function versions.

```
nvcf function version <subcommand> [flags]
```

Subcommands:
- `list`: List versions of a function
- `create`: Create a new version of a function
- `get`: Get details of a specific function version
- `delete`: Delete a function version

### deployment

Manage function deployments.

#### Usage

```
nvcf deployment <subcommand> [flags]
```

#### Subcommands

- `create`: Deploy a function
- `get`: Get deployment details
- `update`: Update a deployment
- `delete`: Delete a deployment

##### deployment create

Deploy a function.

```
nvcf deployment create <function-id> <version-id> [flags]
```

Flags:
- `--min-instances <number>`: Minimum number of instances
- `--max-instances <number>`: Maximum number of instances
- `--gpu <string>`: GPU type to use
- `--instance-type <string>`: Instance type to use

##### deployment get

Get deployment details.

```
nvcf deployment get <function-id> <version-id> [flags]
```

##### deployment update

Update a deployment.

```
nvcf deployment update <function-id> <version-id> [flags]
```

Flags:
- `--min-instances <number>`: New minimum number of instances
- `--max-instances <number>`: New maximum number of instances

##### deployment delete

Delete a deployment.

```
nvcf deployment delete <function-id> <version-id> [flags]
```

Flags:
- `--force`: Force deletion without confirmation
- `--graceful`: Perform a graceful shutdown

### invoke

Invoke a function.

#### Usage

```
nvcf invoke <function-id> [version-id] [flags]
```

Flags:
- `--data <string>`: JSON string to pass as input data
- `--data-file <file>`: File containing JSON input data
- `--async`: Invoke function asynchronously
- `--timeout <duration>`: Timeout for the function invocation

### asset

Manage assets for functions.

#### Usage

```
nvcf asset <subcommand> [flags]
```

#### Subcommands

- `list`: List all assets
- `create`: Create a new asset
- `get`: Get details of an asset
- `delete`: Delete an asset

##### asset create

Create a new asset.

```
nvcf asset create [flags]
```

Flags:
- `--file <file>`: File to upload as an asset
- `--description <string>`: Description of the asset
- `--content-type <string>`: Content type of the asset

### auth

Manage authentication for the CLI.

#### Usage

```
nvcf auth <subcommand> [flags]
```

#### Subcommands

- `login`: Authenticate with NVIDIA Cloud
- `logout`: Log out and remove stored credentials
- `status`: Show the current authentication status

### queue

Manage and view function queues.

#### Usage

```
nvcf queue <subcommand> [flags]
```

#### Subcommands

- `list`: List queues for a function
- `position`: Get position of a request in the queue

##### queue list

List queues for a function.

```
nvcf queue list <function-id> [version-id] [flags]
```

##### queue position

Get position of a request in the queue.

```
nvcf queue position <request-id> [flags]
```

### cluster-group

Manage cluster groups.

#### Usage

```
nvcf cluster-group <subcommand> [flags]
```

#### Subcommands

- `list`: List available cluster groups

##### cluster-group list

List available cluster groups.

```
nvcf cluster-group list [flags]
```

## Examples

1. Create a new function:

```
nvcf function create --name my-function --inference-url https://example.com/function --container-image nvidia/my-function:latest
```

2. Deploy a function:

```
nvcf deployment create abc123 def456 --min-instances 1 --max-instances 5 --gpu A100 --instance-type g4dn.xlarge
```

3. Invoke a function:

```
nvcf invoke abc123 --data '{"input": "Hello, World!"}'
```

4. List all functions:

```
nvcf function list
```

5. Update a deployment:

```
nvcf deployment update abc123 def456 --min-instances 2 --max-instances 10
```

## Getting Help

To get help for any command, you can use the `--help` flag. For example:

```
nvcf --help
nvcf function --help
nvcf function create --help
```

For more detailed documentation, visit the [NVIDIA Cloud Functions documentation](https://docs.nvidia.com/cloud-functions/).

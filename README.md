# NVCF - NVIDIA Cloud Functions CLI

NVCF is a command-line interface (CLI) tool for managing and interacting with NVIDIA Cloud Functions. It provides a seamless way to create, deploy, invoke, and manage serverless functions on NVIDIA's cloud infrastructure.

## Key Features

- Function management (create, list, update, delete)
- Function deployment and invocation
- Asset management
- Authentication handling
- Queue management
- Cluster group operations
- Comprehensive error handling and logging
- Multiple output formats (JSON, table)
- Color-coded output for better readability

## Installation

To install the NVCF CLI, you need to have Go installed on your system. Then, you can use the following command:

```bash
go install github.com/brevdev/nvcf@latest
```

## Usage

After installation, you can use the `nvcf` command to interact with NVIDIA Cloud Functions. Here are some common usage examples:

```bash
# Authenticate with NVIDIA Cloud
nvcf auth login

# List all functions
nvcf function list

# Create a new function
nvcf function create --name my-function --inference-url https://example.com/function

# Invoke a function
nvcf invoke <function-id> --data '{"input": "Hello, World!"}'

# Get deployment details
nvcf deployment get <function-id> <version-id>
```

For a full list of commands and options, use the `--help` flag:

```bash
nvcf --help
```

## Project Structure

The project is organized as follows:

- `cmd/`: Contains the implementation of CLI commands
- `api/`: Implements the API client for interacting with NVIDIA Cloud Functions
- `config/`: Handles configuration management
- `output/`: Manages output formatting and display
- `main.go`: Entry point of the CLI application

## Dependencies

The project uses the following main dependencies:

- `github.com/spf13/cobra`: For building the CLI command structure
- `github.com/brevdev/nvcf-go`: NVIDIA Cloud Functions Go SDK
- `github.com/fatih/color`: For color-coded output
- `github.com/olekukonko/tablewriter`: For table-formatted output

## Contributing

Contributions to the NVCF CLI are welcome. Please follow these steps to contribute:

1. Fork the repository
2. Create a new branch for your feature or bug fix
3. Make your changes and commit them with clear, descriptive messages
4. Push your changes to your fork
5. Submit a pull request to the main repository

Please ensure your code adheres to the existing style and includes appropriate tests and documentation.

## License

This project is licensed under the [MIT License](LICENSE).

## Badges

[![Go Report Card](https://goreportcard.com/badge/github.com/brevdev/nvcf)](https://goreportcard.com/report/github.com/brevdev/nvcf)
[![GoDoc](https://godoc.org/github.com/brevdev/nvcf?status.svg)](https://godoc.org/github.com/brevdev/nvcf)

## Support

For issues, feature requests, or questions, please open an issue on the [GitHub repository](https://github.com/brevdev/nvcf/issues).

For more information about NVIDIA Cloud Functions, visit the [official documentation](https://docs.nvidia.com/cloud-functions/)

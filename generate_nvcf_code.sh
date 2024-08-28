#!/usr/bin/env bash

set -euo pipefail

# Function to display script usage
usage() {
    echo "Usage: $0 <command_name>"
    echo "Generates Go code for the specified nvcf command using cgpt."
    exit 1
}

# Function to log messages
log() {
    echo "[$(date +'%Y-%m-%d %H:%M:%S')] $1" >&2
}

# Function to generate a system prompt for cgpt
generate_system_prompt() {
    cat <<EOF
You are an expert Go developer specializing in CLI tools and SDK libraries. Your task is to generate Go code for the nvcf codebase. Follow these guidelines:
1. Use idiomatic Go patterns and best practices.
2. Implement robust error handling and logging.
3. Add comprehensive comments and documentation.
4. Ensure consistency with the existing codebase structure.
5. Consider adding unit tests for new functionality.
EOF
}

# Function to generate prefill content for cgpt
generate_prefill() {
    local command_name="$1"
    cat <<EOF
<ant-thinking>
To implement the '${command_name}' command for the nvcf CLI, I'll need to:
1. Analyze the existing codebase structure
2. Identify the appropriate package and file for the new command
3. Implement the command functionality
4. Add error handling and logging
5. Write comprehensive comments and documentation
6. Consider adding unit tests
</ant-thinking>

Here's the Go code for the '${command_name}' command:

\`\`\`go
EOF
}

# Function to get the existing codebase structure
get_codebase_structure() {
    tree -L 2 ~/go/src/github.com/tmc/nvcf
}

# Main script logic
main() {
    local command_name="$1"
    local output_file="${command_name// /_}.go"

    log "Generating Go code for nvcf command: $command_name"

    # Generate system prompt and prefill content
    local system_prompt=$(generate_system_prompt)
    local prefill=$(generate_prefill "$command_name")

    # Get the existing codebase structure
    local codebase_structure=$(get_codebase_structure)

    # Execute cgpt command
    log "Executing cgpt to generate Go code"
    (
        echo "Generate Go code for the '${command_name}' command in the nvcf CLI tool. Consider the following codebase structure:"
        echo "$codebase_structure"
        echo "Ensure the new code is consistent with the existing structure and follows Go best practices."
    ) | cgpt -s "$system_prompt" -p "$prefill" > "$output_file"

    log "Go code generated and saved to $output_file"

    echo "Generated Go code contents:"
    echo "------------------------"
    cat "$output_file"
    echo "------------------------"
    echo "You can now integrate this code into the appropriate file in the nvcf codebase."
}

# Check if a command name is provided
if [ $# -eq 0 ]; then
    log "Error: No command name provided"
    usage
fi

# Run the main function
main "$*"
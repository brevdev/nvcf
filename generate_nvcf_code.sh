#!/usr/bin/env bash

set -euo pipefail

# Function to display script usage
usage() {
    echo "Usage: $0 <task_description>"
    echo "Generates Go code for the nvcf codebase using cgpt."
    exit 1
}

# Function to log messages
log() {
    echo "[$(date +'%Y-%m-%d %H:%M:%S')] $1" >&2
}

# Function to generate a system prompt for cgpt
generate_system_prompt() {
    cat <<EOF
You are an expert Go developer specializing in CLI tools and SDK libraries. Your task is to generate Go code for the nvcf codebase based on the given task description. Follow these guidelines:
1. Use idiomatic Go patterns and best practices
2. Implement robust error handling and logging
3. Add comprehensive comments and documentation
4. Ensure code is efficient and follows Go style guidelines
5. Consider potential edge cases and handle them appropriately
6. Integrate well with the existing nvcf codebase structure
EOF
}

# Function to generate prefill content for cgpt
generate_prefill() {
    local task="$1"
    cat <<EOF
<ant-thinking>
To implement the task "${task}" for the nvcf codebase, I'll need to:
1. Analyze the requirements and existing codebase structure
2. Design a suitable implementation approach
3. Write efficient and idiomatic Go code
4. Implement error handling and logging
5. Add comprehensive comments and documentation
6. Consider integration with existing components
7. Handle potential edge cases
</ant-thinking>

Here's the Go code implementation for the task:

\`\`\`go
EOF
}

# Function to get the nvcf codebase structure
get_codebase_structure() {
    find . -type f -name "*.go" | sort | xargs -I {} echo "File: {}"
}

# Main script logic
main() {
    # Check if a task description is provided
    if [ $# -eq 0 ]; then
        log "Error: No task description provided"
        usage
    fi

    local task="$*"
    local output_file="${task// /_}.go"

    log "Generating Go code for task: $task"

    # Generate system prompt and prefill content
    local system_prompt=$(generate_system_prompt)
    local prefill=$(generate_prefill "$task")

    # Get the nvcf codebase structure
    local codebase_structure=$(get_codebase_structure)

    # Execute cgpt command
    log "Executing AI model to generate Go code"
    (
        echo "Generate Go code for the nvcf codebase to perform the following task: $task"
        echo "Existing codebase structure:"
        echo "$codebase_structure"
        echo "Ensure the generated code integrates well with the existing structure."
    ) | cgpt -s "$system_prompt" -p "$prefill" > "$output_file"

    log "Go code generated and saved to $output_file"

    echo "Generated Go code contents:"
    echo "------------------------"
    cat "$output_file"
    echo "------------------------"
    echo "You can find the generated Go code in: $output_file"
}

# Run the main function
main "$@"
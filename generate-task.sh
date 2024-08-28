#!/usr/bin/env bash

set -euo pipefail

# Function to display script usage
usage() {
    echo "Usage: $0 <task_description>"
    echo "Generates Go code to perform the specified task for the nvcf project."
    exit 1
}

# Function to log messages
log() {
    echo "[$(date +'%Y-%m-%d %H:%M:%S')] $1" >&2
}

# Function to generate a system prompt
generate_system_prompt() {
    cat <<EOF
You are an expert Go developer specializing in CLI tools and the nvcf project. Your task is to create Go code that accomplishes the given task. Follow these guidelines:
1. Use Go best practices and idiomatic Go
2. Implement robust error handling and logging
3. Add comprehensive comments and documentation
4. Make the code modular and reusable where possible
5. Consider the existing code structure and recent changes in the git history
6. Ensure the new code integrates well with the existing codebase

Use cgpt for AI tasks: $(cgpt -h 2>&1)
EOF
}

# Function to generate prefill content
generate_prefill() {
    local task="$1"
    cat <<EOF
<ant-thinking>
To implement the task "${task}" for the nvcf project, I'll need to:
1. Analyze the requirements of the task
2. Review the existing code structure and recent git history
3. Design a solution that integrates well with the current codebase
4. Implement the core functionality
5. Add error handling and logging
6. Include documentation and comments
7. Consider potential edge cases and handle them appropriately
</ant-thinking>

Here's the Go code to ${task}:

\`\`\`go
EOF
}

# Function to get recent git history
get_recent_git_history() {
    git log -n 10 --pretty=format:"%h %s" -- ./cmd/*.go main.go
}

# Function to get the script's source code
get_source_code() {
    cat "$0"
}

# Main script logic
main() {
    local task="$1"
    local output_file="${task// /_}.go"

    log "Generating Go code for task: $task"

    # Generate system prompt and prefill content
    local system_prompt=$(generate_system_prompt)
    local prefill=$(generate_prefill "$task")

    # Get recent git history
    local git_history=$(get_recent_git_history)

    # Get the script's source code
    local source_code=$(get_source_code)

    # Execute cgpt command
    log "Executing AI model to generate code"
    (
        echo "Create Go code to perform the following task: $task"
        echo "Recent git history of relevant files:"
        echo "$git_history"
        echo "Current content of main.go:"
        cat main.go
        echo "Current content of cmd/*.go files:"
        cat ./cmd/*.go
        echo "Content of generate-task.sh:"
        echo "$source_code"
    ) | cgpt -s "$system_prompt" -p "$prefill" > "$output_file"

    log "Code generated and saved to $output_file"

    echo "Generated code contents:"
    echo "------------------------"
    cat "$output_file"
    echo "------------------------"
    echo "You can review and edit the generated code in: $output_file"
}

# Check if a task is provided
if [ $# -eq 0 ]; then
    log "Error: No task description provided"
    usage
fi

# Run the main function
main "$*"
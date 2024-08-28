#!/usr/bin/env bash

set -euo pipefail

# Function to display script usage
usage() {
    echo "Usage: $0 <task_description>"
    echo "Generates a bash script to perform the specified task."
    exit 1
}

# Function to log messages
log() {
    echo "[$(date +'%Y-%m-%d %H:%M:%S')] $1" >&2
}

# Function to generate a system prompt
generate_system_prompt() {
    cat <<EOF
You are an expert bash script developer specializing in CLI tools and automation. Your task is to create a bash script that accomplishes the given task. Follow these guidelines:
1. Use bash best practices and modern syntax (bash 4+)
2. Implement robust error handling and logging
3. Add comprehensive comments and usage information
4. Make the script modular and reusable where possible
5. Consider adding command-line argument parsing for flexibility
EOF
}

# Function to generate prefill content
generate_prefill() {
    local task="$1"
    cat <<EOF
<ant-thinking>
To implement a bash script for "${task}", I'll need to:
1. Analyze the requirements of the task
2. Design a modular script structure
3. Implement the core functionality
4. Add error handling and logging
5. Include usage information and help text
6. Consider potential edge cases and handle them appropriately
</ant-thinking>

Here's the bash script to ${task}:

\`\`\`bash
#!/usr/bin/env bash
EOF
}

# Main script logic
main() {
    local task="$1"
    local output_file="${task// /_}.sh"

    log "Generating bash script for task: $task"

    # Generate system prompt and prefill content
    local system_prompt=$(generate_system_prompt)
    local prefill=$(generate_prefill "$task")

    # Execute cgpt command
    log "Executing AI model to generate script"
    (echo "Create a bash script to perform the following task: $task"; echo "$0") | \
    cgpt -s "$system_prompt" -p "$prefill" > "$output_file"

    log "Script generated and saved to $output_file"
    chmod +x "$output_file"
    log "Made $output_file executable"

    echo "Generated script contents:"
    echo "------------------------"
    cat "$output_file"
    echo "------------------------"
    echo "You can run the generated script with: ./$output_file"
}

# Check if a task is provided
if [ $# -eq 0 ]; then
    log "Error: No task description provided"
    usage
fi

# Run the main function
main "$*"

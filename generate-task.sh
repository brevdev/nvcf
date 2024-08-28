#!/usr/bin/env bash

set -euo pipefail

# Function to generate a system prompt
generate_system_prompt() {
    cat <<EOF
You are an expert Go developer specializing in CLI tools and SDK libraries. Your task is to extend and improve the nvcf codebase. Follow these guidelines:
1. Maintain consistency with the existing code style and structure.
2. Implement robust error handling and logging.
3. Add comprehensive comments and documentation.
4. Consider adding unit tests for new functionality.

Use cgpt for AI tasks: <cgpt-usage cmd
EOF
}

# Function to generate prefill content
generate_prefill() {
    local task="$1"
    cat <<EOF
<ant-thinking>
To implement ${task}, I'll need to:
1. Analyze the existing codebase structure
2. Identify the appropriate package and file for the new functionality
3. Implement the feature while maintaining consistency with existing code
4. Add error handling and logging
5. Write comprehensive comments and documentation
6. Consider adding unit tests
</ant-thinking>

Here's the implementation for ${task}:

\`\`\`go
EOF
}

# Function to get the script's source code
get_source_code() {
    cat "$0"
}

# Main script
main() {
    local task="$1"
    local output_file="${task// /_}.go"

    # Generate system prompt and prefill content
    local system_prompt=$(generate_system_prompt)
    local prefill=$(generate_prefill "$task")

    # Get the script's source code
    local source_code=$(get_source_code)

    # Execute cgpt command
    echo "Implementing: $task"
    (
        echo "Implement the following feature for the nvcf codebase: $task"
        echo "Here's the source code of the script that generated this task (you can invoke it again)
        echo "\`\`\`bash"
        echo "$source_code"
        echo "\`\`\`"
        ~/code-to-gpt.sh
    ) | cgpt -s "$system_prompt" -p "$prefill" | tee "$output_file"

    echo "Output saved to $output_file"
}

# Check if a task is provided
if [ $# -eq 0 ]; then
    echo "Usage: $0 <task_description>"
    exit 1
fi

# Run the main function
main "$*"

#!/usr/bin/env bash

set -euo pipefail

# Function to generate a system prompt
generate_system_prompt() {
    cat <<EOF
You are an AI assistant specialized in generating Bash scripts. Your task is to create a script that uses cgpt to generate Go code for the nvcf codebase. Follow these guidelines:
1. Maintain a clear and organized script structure.
2. Implement robust error handling and input validation.
3. Add comprehensive comments and documentation.
4. Ensure the script is efficient and follows Bash best practices.
5. Implement proper quoting and escaping of variables.
6. Use cgpt for AI tasks and handle its output appropriately.

Use cgpt for AI tasks: $(cgpt -h 2>&1)
EOF
}

# Function to generate prefill content
generate_prefill() {
    cat <<EOF
<ant-thinking>
To create a script that uses cgpt to generate Go code for the nvcf codebase, I'll need to:
1. Analyze the existing script structure
2. Identify the key components and functions
3. Adapt the script to use cgpt for generating Go code
4. Add error handling and input validation
5. Write comprehensive comments and documentation
6. Ensure efficient implementation and follow Bash best practices
7. Implement proper quoting and escaping of variables
8. Handle cgpt output appropriately
</ant-thinking>

<ant-scratchpad>
- Review the existing script for reusable components
- Consider potential edge cases and error scenarios
- Plan for extensibility and future improvements
- Ensure proper integration with cgpt
</ant-scratchpad>

Here's the implementation of the new script:

\`\`\`bash
EOF
}

# Function to get the script's source code
get_source_code() {
    cat "$0"
}

# Main script
main() {
    local output_file="generate_nvcf_code.sh"

    # Generate system prompt and prefill content
    local system_prompt=$(generate_system_prompt)
    local prefill=$(generate_prefill)

    # Get the script's source code
    local source_code=$(get_source_code)

    # Execute cgpt command
    echo "Generating script to create Go code for nvcf codebase using cgpt"
    (
        echo "Create a Bash script that uses cgpt to generate Go code for the nvcf codebase. The script should be based on the following source code but adapted to use cgpt for generating Go code:"
        echo "\`\`\`bash"
        echo "$source_code"
        echo "\`\`\`"
        ~/code-to-gpt.sh
    ) | cgpt -s "$system_prompt" -p "$prefill" | tee "$output_file"

    echo "Output saved to $output_file"
    chmod +x "$output_file"
}

# Run the main function
main
#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color
BLUE='\033[0;34m'

# Server URL
SERVER="http://localhost:8080"

# Function to run a query and display results
run_query() {
    local file=$1
    local name=$(basename "$file" .yaml)
    
    echo -e "\n${BLUE}Running test: ${name}${NC}"
    echo "Query file: $file"
    echo "----------------------------------------"
    bat --color=always "$file"
    echo "----------------------------------------"
    
    # Run the query and capture the response
    response=$(curl -s -X POST \
        -H "Content-Type: application/yaml" \
        --data-binary "@$file" \
        "${SERVER}/search")
    
    # Check if the request was successful
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}Success!${NC}"
        # Pretty print the JSON response
        echo "$response" | jq '.'
    else
        echo -e "${RED}Failed to execute query${NC}"
        echo "$response"
    fi
    
    echo "----------------------------------------"
}

# First check if the server is running
if ! curl -s "${SERVER}" > /dev/null; then
    echo -e "${RED}Error: Server is not running at ${SERVER}${NC}"
    echo "Please start the server first with: go run ."
    exit 1
fi

# Check if jq is installed
if ! command -v jq &> /dev/null; then
    echo -e "${RED}Error: jq is not installed${NC}"
    echo "Please install jq first:"
    echo "  macOS: brew install jq"
    echo "  Linux: sudo apt-get install jq"
    exit 1
fi

# Run all test queries
echo -e "${BLUE}Starting test queries...${NC}"

# Get all YAML files in the current directory
for file in $(dirname "$0")/*.yaml; do
    if [ -f "$file" ] && [ "$file" != "$(dirname "$0")/run-tests.sh" ]; then
        run_query "$file"
        # Add a small delay between requests
        sleep 1
    fi
done

echo -e "\n${GREEN}All tests completed!${NC}" 
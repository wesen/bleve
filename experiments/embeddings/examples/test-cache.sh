#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color
BLUE='\033[0;34m'

# Server URL
SERVER="http://localhost:8080"

echo -e "${BLUE}Testing Embeddings Cache Functionality${NC}"

# Check if the server is running
if ! curl -s "${SERVER}" > /dev/null; then
    echo -e "${RED}Error: Server is not running at ${SERVER}${NC}"
    echo "Please start the server first with: go run ."
    exit 1
fi

# Test 1: Get cache stats
echo -e "\n${BLUE}Test 1: Get cache stats${NC}"
curl -s "${SERVER}/cache/stats" | jq '.'

# Test 2: Generate an embedding
echo -e "\n${BLUE}Test 2: Generate an embedding${NC}"
FIRST_TEXT="The quick brown fox jumps over the lazy dog"

curl -s -X POST \
  -H "Content-Type: application/json" \
  -d "{\"text\": \"${FIRST_TEXT}\"}" \
  "${SERVER}/embeddings" | jq '.dimensions, .model'

# Test 3: Generate the same embedding again (should use cache)
echo -e "\n${BLUE}Test 3: Generate the same embedding again (should use cache)${NC}"
echo "First request time (should be slower):"
time curl -s -X POST \
  -H "Content-Type: application/json" \
  -d "{\"text\": \"${FIRST_TEXT}\"}" \
  "${SERVER}/embeddings" > /dev/null

echo "Second request time (should be faster due to cache):"
time curl -s -X POST \
  -H "Content-Type: application/json" \
  -d "{\"text\": \"${FIRST_TEXT}\"}" \
  "${SERVER}/embeddings" > /dev/null

# Test 4: Check cache stats again
echo -e "\n${BLUE}Test 4: Check cache stats again${NC}"
curl -s "${SERVER}/cache/stats" | jq '.'

# Test 5: Clear cache
echo -e "\n${BLUE}Test 5: Clear cache${NC}"
curl -s -X POST "${SERVER}/cache/clear" | jq '.'

# Test 6: Check cache stats after clearing
echo -e "\n${BLUE}Test 6: Check cache stats after clearing${NC}"
curl -s "${SERVER}/cache/stats" | jq '.'

# Test 7: Test vector search
echo -e "\n${BLUE}Test 7: Test vector search${NC}"
curl -s -X POST \
  -H "Content-Type: application/yaml" \
  -d 'query:
  vector:
    field: vector
    text: "What is the meaning of life?"
    model: all-minilm
    k: 10
    boost: 1.0
options:
  size: 10
  highlight:
    fields: [content]' \
  "${SERVER}/search" | jq '.total, .took, .hits[].id, .hits[].score'

echo -e "\n${GREEN}All tests completed!${NC}" 
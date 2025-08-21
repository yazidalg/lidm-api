#!/bin/bash

# Test script for module progress functionality
BASE_URL="http://localhost:3000"

echo "Testing Module Progress Implementation"
echo "====================================="

# First, login to get a token
echo "1. Testing login..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "anjayy@gmail.com",
    "password": "uhuy123!"
  }')

echo "Login response: $LOGIN_RESPONSE"

# Extract token from response
TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.token')

if [ "$TOKEN" = "null" ] || [ -z "$TOKEN" ]; then
    echo "Failed to get token. Trying belajar-login..."
    LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/belajar-login" \
      -H "Content-Type: application/json" \
      -d '{
        "email": "anjayy@gmail.com",
        "password": "uhuy123!"
      }')
    echo "Belajar login response: $LOGIN_RESPONSE"
    TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.token')
fi

if [ "$TOKEN" = "null" ] || [ -z "$TOKEN" ]; then
    echo "Failed to authenticate. Please check credentials."
    exit 1
fi

echo "Token obtained: ${TOKEN:0:20}..."
echo ""

# Test the enhanced module/all endpoint
echo "2. Testing GET /module/all (with progress)..."
MODULE_RESPONSE=$(curl -s -X GET "$BASE_URL/module/all" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json")

echo "Module response:"
echo $MODULE_RESPONSE | jq '.'
echo ""

# Test individual module progress endpoint 
echo "3. Testing GET /module/1/progress..."
MODULE_PROGRESS_RESPONSE=$(curl -s -X GET "$BASE_URL/module/1/progress" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json")

echo "Module progress response:"
echo $MODULE_PROGRESS_RESPONSE | jq '.'
echo ""

# Test progress alias endpoint
echo "4. Testing GET /progress/module/1..."
PROGRESS_ALIAS_RESPONSE=$(curl -s -X GET "$BASE_URL/progress/module/1" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json")

echo "Progress alias response:"
echo $PROGRESS_ALIAS_RESPONSE | jq '.'
echo ""

echo "====================================="
echo "Test completed!"

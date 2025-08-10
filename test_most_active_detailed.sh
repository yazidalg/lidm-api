#!/bin/bash

# Test Most Active User Detailed endpoint

echo "Testing Most Active User Detailed endpoint..."

# Login first to get token (using belajar login)
echo "Getting authentication token..."
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:3000/auth/belajar-login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "azis@belajar.id",
    "password": "password123"
  }')

echo "Login Response: $LOGIN_RESPONSE"

# Extract token from response
TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
    echo "Failed to get authentication token"
    exit 1
fi

echo "Token obtained: ${TOKEN:0:20}..."

# Test the detailed most active user endpoint
echo -e "\n--- Testing Most Active User Detailed ---"
curl -X GET "http://localhost:3000/user-activity/most-active-detailed" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" | jq '.'

echo -e "\nTest completed!"

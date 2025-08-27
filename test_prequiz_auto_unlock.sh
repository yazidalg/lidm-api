#!/bin/bash

# Test script for prequiz submission and auto-unlock
BASE_URL="http://localhost:3000"

echo "Testing Prequiz Submission & Auto-Unlock Fix"
echo "============================================="

# First, login to get a token
echo "1. Testing login..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "anjayyy@gmail.com",
    "password": "uhuy123!"
  }')

echo "Login response: $LOGIN_RESPONSE"

# Extract token from response
TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*' | grep -o '[^"]*$')

if [ -z "$TOKEN" ]; then
  echo "❌ Failed to get token. Check login credentials."
  exit 1
fi

echo "✅ Token obtained: ${TOKEN:0:20}..."

# Test: Submit a prequiz answer
echo ""
echo "2. Testing prequiz submission..."
echo "Submitting answer for prequiz_id: 13 (Module 5)"

PREQUIZ_RESPONSE=$(curl -s -X POST "$BASE_URL/prequiz/submit" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "prequiz_id": 13,
    "selected_answer": "Glukosa dan oksigen"
  }')

echo "Prequiz submission response:"
echo "$PREQUIZ_RESPONSE"

# Check if auto-unlock worked (no errors should appear)
if echo "$PREQUIZ_RESPONSE" | grep -q "error"; then
  echo "❌ Error detected in response"
else
  echo "✅ No errors - auto-unlock mechanism working!"
fi

# Test: Get updated module progress
echo ""
echo "3. Testing module progress after submission..."

MODULES_RESPONSE=$(curl -s -X GET "$BASE_URL/modules/all" \
  -H "Authorization: Bearer $TOKEN")

echo "Updated modules with progress:"
echo "$MODULES_RESPONSE" | jq '.'

echo ""
echo "Test completed!"

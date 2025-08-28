#!/bin/bash

# Simple test untuk is_current_user functionality
BASE_URL="http://localhost:3000"

echo "🧪 Testing is_current_user functionality..."
echo "============================================"

echo "Please provide your JWT token:"
read -r JWT_TOKEN

if [ -z "$JWT_TOKEN" ]; then
    echo "❌ JWT token is required!"
    exit 1
fi

echo -e "\n📊 Testing leaderboard..."
echo "Looking for is_current_user: true in the response"
echo -e "\nResponse:"

curl -s -X GET "$BASE_URL/leaderboard" \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -H "Content-Type: application/json"

echo -e "\n\n✅ Test completed!"
echo "💡 Look for \"is_current_user\": true in the JSON response above"
echo "📝 This should appear for exactly one user (the one who owns the JWT token)"

#!/bin/bash

# Simple Leaderboard API Test Script
# Update JWT_TOKEN with your actual token
BASE_URL="http://localhost:3000"
JWT_TOKEN="YOUR_JWT_TOKEN_HERE"

echo "🏆 Testing Leaderboard API..."
echo "==============================="

# Test 1: Get all leaderboard
echo "📊 Test 1: Get all leaderboard"
curl -X GET "$BASE_URL/leaderboard" \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -H "Content-Type: application/json"

echo -e "\n\n"

# Test 2: Get leaderboard for module 1
echo "📊 Test 2: Get leaderboard for module 1"
curl -X GET "$BASE_URL/leaderboard?module_id=1" \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -H "Content-Type: application/json"

echo -e "\n\n"

# Test 3: Get leaderboard for solo quiz
echo "📊 Test 3: Get leaderboard for solo quiz"
curl -X GET "$BASE_URL/leaderboard?quiz_type=solo" \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -H "Content-Type: application/json"

echo -e "\n\n"

# Test 4: Get user rank for user ID 1
echo "👤 Test 4: Get user rank for user ID 1"
curl -X GET "$BASE_URL/leaderboard/user/1" \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -H "Content-Type: application/json"

echo -e "\n\n"

echo "✅ Basic leaderboard tests completed!"
echo "💡 Remember to replace YOUR_JWT_TOKEN_HERE with your actual JWT token"

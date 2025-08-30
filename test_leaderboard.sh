#!/bin/bash

# Leaderboard API Test Script
# Make sure to update these variables with your actual values
BASE_URL="http://localhost:3000"
JWT_TOKEN="YOUR_JWT_TOKEN_HERE"

echo "🏆 Testing Leaderboard API..."
echo "==============================="

# Test 1: Get all leaderboard
echo "📊 Test 1: Get all leaderboard"
curl -X GET "$BASE_URL/leaderboard" \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -H "Content-Type: application/json" \
  | jq '.'

echo -e "\n"

# Test 2: Get leaderboard for module 1
echo "📊 Test 2: Get leaderboard for module 1"
curl -X GET "$BASE_URL/leaderboard?module_id=1" \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -H "Content-Type: application/json" \
  | jq '.'

echo -e "\n"

# Test 3: Get leaderboard for solo quiz
echo "📊 Test 3: Get leaderboard for solo quiz"
curl -X GET "$BASE_URL/leaderboard?quiz_type=solo" \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -H "Content-Type: application/json" \
  | jq '.'

echo -e "\n"

# Test 4: Get leaderboard for matchmaking quiz
echo "📊 Test 4: Get leaderboard for matchmaking quiz"  
curl -X GET "$BASE_URL/leaderboard?quiz_type=matchmaking" \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -H "Content-Type: application/json" \
  | jq '.'

echo -e "\n"

# Test 5: Get leaderboard for module 1 solo quiz
echo "📊 Test 5: Get leaderboard for module 1 solo quiz"
curl -X GET "$BASE_URL/leaderboard?module_id=1&quiz_type=solo" \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -H "Content-Type: application/json" \
  | jq '.'

echo -e "\n"

# Test 6: Get user rank for user ID 1
echo "👤 Test 6: Get user rank for user ID 1"
curl -X GET "$BASE_URL/leaderboard/user/1" \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -H "Content-Type: application/json" \
  | jq '.'

echo -e "\n"

# Test 7: Get user rank for user ID 1 in module 1
echo "👤 Test 7: Get user rank for user ID 1 in module 1"
curl -X GET "$BASE_URL/leaderboard/user/1?module_id=1" \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -H "Content-Type: application/json" \
  | jq '.'

echo -e "\n"

# Test 8: Get user rank for user ID 1 in solo quiz
echo "👤 Test 8: Get user rank for user ID 1 in solo quiz"
curl -X GET "$BASE_URL/leaderboard/user/1?quiz_type=solo" \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -H "Content-Type: application/json" \
  | jq '.'

echo -e "\n"

echo "✅ All leaderboard tests completed!"

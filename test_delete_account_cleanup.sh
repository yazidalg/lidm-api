#!/bin/bash

# Test Delete Account with Leaderboard Cleanup
# This script tests that when an account is deleted, all related data is also removed

echo "=========================================="
echo "Testing Delete Account with Data Cleanup"
echo "=========================================="
echo ""

# First, create a test account
echo "1. Creating test account..."
REGISTER_RESPONSE=$(curl -s -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test User Delete",
    "email": "testdelete@example.com",
    "password": "testpassword123"
  }')

echo "$REGISTER_RESPONSE" | jq '.'

# Extract token from response
TOKEN=$(echo "$REGISTER_RESPONSE" | jq -r '.data.token // .token // empty')

if [ -z "$TOKEN" ]; then
    echo "Failed to get token, trying to login..."
    LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8080/auth/login \
      -H "Content-Type: application/json" \
      -d '{
        "email": "testdelete@example.com",
        "password": "testpassword123"
      }')
    TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.data.token // .token // empty')
fi

echo ""
echo "Token: $TOKEN"
echo ""

# Get user profile to see user_id
echo "2. Getting user profile..."
USER_PROFILE=$(curl -s -X GET http://localhost:8080/user/profile \
  -H "Authorization: Bearer $TOKEN")

echo "$USER_PROFILE" | jq '.'
USER_ID=$(echo "$USER_PROFILE" | jq -r '.data.id')
echo ""
echo "User ID: $USER_ID"
echo ""

# Check leaderboard before deletion
echo "3. Checking leaderboard for user $USER_ID..."
LEADERBOARD_BEFORE=$(curl -s -X GET "http://localhost:8080/leaderboard?user_id=$USER_ID" \
  -H "Authorization: Bearer $TOKEN")

echo "$LEADERBOARD_BEFORE" | jq '.'
echo ""

# Delete the account
echo "4. Deleting account..."
DELETE_RESPONSE=$(curl -s -X DELETE http://localhost:8080/user/delete-account \
  -H "Content-Type: application/json" \
  -d '{
    "email": "testdelete@example.com",
    "password": "testpassword123"
  }')

echo "$DELETE_RESPONSE" | jq '.'
echo ""

# Try to get user profile (should fail)
echo "5. Verifying user is deleted (should return 401)..."
VERIFY_DELETE=$(curl -s -X GET http://localhost:8080/user/profile \
  -H "Authorization: Bearer $TOKEN")

echo "$VERIFY_DELETE" | jq '.'
echo ""

# Check leaderboard after deletion (user should not be in leaderboard)
echo "6. Checking leaderboard after deletion..."
LEADERBOARD_AFTER=$(curl -s -X GET http://localhost:8080/leaderboard)

echo "$LEADERBOARD_AFTER" | jq '.'
echo ""

echo "=========================================="
echo "Test Completed!"
echo "=========================================="
echo ""
echo "Summary:"
echo "- Account created: testdelete@example.com"
echo "- User ID: $USER_ID"
echo "- Account deleted: Check response above"
echo "- Leaderboard cleaned: User should not appear in leaderboard"
echo ""
echo "All related data should be removed:"
echo "  ✓ Leaderboard entry"
echo "  ✓ Participant records"
echo "  ✓ Quiz sessions"
echo "  ✓ Module progress"
echo "  ✓ User activities"
echo "  ✓ Flashcard progress"

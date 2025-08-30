#!/bin/bash

# Test script untuk memverifikasi is_current_user
BASE_URL="http://localhost:3000"

echo "🧪 Testing is_current_user functionality..."
echo "============================================"

echo "Please provide your JWT token:"
read -r JWT_TOKEN

if [ -z "$JWT_TOKEN" ]; then
    echo "❌ JWT token is required!"
    exit 1
fi

echo -e "\n📊 Testing leaderboard with current user marking..."

# Test leaderboard endpoint
echo "Request: GET /leaderboard"
echo "Expected: One user should have is_current_user: true"
echo -e "\nResponse:"

curl -s -X GET "$BASE_URL/leaderboard" \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -H "Content-Type: application/json" \
  | jq '
    .juara1.is_current_user as $j1 |
    .juara2.is_current_user as $j2 |
    .juara3.is_current_user as $j3 |
    (.leaderboard[] | select(.is_current_user == true)) as $current |
    {
      "summary": {
        "juara1_is_current": $j1,
        "juara2_is_current": $j2, 
        "juara3_is_current": $j3,
        "current_user_found": ($j1 or $j2 or $j3 or ($current != null)),
        "current_user_position": (
          if $j1 then "juara1"
          elif $j2 then "juara2"
          elif $j3 then "juara3"
          elif $current then "rank_\($current.rank)"
          else "not_found"
          end
        )
      },
      "current_user_details": (
        if $j1 then .juara1
        elif $j2 then .juara2
        elif $j3 then .juara3
        else $current
        end
      )
    }
  '

echo -e "\n\n✅ Test completed!"
echo "💡 Check if current_user_found is true and current_user_details shows the logged-in user"

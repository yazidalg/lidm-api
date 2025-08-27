#!/bin/bash

echo "🧪 Testing Database Trigger Auto-Unlock..."

# Get JWT token
echo "Getting JWT token..."
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:3000/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "anjayyy@gmail.com",
    "password": "123456"
  }')

JWT_TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -z "$JWT_TOKEN" ]; then
    echo "❌ Failed to get JWT token"
    exit 1
fi

echo "✅ JWT Token obtained"

# Test 1: Submit prequiz (should trigger auto-unlock if module becomes complete)
echo -e "\n📝 Testing prequiz submission (trigger test)..."
PREQUIZ_RESPONSE=$(curl -s -X POST http://localhost:3000/prequiz/submit \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -d '{
    "prequiz_id": 13,
    "selected_answer": "Glukosa dan oksigen"
  }')

echo "Prequiz Response: $PREQUIZ_RESPONSE"

# Test 2: Check modules status after submission
echo -e "\n📊 Checking modules status after trigger..."
MODULES_RESPONSE=$(curl -s -X GET http://localhost:3000/modules/all \
  -H "Authorization: Bearer $JWT_TOKEN")

echo "Modules Response (focusing on Module 5):"
echo $MODULES_RESPONSE | python3 -m json.tool | grep -A 20 -B 5 '"id": 5'

echo -e "\n🎯 Test completed! Check if Module 5 is now unlocked."

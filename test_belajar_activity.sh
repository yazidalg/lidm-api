#!/bin/bash

# Script untuk test login dengan akun belajar dan activity tracking
BASE_URL="http://localhost:8080"

echo "=== Testing Belajar Login & Activity Tracking ==="
echo

# 1. Logout dulu untuk clear session
echo "1. Logout untuk clear session..."
curl -X POST "$BASE_URL/auth/logout" \
  -H "Content-Type: application/json" \
  -c cookies.txt \
  -s
echo "✅ Logout done"
echo

# 2. Login dengan akun belajar (simulasi)
echo "2. Login dengan akun belajar..."
read -p "Masukkan Belajar ID Token: " BELAJAR_TOKEN

LOGIN_RESPONSE=$(curl -X POST "$BASE_URL/auth/belajar-login" \
  -H "Content-Type: application/json" \
  -d "{\"id_token\": \"$BELAJAR_TOKEN\"}" \
  -c cookies.txt \
  -s)

echo "Login Response: $LOGIN_RESPONSE"
echo

# Extract token dari response
TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
  echo "❌ Login gagal, tidak dapat token"
  exit 1
fi

echo "✅ Login berhasil! Token: ${TOKEN:0:50}..."
echo

# 3. Test beberapa endpoint untuk trigger activity
echo "3. Testing various endpoints to trigger activities..."

echo "   - Get user profile..."
curl -X GET "$BASE_URL/user/profile" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -s | jq .

echo "   - Get all modules..."
curl -X GET "$BASE_URL/module/all" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -s | jq .

echo "   - Get all lessons..."
curl -X GET "$BASE_URL/lesson/all" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -s | jq .

echo "   - Get all quizzes..."
curl -X GET "$BASE_URL/quiz/all" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -s | jq .

echo

# 4. Check user activities
echo "4. Check user activities..."
curl -X GET "$BASE_URL/user-activity/my-activities" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -s | jq .

echo

# 5. Logout
echo "5. Logout..."
curl -X POST "$BASE_URL/auth/logout" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -s

echo "✅ Test completed!"
echo "Check database untuk melihat activity logs!"

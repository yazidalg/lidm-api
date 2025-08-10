#!/bin/bash

# Script untuk test user activity tracking
BASE_URL="http://localhost:3000"

echo "=== Testing User Activity Tracking ==="
echo

# 1. Logout dulu untuk clear session
echo "1. Logout untuk clear session..."
curl -X POST "$BASE_URL/auth/logout" \
  -H "Content-Type: application/json" \
  -c cookies.txt \
  -v
echo -e "\n"

# 2. Login dengan Google (simulasi)
echo "2. Login dengan Google..."
read -p "Masukkan Google ID Token: " GOOGLE_TOKEN

LOGIN_RESPONSE=$(curl -X POST "$BASE_URL/auth/google" \
  -H "Content-Type: application/json" \
  -d "{\"id_token\": \"$GOOGLE_TOKEN\"}" \
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

# 3. Test get user profile (untuk trigger activity tracking)
echo "3. Get user profile..."
curl -X GET "$BASE_URL/user/profile" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -v
echo -e "\n"

# 4. Test get modules (untuk trigger module activity)
echo "4. Get all modules..."
curl -X GET "$BASE_URL/module/all" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -v
echo -e "\n"

# 5. Test get lessons (untuk trigger lesson activity)
echo "5. Get all lessons..."
curl -X GET "$BASE_URL/lesson/all" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -v
echo -e "\n"

# 6. Test logout
echo "6. Logout..."
curl -X POST "$BASE_URL/auth/logout" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -v
echo -e "\n"

echo "=== Test Selesai ==="
echo "Cek database untuk melihat activity logs yang tercatat!"

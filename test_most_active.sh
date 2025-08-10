#!/bin/bash

# Script untuk test most active users
BASE_URL="http://localhost:3000"

echo "=== Tes Endpoint Pengguna Paling Aktif ==="
echo

# Pertama, pastikan kita login sebagai admin untuk mengakses endpoint most-active
echo "1. Login sebagai admin dulu..."

# Login dengan akun admin (perlu token yang valid)
echo "Silakan login dengan kredensial admin untuk mendapatkan token admin:"
read -p "Masukkan JWT token admin: " ADMIN_TOKEN

if [ -z "$ADMIN_TOKEN" ]; then
  echo "❌ Token tidak boleh kosong"
  exit 1
fi

echo "✅ Menggunakan token admin: ${ADMIN_TOKEN:0:50}..."
echo

# Test most active users endpoint
echo "2. Dapatkan Pengguna Paling Aktif..."
curl -X GET "$BASE_URL/user-activity/most-active?limit=10" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -s | jq .

echo

# Test most active user dengan detail lengkap (hanya 1 user)
echo "3. Dapatkan Pengguna Paling Aktif dengan Detail Lengkap..."
curl -X GET "$BASE_URL/user-activity/most-active-detailed" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -s | jq .

echo

# Test recent activities
echo "4. Dapatkan Aktivitas Terbaru..."
curl -X GET "$BASE_URL/user-activity/recent?limit=20" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -s | jq .

echo

# Test activity stats
echo "5. Dapatkan Statistik Aktivitas..."
curl -X GET "$BASE_URL/user-activity/stats" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -s | jq .

echo

# Test dashboard (yang include most active users)
echo "6. Dapatkan Data Dashboard..."
curl -X GET "$BASE_URL/dashboard/" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -s | jq .

echo
echo "=== Tes Pengguna Paling Aktif Selesai ==="

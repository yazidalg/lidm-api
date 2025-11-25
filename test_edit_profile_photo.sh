#!/bin/bash

# Test Edit Profile dengan Photo Profile URL
# Pastikan sudah ada TOKEN yang valid

if [ -z "$TOKEN" ]; then
    echo "Error: TOKEN environment variable is not set"
    echo "Please run: export TOKEN=your_jwt_token_here"
    exit 1
fi

echo "=========================================="
echo "Testing Edit Profile with Photo Profile"
echo "=========================================="
echo ""

echo "1. Edit Profile - Update Name, Email, dan Photo Profile URL"
curl -X PUT http://localhost:8080/user/edit-profile \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated User Name",
    "email": "updated@email.com",
    "photo_profile": "https://example.com/profile-photo.jpg"
  }'

echo ""
echo ""
echo "=========================================="
echo "2. Edit Profile - Hanya Update Name (Photo Profile tidak berubah)"
curl -X PUT http://localhost:8080/user/edit-profile \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Another Name Update",
    "email": "updated@email.com"
  }'

echo ""
echo ""
echo "=========================================="
echo "3. Edit Profile - Update dengan Photo Profile Kosong (tidak akan update photo)"
curl -X PUT http://localhost:8080/user/edit-profile \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Final Name",
    "email": "final@email.com",
    "photo_profile": ""
  }'

echo ""
echo ""
echo "=========================================="
echo "4. Cek Profile Setelah Update"
curl -X GET http://localhost:8080/user/profile \
  -H "Authorization: Bearer $TOKEN"

echo ""
echo ""
echo "=========================================="
echo "Test Completed!"
echo "=========================================="

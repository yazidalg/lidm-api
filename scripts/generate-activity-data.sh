#!/bin/bash

# Script untuk generate activity data
BASE_URL="http://localhost:3000"

echo "🚀 Generating Activity Data for LIDM..."

# Step 1: Register beberapa user test
echo "📝 Creating test users..."

curl -s -X POST $BASE_URL/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Rina Septiani",
    "email": "rina@test.com",
    "password": "test123",
    "class": "XII RPL 1",
    "role": "user"
  }' > /dev/null

curl -s -X POST $BASE_URL/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Budi Prakoso",
    "email": "budi@test.com", 
    "password": "test123",
    "class": "XII RPL 2",
    "role": "user"
  }' > /dev/null

curl -s -X POST $BASE_URL/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Sari Wulandari",
    "email": "sari@test.com",
    "password": "test123", 
    "class": "XII RPL 1",
    "role": "user"
  }' > /dev/null

echo "✅ Test users created!"

# Step 2: Login users dan generate activities
users=("rina@test.com" "budi@test.com" "sari@test.com")

for email in "${users[@]}"; do
    echo "🔐 Logging in $email..."
    
    # Login user
    response=$(curl -s -X POST $BASE_URL/auth/login \
        -H "Content-Type: application/json" \
        -d "{\"email\": \"$email\", \"password\": \"test123\"}")
    
    # Extract token (assuming the response has "token" field)
    token=$(echo $response | grep -o '"token":"[^"]*' | grep -o '[^"]*$')
    
    if [ ! -z "$token" ]; then
        echo "✅ Login successful for $email"
        
        # Generate multiple activities for this user
        echo "📊 Generating activities for $email..."
        
        # Multiple lesson views (simulate active learning)
        for i in {1..5}; do
            curl -s -X GET $BASE_URL/lesson/1 \
                -H "Authorization: Bearer $token" > /dev/null
            sleep 1
        done
        
        # Multiple module views
        for i in {1..3}; do
            curl -s -X GET $BASE_URL/module/1 \
                -H "Authorization: Bearer $token" > /dev/null
            sleep 1
        done
        
        # Profile updates
        curl -s -X PUT $BASE_URL/user \
            -H "Authorization: Bearer $token" \
            -H "Content-Type: application/json" \
            -d "{\"name\": \"Updated $(echo $email | cut -d'@' -f1)\"}" > /dev/null
        
        # Logout
        curl -s -X POST $BASE_URL/auth/logout \
            -H "Authorization: Bearer $token" > /dev/null
        
        echo "✅ Activities generated for $email"
    else
        echo "❌ Failed to login $email"
    fi
    
    sleep 2
done

echo ""
echo "🎉 Activity data generation completed!"
echo ""
echo "📊 Now you can check the most active users:"
echo "   GET $BASE_URL/activity/most-active"
echo ""
echo "🔑 Login as admin first:"
echo "   POST $BASE_URL/auth/login"
echo '   {"email": "admin@lidm.com", "password": "admin123"}'

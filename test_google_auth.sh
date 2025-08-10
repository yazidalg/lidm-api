#!/bin/bash

# Test Google Auth Script
# Usage: ./test_google_auth.sh [ID_TOKEN]

echo "🚀 LIDM Google Auth Tester"
echo "=========================="

# Check if ID token is provided
if [ -z "$1" ]; then
    echo "❌ Error: Google ID token required"
    echo ""
    echo "Usage: ./test_google_auth.sh [GOOGLE_ID_TOKEN]"
    echo ""
    echo "How to get Google ID Token:"
    echo "1. Go to: https://developers.google.com/oauthplayground/"
    echo "2. Select 'Google OAuth2 API v2'"
    echo "3. Select 'https://www.googleapis.com/auth/userinfo.profile'"
    echo "4. Click 'Authorize APIs'"
    echo "5. Login with Google"
    echo "6. Click 'Exchange authorization code for tokens'"
    echo "7. Copy the 'id_token' value"
    echo ""
    exit 1
fi

ID_TOKEN="$1"
API_URL="http://localhost:3000/auth/google"

echo "📡 Testing Google Auth endpoint..."
echo "URL: $API_URL"
echo ""

# Test the Google auth endpoint
RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}\n" -X POST "$API_URL" \
    -H "Content-Type: application/json" \
    -d "{\"id_token\": \"$ID_TOKEN\"}")

# Extract HTTP code and response body
HTTP_CODE=$(echo "$RESPONSE" | grep "HTTP_CODE:" | cut -d: -f2)
BODY=$(echo "$RESPONSE" | grep -v "HTTP_CODE:")

echo "📊 Response:"
echo "HTTP Status: $HTTP_CODE"
echo "Body:"
echo "$BODY" | jq . 2>/dev/null || echo "$BODY"

if [ "$HTTP_CODE" = "200" ]; then
    echo ""
    echo "✅ Google Auth test PASSED!"
    
    # Extract JWT token if available
    JWT_TOKEN=$(echo "$BODY" | jq -r '.token // empty' 2>/dev/null)
    if [ ! -z "$JWT_TOKEN" ] && [ "$JWT_TOKEN" != "null" ]; then
        echo "🔑 JWT Token: $JWT_TOKEN"
        echo ""
        echo "💡 You can now use this JWT token for authenticated requests:"
        echo "curl -H \"Authorization: Bearer $JWT_TOKEN\" http://localhost:3000/user/profile"
    fi
else
    echo ""
    echo "❌ Google Auth test FAILED!"
fi

echo ""
echo "🔧 Troubleshooting:"
echo "- Make sure server is running on port 3000"
echo "- Check if GOOGLE_CLIENT_ID is set in .env"
echo "- Verify Google ID token is valid and not expired"
echo "- Check Google Cloud Console credentials settings"

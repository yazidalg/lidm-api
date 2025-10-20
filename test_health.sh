#!/bin/bash

# Health Check Test Script for LIDM Backend
# This script tests all health check endpoints

# Configuration
BASE_URL="http://localhost:8080"  # Change this to your actual URL
# For Cloud Run: BASE_URL="https://your-service-url"

echo "🔍 Testing Health Check Endpoints for LIDM Backend"
echo "📍 Base URL: $BASE_URL"
echo ""

# Function to test endpoint
test_endpoint() {
    local endpoint=$1
    local description=$2
    
    echo "🧪 Testing $description..."
    echo "   URL: $BASE_URL$endpoint"
    
    response=$(curl -s -w "\nHTTP_CODE:%{http_code}\nTIME:%{time_total}" "$BASE_URL$endpoint")
    
    # Extract HTTP code and time
    http_code=$(echo "$response" | grep "HTTP_CODE:" | cut -d: -f2)
    time_total=$(echo "$response" | grep "TIME:" | cut -d: -f2)
    json_response=$(echo "$response" | sed '/HTTP_CODE:/d' | sed '/TIME:/d')
    
    echo "   Status: $http_code"
    echo "   Time: ${time_total}s"
    
    if [ "$http_code" = "200" ]; then
        echo "   ✅ SUCCESS"
    elif [ "$http_code" = "503" ]; then
        echo "   ⚠️  SERVICE UNAVAILABLE (Database or dependency issue)"
    else
        echo "   ❌ FAILED"
    fi
    
    # Pretty print JSON response
    echo "$json_response" | python3 -m json.tool 2>/dev/null || echo "$json_response"
    echo ""
}

# Test all endpoints
test_endpoint "/health" "Basic Health Check"
test_endpoint "/ready" "Readiness Check (Database + Environment)"
test_endpoint "/healthy" "Liveness Check (Simple)"

echo "🏁 Health check tests completed!"
echo ""
echo "📋 Expected Results:"
echo "   /health   - Should return 200 with basic status and uptime"
echo "   /ready    - Should return 200 if database is connected, 503 if not"
echo "   /healthy  - Should return 200 if app is running, 503 if database issues"
echo ""
echo "💡 Troubleshooting:"
echo "   - If all return 503: Check database connection and environment variables"
echo "   - If /health works but others fail: Database connection issue"
echo "   - If none work: Application not running or wrong URL"

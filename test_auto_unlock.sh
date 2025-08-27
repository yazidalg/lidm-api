#!/bin/bash

# Test Auto-Unlock Functionality
# This script demonstrates how modules are automatically unlocked when all prequizzes are answered

echo "🚀 Testing Auto-Unlock Module Functionality"
echo "============================================="

# Set base URL
BASE_URL="http://localhost:3000"

# Set authentication token (replace with valid token)
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3RAZXhhbXBsZS5jb20iLCJ1c2VyX2lkIjoxLCJleHAiOjE3MzUzNzc0NTR9.8v5lCeUvNWVs6N3qKl7CbgmnZr5jdOLsE-_7zGzN9qQ"

echo ""
echo "📋 Step 1: Check initial module status"
curl -s -H "Authorization: Bearer $TOKEN" \
     -H "Content-Type: application/json" \
     "$BASE_URL/module/all" | jq -r '.data[] | "Module \(.id): \(.title) - Unlocked: \(.is_unlocked) - Complete: \(.is_complete)"'

echo ""
echo "📝 Step 2: Get prequizzes for Module 1"
MODULE_1_PREQUIZZES=$(curl -s -H "Authorization: Bearer $TOKEN" \
                      -H "Content-Type: application/json" \
                      "$BASE_URL/prequiz/module/1")

echo "Prequizzes found:"
echo "$MODULE_1_PREQUIZZES" | jq -r '.data[] | "Prequiz \(.id): \(.question)"'

echo ""
echo "📚 Step 3: Answer all prequizzes in Module 1"

# Get prequiz IDs from the response
PREQUIZ_IDS=$(echo "$MODULE_1_PREQUIZZES" | jq -r '.data[].id')

for prequiz_id in $PREQUIZ_IDS; do
    echo "Answering prequiz $prequiz_id..."
    
    # Get the prequiz details to find the correct answer
    PREQUIZ_DETAILS=$(curl -s -H "Authorization: Bearer $TOKEN" \
                      -H "Content-Type: application/json" \
                      "$BASE_URL/prequiz/$prequiz_id")
    
    CORRECT_ANSWER=$(echo "$PREQUIZ_DETAILS" | jq -r '.data.correct_answer')
    
    # Submit the correct answer
    RESPONSE=$(curl -s -H "Authorization: Bearer $TOKEN" \
               -H "Content-Type: application/json" \
               -X POST \
               -d "{\"prequiz_id\": $prequiz_id, \"selected_answer\": \"$CORRECT_ANSWER\"}" \
               "$BASE_URL/prequiz/submit")
    
    echo "$RESPONSE" | jq -r '"Answer submitted for prequiz \(.data.prequiz_id): \(.data.is_correct)"'
done

echo ""
echo "⏳ Step 4: Wait a moment for auto-unlock processing..."
sleep 2

echo ""
echo "🔓 Step 5: Check module status after completing Module 1"
curl -s -H "Authorization: Bearer $TOKEN" \
     -H "Content-Type: application/json" \
     "$BASE_URL/module/all" | jq -r '.data[] | "Module \(.id): \(.title) - Unlocked: \(.is_unlocked) - Complete: \(.is_complete) - Progress: \(.progress)%"'

echo ""
echo "✅ Test completed! Module 2 should now be unlocked if all Module 1 prequizzes were answered correctly."
echo ""
echo "💡 Expected result: Module 1 should show as complete (100% progress) and Module 2 should be unlocked."

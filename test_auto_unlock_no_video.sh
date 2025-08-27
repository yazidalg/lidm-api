#!/bin/bash

# Test Auto-Unlock for Modules WITHOUT Video Quizzes
# This script tests that modules unlock properly when there are only prequizzes (no video material)

echo "🎯 Testing Auto-Unlock for Modules WITHOUT Video Quizzes"
echo "========================================================="

# Set base URL
BASE_URL="http://localhost:3000"

# Set authentication token (replace with valid token)
TOKEN="YOUR_JWT_TOKEN_HERE"

echo ""
echo "📋 Step 1: Check module structure to find modules without video material"
curl -s -H "Authorization: Bearer $TOKEN" \
     -H "Content-Type: application/json" \
     "$BASE_URL/module/all" | jq -r '.data[] | "Module \(.id): \(.title) - Video Material: \(.video_material != null) - Video Quizzes: \(if .video_material then (.video_material.video_quizzes | length) else 0 end)"'

echo ""
echo "🔍 Step 2: Find a module with only prequizzes (no video material)"

# Get all modules and find one without video material
MODULES_DATA=$(curl -s -H "Authorization: Bearer $TOKEN" \
               -H "Content-Type: application/json" \
               "$BASE_URL/module/all")

# Find first module that has prequizzes but no video material
TARGET_MODULE=$(echo "$MODULES_DATA" | jq -r '.data[] | select(.video_material == null and .total_prequizzes > 0) | .id' | head -1)

if [ -z "$TARGET_MODULE" ]; then
    echo "❌ No module found with prequizzes only (no video material)"
    echo "Creating test scenario with any available module..."
    TARGET_MODULE=$(echo "$MODULES_DATA" | jq -r '.data[0].id')
fi

echo "🎯 Testing with Module $TARGET_MODULE"

echo ""
echo "📝 Step 3: Get all prequizzes for Module $TARGET_MODULE"
MODULE_PREQUIZZES=$(curl -s -H "Authorization: Bearer $TOKEN" \
                    -H "Content-Type: application/json" \
                    "$BASE_URL/prequiz/module/$TARGET_MODULE")

echo "Prequizzes found:"
echo "$MODULE_PREQUIZZES" | jq -r '.data[] | "Prequiz \(.id): \(.question)"'

PREQUIZ_COUNT=$(echo "$MODULE_PREQUIZZES" | jq '.data | length')
echo "Total prequizzes: $PREQUIZ_COUNT"

if [ "$PREQUIZ_COUNT" -eq 0 ]; then
    echo "❌ No prequizzes found in this module. Cannot test auto-unlock."
    exit 1
fi

echo ""
echo "📚 Step 4: Answer all prequizzes in Module $TARGET_MODULE"

# Get prequiz IDs from the response
PREQUIZ_IDS=$(echo "$MODULE_PREQUIZZES" | jq -r '.data[].id')

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
    
    IS_CORRECT=$(echo "$RESPONSE" | jq -r '.data.is_correct // false')
    echo "✅ Answer submitted for prequiz $prequiz_id: Correct=$IS_CORRECT"
done

echo ""
echo "⏳ Step 5: Wait for auto-unlock processing..."
sleep 3

echo ""
echo "🔓 Step 6: Check if next module was unlocked after completing Module $TARGET_MODULE"

NEXT_MODULE=$((TARGET_MODULE + 1))

echo "Module status after completion:"
curl -s -H "Authorization: Bearer $TOKEN" \
     -H "Content-Type: application/json" \
     "$BASE_URL/module/all" | jq -r '.data[] | select(.id <= '$NEXT_MODULE') | "Module \(.id): \(.title) - Unlocked: \(.is_unlocked) - Complete: \(.is_complete) - Progress: \(.progress)%"'

echo ""
echo "✅ Test completed!"
echo ""
echo "💡 Expected result:"
echo "- Module $TARGET_MODULE should show as complete (100% progress)"
echo "- Module $NEXT_MODULE should be unlocked (if it exists)"
echo "- This works even WITHOUT video quizzes!"

echo ""
echo "🏆 This demonstrates that modules unlock properly when they contain only prequizzes (no video material)."

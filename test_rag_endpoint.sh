#!/bin/bash

echo "🚀 Testing RAG Endpoint for User Activity (No Auth Required)"
echo "==========================================================="

echo "📊 Testing RAG Endpoint (No Authentication Required)..."
echo "🎯 Endpoint: GET /user-activity/for-rag"

RAG_RESPONSE=$(curl -s -X GET "http://localhost:3000/user-activity/for-rag?limit=10" \
  -H "Content-Type: application/json")

echo ""
echo "📈 Statistik Aktivitas:"
echo "$RAG_RESPONSE" | jq '.data.statistics'

echo ""
echo "🎯 Contoh Aktivitas dengan Enhanced Metadata:"
echo "$RAG_RESPONSE" | jq '.data.activities[0:2] | .[] | {
  id, 
  activity_type, 
  description, 
  timestamp,
  time_period,
  day_of_week,
  is_learning_activity,
  learning_category,
  user_intent,
  session_context,
  engagement_type,
  content_type,
  metadata: .metadata | {action, learning_activity, knowledge_area, user_intent, session_context}
}'

echo ""
echo "📋 RAG Context Information:"
echo "$RAG_RESPONSE" | jq '.data.rag_context'

echo ""
echo "🔍 Total Activities Found:"
echo "$RAG_RESPONSE" | jq '.data.statistics.total_activities'

echo ""
echo "📚 Learning Activity Breakdown:"
echo "$RAG_RESPONSE" | jq '.data.statistics.activity_breakdown'

echo ""
echo "✅ RAG Endpoint Test Completed!"
echo "🎉 Data siap untuk digunakan oleh AI/Knowledge System"
echo "⭐ Endpoint ini tidak memerlukan authentication token"

# Optional: Test dengan filter user_id
echo ""
echo "🔍 Testing with specific user_id filter..."
RAG_USER_RESPONSE=$(curl -s -X GET "http://localhost:3000/user-activity/for-rag?limit=5&user_id=11" \
  -H "Content-Type: application/json")

echo "📊 Activities for specific user:"
echo "$RAG_USER_RESPONSE" | jq '.data.statistics'

#!/bin/bash

echo "🔍 Melihat Data dari Endpoint /lesson/all untuk RAG"
echo "================================================="

# Login
echo "1. Login untuk mendapatkan token..."
TOKEN=$(curl -s -X POST http://localhost:3000/auth/belajar-login \
  -H "Content-Type: application/json" \
  -d '{"email": "azis@belajar.id", "password": "password123"}' | \
  grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
    echo "❌ Login gagal! Pastikan server berjalan dan credentials benar"
    exit 1
fi

echo "✅ Login berhasil!"

# Get lesson data
echo ""
echo "2. 📚 Data dari /lesson/all:"
echo "==========================="
curl -s -X GET "http://localhost:3000/lesson/all" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" | jq '.data[] | {
    id,
    title,
    content: (.content | if length > 100 then .[0:100] + "..." else . end),
    module: {
      id: .Module.ID,
      title: .Module.title,
      description: (.Module.description | if length > 50 then .[0:50] + "..." else . end)
    },
    sort_order: .SortOrder,
    created_at: .CreatedAt
  }'

echo ""
echo "3. 📊 Data dari /module/all:"
echo "============================"
curl -s -X GET "http://localhost:3000/module/all" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" | jq '.data[] | {
    id: .ID,
    title,
    description: (.description | if length > 100 then .[0:100] + "..." else . end),
    icon,
    thumbnail,
    lessons_count: (.lessons | length),
    lessons_titles: [.lessons[]?.title]
  }'

echo ""
echo "4. 🎯 Contoh Aktivitas yang Terekam untuk RAG:"
echo "=============================================="
curl -s -X GET "http://localhost:3000/user-activity/for-rag?limit=5" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" | jq '.data.activities[] | select(.activity_type == "lihat_pelajaran") | {
    activity_type,
    description,
    timestamp,
    learning_category,
    user_intent,
    session_context,
    metadata: .metadata | {action, content_type, knowledge_area, learning_activity}
  }'

echo ""
echo "✅ Data lesson dan module sudah siap untuk RAG system!"
echo "🔗 Data ini mencakup:"
echo "   - Daftar semua lesson dengan title, content, dan module info"
echo "   - Module information dengan lessons yang ada di dalamnya"  
echo "   - Activity tracking dengan enhanced metadata untuk AI analysis"

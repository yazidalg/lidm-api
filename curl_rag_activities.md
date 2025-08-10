# Curl untuk Endpoint RAG User Activity

## 1. Endpoint RAG - Data Aktivitas untuk AI/Knowledge System (Tanpa Auth)

⭐ **Endpoint ini tidak memerlukan authentication token** - dirancang khusus untuk sistem AI/RAG

### a. Mendapatkan Semua Aktivitas untuk RAG
```bash
curl -X GET "http://localhost:3000/user-activity/for-rag?limit=50" \
  -H "Content-Type: application/json" | jq '.'
```

### b. Mendapatkan Aktivitas User Tertentu untuk RAG
```bash
curl -X GET "http://localhost:3000/user-activity/for-rag?limit=100&user_id=11" \
  -H "Content-Type: application/json" | jq '.'
```

## 2. Login untuk Endpoints Lain (Optional)

Jika kamu butuh mengakses endpoint lain yang memerlukan auth, gunakan login ini:

```bash
# Login dengan akun belajar
curl -X POST http://localhost:3000/auth/belajar-login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "azis@belajar.id",
    "password": "password123"
  }'
```

## 3. Contoh Response RAG Endpoint:

```json
{
  "message": "Data aktivitas untuk RAG berhasil diambil",
  "data": {
    "activities": [
      {
        "id": 25,
        "user_id": 11,
        "activity_type": "lihat_pelajaran",
        "description": "Melihat daftar semua pelajaran",
        "timestamp": "2025-08-04T10:30:00Z",
        "date": "2025-08-04",
        "time": "10:30:00",
        "time_period": "morning",
        "day_of_week": "Sunday",
        "is_weekend": true,
        "is_learning_activity": true,
        "content_type": "lesson_list",
        "user_intent": "browse_available_lessons",
        "session_context": "lesson_discovery",
        "engagement_type": "content_consumption",
        "learning_category": "content_consumption",
        "knowledge_acquisition": true,
        "metadata": {
          "path": "/lesson/all",
          "method": "GET",
          "response_status": 200,
          "action": "view_all_lessons",
          "content_type": "lesson_list",
          "learning_activity": true,
          "knowledge_area": "lessons",
          "user_intent": "browse_available_lessons",
          "session_context": "lesson_discovery"
        }
      },
      {
        "id": 24,
        "user_id": 11,
        "activity_type": "selesai_pelajaran",
        "description": "Menyelesaikan pelajaran",
        "timestamp": "2025-08-04T09:15:00Z",
        "date": "2025-08-04",
        "time": "09:15:00",
        "time_period": "morning",
        "day_of_week": "Sunday",
        "is_weekend": true,
        "is_learning_activity": true,
        "content_type": "lesson_completion",
        "user_intent": "complete_learning_objective",
        "session_context": "lesson_completion",
        "engagement_type": "achievement",
        "learning_category": "content_completion",
        "achievement": true,
        "metadata": {
          "lesson_id": "5",
          "action": "complete_lesson",
          "content_type": "lesson_completion",
          "learning_activity": true,
          "achievement": true,
          "completion_timestamp": "2025-08-04T09:15:00Z",
          "user_intent": "complete_learning_objective",
          "session_context": "lesson_completion",
          "engagement_type": "achievement",
          "learning_milestone": true
        }
      }
    ],
    "statistics": {
      "total_activities": 15,
      "time_range": {
        "start": "2025-08-03T08:00:00Z",
        "end": "2025-08-04T10:30:00Z"
      },
      "activity_breakdown": {
        "lihat_pelajaran": 8,
        "selesai_pelajaran": 3,
        "lihat_modul": 2,
        "masuk": 2
      },
      "learning_activity_count": 13,
      "learning_percentage": 86.67
    },
    "rag_context": {
      "purpose": "learning_behavior_analysis",
      "data_enrichment": "enhanced_metadata_for_ai",
      "temporal_analysis": true,
      "learning_pattern_detection": true
    }
  }
}
```

## 4. Data yang Tersedia untuk RAG/AI:

### Learning Context:
- **learning_category**: Kategori pembelajaran (content_consumption, content_completion, curriculum_exploration, knowledge_assessment)
- **user_intent**: Intensi pengguna (study_lesson, complete_learning_objective, test_knowledge, dll)
- **session_context**: Konteks sesi (lesson_discovery, active_learning, assessment, dll)
- **engagement_type**: Tipe keterlibatan (content_consumption, achievement, knowledge_evaluation, dll)

### Temporal Analysis:
- **time_period**: Periode waktu (morning, afternoon, evening, night)
- **day_of_week**: Hari dalam minggu
- **is_weekend**: Boolean apakah weekend
- **timestamp**: Timestamp lengkap aktivitas

### Learning Patterns:
- **knowledge_acquisition**: Boolean untuk aktivitas akuisisi pengetahuan
- **achievement**: Boolean untuk pencapaian/completion
- **skill_evaluation**: Boolean untuk evaluasi skill (quiz)
- **structured_learning**: Boolean untuk pembelajaran terstruktur

### Content Information:
- **content_type**: Tipe konten (lesson_list, lesson_detail, module_list, quiz_participation, dll)
- **knowledge_area**: Area pengetahuan yang diakses
- **learning_milestone**: Boolean untuk milestone pembelajaran

## 5. Penggunaan untuk RAG/AI:

### Analisis Pola Belajar:
- Waktu belajar favorit user
- Jenis konten yang paling sering diakses
- Tingkat completion rate
- Pola engagement mingguan/harian

### Rekomendasi Pembelajaran:
- Berdasarkan aktivitas sebelumnya
- Waktu optimal untuk notifikasi
- Konten yang relevan dengan minat user

### Personalisasi:
- Gaya belajar user (visual, assessment-focused, completion-oriented)
- Preferensi konten dan timing
- Level engagement dan motivasi

## 6. Script Test RAG Endpoint:

```bash
#!/bin/bash

# Login
TOKEN=$(curl -s -X POST http://localhost:3000/auth/belajar-login \
  -H "Content-Type: application/json" \
  -d '{"email": "azis@belajar.id", "password": "password123"}' | \
  grep -o '"token":"[^"]*"' | cut -d'"' -f4)

echo "Token: $TOKEN"

# Generate some learning activities
echo "Generating learning activities..."
curl -s -X GET "http://localhost:3000/lesson/all" \
  -H "Authorization: Bearer $TOKEN" > /dev/null

curl -s -X GET "http://localhost:3000/module/all" \
  -H "Authorization: Bearer $TOKEN" > /dev/null

# Get RAG data
echo "Getting RAG data..."
curl -X GET "http://localhost:3000/user-activity/for-rag?limit=20" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" | jq '.data.activities[0:3]'

echo "Learning statistics:"
curl -X GET "http://localhost:3000/user-activity/for-rag?limit=20" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" | jq '.data.statistics'
```

## 7. Metadata yang Diperkaya untuk RAG:

Setiap aktivitas sekarang menyimpan metadata yang kaya seperti:
- **action**: Tindakan spesifik yang dilakukan
- **learning_activity**: Boolean marker untuk aktivitas pembelajaran
- **user_intent**: Intensi user yang terdeteksi
- **session_context**: Konteks sesi pembelajaran
- **engagement_type**: Tipe keterlibatan user
- **knowledge_area**: Area pengetahuan yang diakses
- **learning_milestone**: Marker untuk pencapaian pembelajaran
- **completion_timestamp**: Waktu completion untuk tracking progress

Data ini sangat berguna untuk AI/RAG system dalam memahami pola belajar user dan memberikan insights yang personal!

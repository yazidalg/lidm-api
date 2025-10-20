# 📝 Panduan Cara Menjawab Quiz

Dokumentasi lengkap untuk menjawab **Video Quiz** dan **Module Quiz** di LIDM Backend.

---

## 🎥 **1. VIDEO QUIZ**

Video Quiz adalah quiz yang muncul di tengah video pembelajaran pada timestamp tertentu.

### **Endpoint**
```
POST /video-quiz/submit
```

### **Request Body**
```json
{
  "video_quiz_id": 1,           // ID dari video quiz yang akan dijawab
  "selected_answer": "A",       // Pilihan jawaban: A, B, C, atau D
  "response_time": 15           // Waktu respons dalam detik
}
```

### **Headers**
```
Content-Type: application/json
Authorization: Bearer {JWT_TOKEN}
```

### **Contoh cURL**

#### ✅ Jawaban Benar
```bash
curl -X POST "http://localhost:8080/video-quiz/submit" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "video_quiz_id": 1,
    "selected_answer": "A",
    "response_time": 15
  }'
```

#### ❌ Jawaban Salah (untuk testing)
```bash
curl -X POST "http://localhost:8080/video-quiz/submit" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "video_quiz_id": 1,
    "selected_answer": "B",
    "response_time": 20
  }'
```

### **Response Success**
```json
{
  "message": "Answer submitted successfully",
  "data": {
    "id": 123,
    "user_id": 1,
    "video_quiz_id": 1,
    "selected_answer": "A",
    "is_correct": true,
    "response_time": 15,
    "submitted_at": "2025-10-18T10:30:00Z"
  }
}
```

### **Cara Mendapatkan Video Quiz**
```bash
# Dapatkan semua video quiz berdasarkan video material
curl -X GET "http://localhost:8080/video-quiz/video-material/3" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Dapatkan video quiz by ID
curl -X GET "http://localhost:8080/video-quiz/1" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### **Lihat Jawaban User**
```bash
# Semua jawaban user
curl -X GET "http://localhost:8080/video-quiz/user-answers" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Jawaban untuk video material tertentu
curl -X GET "http://localhost:8080/video-quiz/user-answers/3" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

## 📚 **2. MODULE QUIZ (Quiz Session)**

Module Quiz adalah quiz yang diambil setelah menyelesaikan modul pembelajaran.

### **Flow Quiz Session**
1. **Create Quiz Session** → Membuat sesi quiz baru
2. **Join Quiz** (opsional untuk multiplayer) → Bergabung dengan kode invite
3. **Answer Question** → Menjawab pertanyaan satu per satu
4. **Finish Quiz** → Menyelesaikan quiz
5. **Get Results** → Melihat hasil quiz

---

### **Step 1: Create Quiz Session**

#### **Endpoint**
```
POST /quiz-sessions/
```

#### **Request Body**
```json
{
  "mode": "single_player",      // "single_player" atau "multiplayer"
  "module_id": 1                // ID modul yang akan di-quiz
}
```

#### **Contoh cURL**
```bash
curl -X POST "http://localhost:8080/quiz-sessions/" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "mode": "single_player",
    "module_id": 1
  }'
```

#### **Response**
```json
{
  "message": "Quiz session created successfully",
  "data": {
    "quiz_id": 123,
    "mode": "single_player",
    "module_id": 1,
    "status": "waiting",
    "invite_code": "ABC12345",
    "created_at": "2025-10-18T10:30:00Z"
  }
}
```

---

### **Step 2: Join Quiz (Multiplayer Only)**

#### **Endpoint**
```
POST /quiz-sessions/join
```

#### **Request Body**
```json
{
  "invite_code": "ABC12345"     // Kode invite 8 karakter
}
```

#### **Contoh cURL**
```bash
curl -X POST "http://localhost:8080/quiz-sessions/join" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "invite_code": "ABC12345"
  }'
```

---

### **Step 3: Answer Question** ⭐ **PENTING**

#### **Endpoint**
```
POST /quiz-sessions/answer
```

#### **Request Body**
```json
{
  "question_id": 1,             // ID pertanyaan yang dijawab
  "user_answer": "A",           // Pilihan jawaban: A, B, C, atau D
  "response_time": 5000         // Waktu respons dalam MILIDETIK (ms)
}
```

#### **Contoh cURL**
```bash
# Menjawab pertanyaan pertama
curl -X POST "http://localhost:8080/quiz-sessions/answer" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "question_id": 1,
    "user_answer": "A",
    "response_time": 5000
  }'

# Menjawab pertanyaan kedua
curl -X POST "http://localhost:8080/quiz-sessions/answer" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "question_id": 2,
    "user_answer": "C",
    "response_time": 7500
  }'
```

#### **Response**
```json
{
  "message": "Answer submitted successfully",
  "data": {
    "quiz_session": {
      "quiz_id": 123,
      "status": "in_progress",
      "current_question": 2,
      "total_questions": 10
    },
    "is_correct": true,
    "points_earned": 100
  }
}
```

---

### **Step 4: Finish Quiz**

#### **Endpoint**
```
POST /quiz-sessions/:quiz_id/finish
```

#### **Contoh cURL**
```bash
curl -X POST "http://localhost:8080/quiz-sessions/123/finish" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

### **Step 5: Get Quiz Results**

#### **Endpoint**
```
GET /quiz-sessions/:quiz_id/results
```

#### **Contoh cURL**
```bash
curl -X GET "http://localhost:8080/quiz-sessions/123/results" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### **Response**
```json
{
  "message": "Quiz results retrieved successfully",
  "data": {
    "quiz_id": 123,
    "user_id": 1,
    "module_id": 1,
    "total_questions": 10,
    "correct_answers": 8,
    "wrong_answers": 2,
    "score": 80,
    "total_time": 120,
    "status": "completed",
    "completed_at": "2025-10-18T10:35:00Z"
  }
}
```

---

## 🔑 **Perbedaan Utama**

| Aspek | Video Quiz | Module Quiz |
|-------|------------|-------------|
| **Endpoint Submit** | `/video-quiz/submit` | `/quiz-sessions/answer` |
| **ID Quiz** | `video_quiz_id` | `question_id` |
| **Response Time** | Dalam **detik** | Dalam **milidetik** (ms) |
| **Flow** | Langsung submit | Create → Answer → Finish |
| **Context** | Saat menonton video | Setelah selesai modul |

---

## 📋 **Checklist untuk Frontend**

### **Video Quiz:**
- [ ] Dapatkan `video_quiz_id` dari video material
- [ ] Tampilkan quiz pada `timestamp_start`
- [ ] Hitung `response_time` dalam **detik**
- [ ] Submit jawaban ke `/video-quiz/submit`
- [ ] Tampilkan feedback (benar/salah)

### **Module Quiz:**
- [ ] Create quiz session dengan `/quiz-sessions/`
- [ ] Simpan `quiz_id` dari response
- [ ] Loop untuk setiap pertanyaan:
  - Tampilkan soal
  - Hitung `response_time` dalam **milidetik**
  - Submit ke `/quiz-sessions/answer`
- [ ] Finish quiz dengan `/quiz-sessions/:quiz_id/finish`
- [ ] Tampilkan hasil dari `/quiz-sessions/:quiz_id/results`

---

## 🎯 **Tips Debugging**

### **Masalah Umum:**

1. **401 Unauthorized** → Token JWT expired atau tidak valid
   ```bash
   # Login ulang untuk dapat token baru
   curl -X POST "http://localhost:8080/auth/login" \
     -H "Content-Type: application/json" \
     -d '{"Email": "user@example.com", "Password": "password123"}'
   ```

2. **400 Bad Request** → Validasi gagal
   - Cek format jawaban (harus A, B, C, atau D)
   - Pastikan `response_time` > 0
   - Untuk module quiz: gunakan **milidetik** bukan detik

3. **404 Not Found** → ID tidak ditemukan
   - Pastikan `video_quiz_id` atau `question_id` valid
   - Cek dengan GET endpoint terlebih dahulu

---

## 📞 **Butuh Bantuan?**

- Lihat log server di terminal
- Cek file: `server.log` untuk detail error
- Test dengan script: `test_*.sh` yang ada di root project

**Happy Coding! 🚀**

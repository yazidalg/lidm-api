# 📚 Panduan Presentasi LIDM - Hal yang Perlu Diketahui

## 🎯 Overview Aplikasi

Aplikasi LIDM (Learning Interactive Digital Media) adalah platform pembelajaran interaktif untuk siswa SD kelas 4 yang fokus pada materi **Fotosintesis**. Aplikasi terdiri dari 2 komponen utama:

1. **lidm-api** (Backend API) - Go/Gin framework dengan MySQL database
2. **lidm-ai** (AI Service) - Flask/Python service untuk chat AI dengan Groq API

---

## 🏗️ Arsitektur Sistem

### 1. **lidm-api** (Backend Utama)

**Teknologi:**
- **Language**: Go 1.24.4
- **Framework**: Gin (HTTP router)
- **Database**: MySQL (via GORM)
- **Real-time**: Socket.IO v2 untuk multiplayer quiz
- **Authentication**: JWT (JSON Web Token)
- **Deployment**: Docker + Cloud Run ready

**Struktur:**
```
lidm-api/
├── cmd/main.go              # Entry point aplikasi
├── internal/
│   ├── app/
│   │   ├── models/          # Database models (User, Quiz, Module, dll)
│   │   ├── handlers/        # HTTP request handlers
│   │   ├── services/        # Business logic
│   │   ├── repositories/    # Database operations
│   │   └── socket/          # Socket.IO event handlers
│   ├── routes/              # API route definitions
│   ├── config/              # Configuration & DB connection
│   └── realtime/socketio/   # Socket.IO server setup
```

### 2. **lidm-ai** (AI Chat Service)

**Teknologi:**
- **Language**: Python 3.11
- **Framework**: Flask
- **AI Model**: Groq API (Llama 4 Scout 17B)
- **Database**: MongoDB (untuk chat history)
- **Deployment**: Docker + Cloud Run ready

**Fitur:**
- Chat AI khusus untuk pembelajaran fotosintesis
- Multimodal support (text + image)
- Session management dengan MongoDB
- Upload file untuk analisis gambar

---

## 📊 Database Models Utama

### Core Entities:

1. **User**
   - Email, Name, Password
   - Role (User/Admin/Teacher)
   - Point, TotalXP, Lives (nyawa untuk single player)
   - CurrentStreak, MaxStreak (daily login streak)
   - Position tracking untuk leaderboard

2. **Module**
   - Title, Description, Thumbnail, Icon
   - OffsetX, OffsetY (untuk positioning di UI)
   - Relasi: Quizzes, VideoMaterial, ARExperiments, Prequizzes, Flashcards

3. **ModuleProgress**
   - Track progress per user per module
   - `is_unlocked`: Apakah module bisa diakses
   - `is_complete`: Apakah module sudah selesai (100%)
   - `progress`: Persentase completion (0-100)
   - Auto-unlock: Module 1 selalu unlocked, module berikutnya unlock otomatis saat module sebelumnya 100%

4. **Quiz & Question**
   - Quiz bisa solo atau multiplayer
   - Question dengan tipe: "regular" atau "hots" (Higher Order Thinking Skills)
   - Options (A, B, C, D), CorrectAnswer, Explanation
   - ReadTime & AnswerTime untuk timer

5. **Prequiz**
   - Quiz sebelum belajar materi (pre-assessment)
   - Terkait dengan Module
   - User answers tersimpan untuk progress tracking

6. **VideoQuiz**
   - Quiz setelah menonton video
   - Terkait dengan VideoMaterial
   - User answers untuk progress tracking

7. **Flashcard**
   - Sistem flashcard dengan FSRS algorithm (spaced repetition)
   - UserFlashcardProgress untuk tracking review
   - Interval: 1m, 5m, 7h, 10h, dll

8. **Participant**
   - User yang ikut quiz
   - TotalScore, IsFinished
   - Untuk leaderboard calculation

9. **Leaderboard**
   - Ranking berdasarkan total score dari semua quiz
   - Support filter by module_id dan quiz_type

10. **UserActivity**
    - Tracking aktivitas user (login, belajar, quiz, dll)
    - Support RAG (Retrieval Augmented Generation) untuk AI

---

## 🔑 Fitur-Fitur Utama

### 1. **Authentication & Authorization**
- Register/Login dengan email & password
- Google OAuth login
- JWT token untuk session
- Role-based access (User/Admin/Teacher)
- Forgot password dengan OTP via email

### 2. **Module System dengan Progressive Unlock**
- **Module 1 selalu unlocked** untuk semua user baru
- **Auto-unlock**: Saat user menyelesaikan 100% module (semua prequiz + video quiz), module berikutnya otomatis unlock
- **Progress tracking**: Real-time progress calculation
- Sequential unlocking: 1→2→3→4...

**Flow:**
```
User baru → Module 1 unlocked
User selesaikan semua prequiz + video quiz di Module 1 → Progress 100%
→ Module 2 otomatis unlocked
```

### 3. **Quiz System**

#### A. **Prequiz** (Pre-assessment)
- Quiz sebelum belajar materi
- Submit answer → Update progress
- Auto-unlock next module jika module selesai

#### B. **Video Quiz** (Post-video assessment)
- Quiz setelah menonton video
- Submit answer → Update progress
- Auto-unlock next module jika module selesai

#### C. **Solo Quiz** (Single Player)
- Quiz untuk latihan sendiri
- Lives system: User punya 5 nyawa
- Lives berkurang jika salah (hanya di solo mode)
- Lives reset harian

#### D. **Multiplayer Quiz** (Matchmaking)
- Real-time quiz dengan Socket.IO
- Matchmaking system: User mencari opponent
- Head-to-head competition
- Real-time scoring & leaderboard update

**Socket.IO Events:**
- `create_lobby` - Buat room quiz
- `join_lobby` - Join dengan invite code
- `start_quiz` - Mulai quiz
- `question` - Soal baru muncul
- `submit_answer` - Submit jawaban
- `answer_result` - Hasil jawaban (correct/incorrect)
- `quiz_completed` - Quiz selesai dengan final scores

### 4. **Flashcard System dengan FSRS**
- Spaced repetition algorithm (FSRS)
- Review flashcards berdasarkan interval
- Get due flashcards (yang sudah waktunya di-review)
- Retention statistics

### 5. **Leaderboard System**
- Ranking berdasarkan total score
- Top 3 ditampilkan terpisah (juara1, juara2, juara3)
- Rank 4+ di array leaderboard
- Filter by module_id dan quiz_type (solo/matchmaking)
- Position tracking: increasing/decreasing/stable

### 6. **User Activity Tracking**
- Track semua aktivitas user (login, belajar, quiz, dll)
- Support RAG endpoint untuk AI service
- Dashboard analytics

### 7. **AI Chat Service (lidm-ai)**
- Chat AI khusus fotosintesis
- Menggunakan Groq API (Llama 4 Scout 17B)
- Session management dengan MongoDB
- Multimodal: Text + Image support
- Upload file untuk analisis gambar

**Endpoints:**
- `GET /chat` - Init new chat session
- `POST /chat` - Send message
- `GET /chat/<session_id>` - Get chat history
- `POST /chat/<session_id>/regenerate` - Regenerate last response
- `POST /upload` - Upload image/file

**System Prompt:**
- AI bernama "Alsan"
- Fokus hanya pada fotosintesis
- Bahasa sederhana untuk SD kelas 4
- Mengurangi miskonsepsi siswa
- Emoticon friendly (tapi tidak julur lidah)

---

## 🔌 API Endpoints Penting

### Authentication
```
POST /auth/register      - Register user baru
POST /auth/login          - Login
POST /auth/google         - Google OAuth login
POST /auth/logout         - Logout
GET  /auth/verify/:token  - Verify email
```

### Module
```
GET  /module/all          - Get semua module dengan unlock status & progress
GET  /module/:id          - Get detail module
POST /module              - Create module (Admin/Teacher)
PUT  /module/:id          - Update module (Admin/Teacher)
DELETE /module/:id        - Delete module (Admin/Teacher)
```

### Prequiz
```
GET  /prequiz/submaterial/:sub_material_id  - Get prequizzes
POST /prequiz/submit                         - Submit answer
```

### Video Quiz
```
GET  /video-quiz/video-material/:id  - Get video quizzes
POST /video-quiz/submit               - Submit answer
```

### Quiz (Solo/Multiplayer)
```
POST /quiz                    - Create quiz
GET  /quiz/:id                - Get quiz detail
POST /quiz-session/           - Create quiz session
POST /quiz-session/join       - Join quiz session
POST /quiz-session/answer     - Submit answer
GET  /quiz-session/:id/results - Get results
```

### Leaderboard
```
GET /leaderboard              - Get leaderboard
GET /leaderboard/user/:id     - Get user rank
```

### Flashcard
```
GET  /flashcard/all                    - Get all flashcards
GET  /flashcard/due                    - Get due flashcards
POST /flashcard/:id/review             - Review flashcard
POST /flashcard/module/:id/initialize  - Initialize module flashcards
GET  /flashcard/stats                  - Get retention stats
```

### User
```
GET /user/profile  - Get current user profile
```

### Socket.IO
```
WS /socket.io      - Socket.IO connection
WS /ws/matchmaking - Matchmaking endpoint
WS /ws/prequiz     - Prequiz real-time
WS /ws/quiz-session/:id - Quiz session real-time
```

---

## 🎮 Flow Aplikasi

### Flow User Baru:
1. Register → Email verification
2. Login → Dapat JWT token
3. Module 1 otomatis unlocked
4. Mulai belajar:
   - Baca materi
   - Jawab prequiz
   - Tonton video
   - Jawab video quiz
   - Review flashcards
5. Progress 100% → Module 2 unlock
6. Repeat untuk module berikutnya

### Flow Quiz Multiplayer:
1. User A create lobby → Dapat invite code
2. User B join dengan invite code
3. Server match 2 players → Quiz start
4. Real-time questions via Socket.IO
5. Players submit answers
6. Real-time scoring
7. Quiz completed → Update leaderboard

### Flow AI Chat:
1. User init chat → Dapat session_id
2. User kirim pertanyaan tentang fotosintesis
3. AI (Alsan) jawab dengan bahasa sederhana
4. User bisa upload gambar untuk analisis
5. Chat history tersimpan di MongoDB

---

## 🔐 Security & Authentication

- **JWT Token**: Semua endpoint (kecuali public) require JWT
- **Role-based**: Admin/Teacher bisa CRUD content, User hanya read
- **Password**: Hashed dengan bcrypt
- **Email Verification**: Required untuk aktivasi akun
- **OTP**: Untuk forgot password (6 digit, expires 10 menit)

---

## 📈 Progress & Unlock System

### Progress Calculation:
```
Progress = (Completed Items / Total Items) × 100%

Completed Items = Prequizzes answered + Video quizzes answered
Total Items = Total prequizzes + Total video quizzes in module
```

### Unlock Logic:
- Module 1: Always unlocked (default)
- Module N+1: Unlocked jika Module N progress = 100%
- Auto-check setelah setiap prequiz/video quiz submission

---

## 🎯 Points & Scoring System

- **Quiz Points**: 10 points per correct answer
- **Total Score**: Sum dari semua quiz yang selesai
- **Leaderboard**: Ranking berdasarkan total score
- **XP System**: TotalXP untuk tracking overall progress
- **Lives**: 5 nyawa untuk solo quiz (reset harian)

---

## 🚀 Deployment

### lidm-api:
- Docker container
- Cloud Run ready (PORT env variable)
- Health check: `/health`, `/ready`, `/healthy`
- Database migration otomatis saat startup

### lidm-ai:
- Docker container
- Cloud Run ready (PORT=8080)
- MongoDB connection via MONGO_URI
- Groq API key via GROQ_API_KEY

---

## 💡 Hal Penting untuk Presentasi

### 1. **Value Proposition**
- Platform pembelajaran interaktif untuk SD kelas 4
- Fokus pada fotosintesis dengan pendekatan gamifikasi
- AI tutor (Alsan) untuk personalized learning
- Progressive unlock system untuk engagement
- Real-time multiplayer quiz untuk kompetisi

### 2. **Teknologi Stack**
- **Backend**: Go (performant, scalable)
- **AI Service**: Python Flask (flexible untuk AI integration)
- **Database**: MySQL (structured data) + MongoDB (chat history)
- **Real-time**: Socket.IO (low latency multiplayer)
- **AI Model**: Groq Llama 4 (fast inference)

### 3. **Fitur Unik**
- **Auto-unlock system**: Gamifikasi dengan progressive unlock
- **FSRS Flashcard**: Spaced repetition untuk retention
- **Real-time multiplayer**: Head-to-head quiz competition
- **AI Tutor**: Personalized learning dengan chat AI
- **Activity tracking**: Analytics untuk insights

### 4. **Scalability**
- Docker containerization
- Cloud Run ready (auto-scaling)
- Stateless API design
- Efficient database queries dengan indexing

### 5. **User Experience**
- Module 1 always accessible (no barrier to entry)
- Clear progress tracking (0-100%)
- Real-time feedback (quiz results, leaderboard)
- AI support untuk questions
- Gamification (points, lives, streaks)

---

## ❓ Pertanyaan yang Mungkin Ditanyakan Juri

### Q: Bagaimana sistem auto-unlock bekerja?
**A:** 
- Module 1 selalu unlocked untuk user baru
- Saat user menyelesaikan semua prequiz + video quiz di module (progress 100%), sistem otomatis unlock module berikutnya
- Progress dihitung real-time: (completed items / total items) × 100%
- Check dilakukan setiap kali user submit prequiz/video quiz answer

### Q: Bagaimana scoring system bekerja?
**A:**
- 10 points per correct answer di quiz
- Total score = sum dari semua quiz yang sudah selesai (is_finished = true)
- Leaderboard ranking berdasarkan total score (descending)
- User yang belum pernah ikut quiz = score 0

### Q: Bagaimana AI chat service terintegrasi?
**A:**
- Service terpisah (lidm-ai) menggunakan Flask + Groq API
- Fokus khusus pada fotosintesis dengan system prompt yang ketat
- Chat history tersimpan di MongoDB untuk context
- Support multimodal (text + image) untuk analisis gambar

### Q: Bagaimana real-time multiplayer quiz bekerja?
**A:**
- Menggunakan Socket.IO untuk real-time communication
- Matchmaking: User create lobby → dapat invite code → opponent join
- Server manage quiz state, questions, timers
- Real-time events: question, answer_result, quiz_completed
- Scoring real-time dan update leaderboard setelah selesai

### Q: Apa perbedaan solo quiz dan multiplayer quiz?
**A:**
- **Solo Quiz**: 
  - User latihan sendiri
  - Lives system (5 nyawa, berkurang jika salah)
  - Lives reset harian
- **Multiplayer Quiz**:
  - Head-to-head dengan opponent
  - No lives (kompetisi langsung)
  - Real-time scoring
  - Update leaderboard

### Q: Bagaimana flashcard system dengan FSRS bekerja?
**A:**
- FSRS (Free Spaced Repetition Scheduler) algorithm
- Interval review: 1m → 5m → 7h → 10h → dll (berdasarkan performance)
- User review flashcard dengan grade (1-4)
- System calculate next review date berdasarkan grade
- Get due flashcards untuk review yang sudah waktunya

### Q: Bagaimana progress tracking bekerja?
**A:**
- ModuleProgress table track per user per module
- Progress = (prequizzes answered + video quizzes answered) / (total prequizzes + total video quizzes) × 100%
- Auto-update setiap kali user submit prequiz/video quiz answer
- Real-time calculation, tidak perlu manual update

### Q: Apa yang terjadi jika user skip module?
**A:**
- User tidak bisa skip module karena sequential unlock
- Harus selesaikan Module 1 (100%) untuk unlock Module 2
- Ini untuk memastikan user memahami materi secara berurutan
- Gamification element untuk engagement

### Q: Bagaimana sistem handle concurrent users?
**A:**
- Stateless API design (JWT token)
- Database transactions untuk data consistency
- Socket.IO rooms untuk isolate multiplayer sessions
- Goroutines untuk async operations (video quiz progress update)

### Q: Apa kelebihan menggunakan Go untuk backend?
**A:**
- Performance tinggi (compiled language)
- Concurrent handling dengan goroutines
- Low memory footprint
- Type safety
- Good untuk real-time systems (Socket.IO)

---

## 📝 Checklist Sebelum Presentasi

- [ ] Pahami flow aplikasi dari register sampai selesai module
- [ ] Pahami perbedaan solo quiz vs multiplayer quiz
- [ ] Pahami auto-unlock system dan progress calculation
- [ ] Pahami AI chat service dan integrasinya
- [ ] Pahami leaderboard system dan scoring
- [ ] Pahami Socket.IO events untuk multiplayer
- [ ] Siapkan demo flow: Register → Login → Module 1 → Prequiz → Video Quiz → Module 2 unlock
- [ ] Siapkan demo multiplayer quiz
- [ ] Siapkan demo AI chat
- [ ] Siapkan jawaban untuk pertanyaan teknis (Go, Socket.IO, FSRS, dll)

---

## 🎓 Key Takeaways

1. **Progressive Learning**: Sequential module unlock memastikan user belajar step-by-step
2. **Gamification**: Points, lives, streaks, leaderboard untuk engagement
3. **Personalized AI**: Chat AI khusus fotosintesis dengan bahasa sederhana
4. **Real-time Competition**: Multiplayer quiz untuk motivasi
5. **Spaced Repetition**: FSRS flashcard untuk retention
6. **Scalable Architecture**: Go backend + Docker + Cloud Run ready

---

**Good luck dengan presentasinya! 🚀**





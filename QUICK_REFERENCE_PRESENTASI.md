# ⚡ Quick Reference - Presentasi LIDM

## 🎯 Elevator Pitch (30 detik)
"LIDM adalah platform pembelajaran interaktif untuk SD kelas 4 yang fokus pada fotosintesis. Aplikasi menggunakan gamifikasi dengan progressive unlock system, AI tutor bernama Alsan, dan real-time multiplayer quiz untuk meningkatkan engagement siswa."

---

## 🏗️ Arsitektur Singkat

```
┌─────────────┐         ┌─────────────┐
│   Mobile    │────────▶│  lidm-api   │ (Go/Gin + MySQL)
│   App       │         │  (Backend)  │
└─────────────┘         └──────┬──────┘
                               │
                               ▼
                        ┌─────────────┐
                        │  lidm-ai    │ (Flask + Groq AI)
                        │ (AI Chat)   │
                        └─────────────┘
```

---

## 📊 Core Features (1 menit)

### 1. **Progressive Module Unlock** 🔓
- Module 1 selalu unlocked
- Auto-unlock module berikutnya saat progress 100%
- Progress = (completed prequiz + video quiz) / total × 100%

### 2. **Quiz System** 🎮
- **Prequiz**: Pre-assessment sebelum belajar
- **Video Quiz**: Post-video assessment
- **Solo Quiz**: Latihan sendiri (5 lives)
- **Multiplayer Quiz**: Real-time head-to-head (Socket.IO)

### 3. **AI Tutor (Alsan)** 🤖
- Chat AI khusus fotosintesis
- Groq Llama 4 model
- Bahasa sederhana untuk SD kelas 4
- Multimodal (text + image)

### 4. **Gamification** 🏆
- Points system (10 per correct answer)
- Leaderboard dengan ranking
- Daily streak tracking
- Lives system untuk solo quiz

### 5. **Flashcard dengan FSRS** 📚
- Spaced repetition algorithm
- Review berdasarkan interval (1m, 5m, 7h, dll)
- Retention statistics

---

## 🔑 Key Technical Points

### Backend (lidm-api)
- **Language**: Go 1.24.4
- **Framework**: Gin
- **Database**: MySQL (GORM)
- **Real-time**: Socket.IO v2
- **Auth**: JWT

### AI Service (lidm-ai)
- **Language**: Python 3.11
- **Framework**: Flask
- **AI**: Groq API (Llama 4)
- **Database**: MongoDB (chat history)

---

## 📈 User Flow (Demo Flow)

```
1. Register → Email verification
2. Login → JWT token
3. Module 1 unlocked (default)
4. Belajar:
   ├─ Baca materi
   ├─ Jawab prequiz → Update progress
   ├─ Tonton video
   ├─ Jawab video quiz → Update progress
   └─ Review flashcards
5. Progress 100% → Module 2 auto-unlock
6. Repeat untuk module berikutnya
```

---

## 🎮 Quiz Flow

### Solo Quiz:
```
Start → Question → Answer → Result → Next Question
         ↓ (salah)
      Lives -1 → Continue
```

### Multiplayer Quiz:
```
User A: Create Lobby → Invite Code
User B: Join Lobby → Match Found
Server: Start Quiz → Real-time Questions
Both: Submit Answers → Real-time Scoring
Server: Quiz Completed → Update Leaderboard
```

---

## 💡 Value Propositions

1. **Progressive Learning**: Sequential unlock memastikan step-by-step learning
2. **Personalized AI**: Chat AI dengan bahasa sesuai level siswa
3. **Gamification**: Points, leaderboard, streaks untuk engagement
4. **Real-time Competition**: Multiplayer quiz untuk motivasi
5. **Spaced Repetition**: FSRS untuk retention jangka panjang

---

## ❓ FAQ - Jawaban Singkat

**Q: Bagaimana auto-unlock bekerja?**
A: Module 1 default unlocked. Saat user selesaikan semua prequiz + video quiz (100%), module berikutnya otomatis unlock.

**Q: Perbedaan solo vs multiplayer?**
A: Solo = latihan sendiri (5 lives). Multiplayer = head-to-head real-time (no lives, kompetisi langsung).

**Q: Bagaimana scoring?**
A: 10 points per correct answer. Total score = sum semua quiz selesai. Leaderboard ranking berdasarkan total score.

**Q: AI chat fokus apa?**
A: Khusus fotosintesis dengan bahasa sederhana SD kelas 4. Mengurangi miskonsepsi siswa.

**Q: Kenapa Go untuk backend?**
A: Performance tinggi, concurrent handling (goroutines), type safety, good untuk real-time systems.

**Q: Bagaimana progress tracking?**
A: Real-time calculation: (completed items / total items) × 100%. Auto-update setiap submit answer.

---

## 🎯 Demo Checklist

- [ ] Register & Login
- [ ] Module 1 unlocked (default)
- [ ] Jawab prequiz → Progress update
- [ ] Jawab video quiz → Progress update
- [ ] Progress 100% → Module 2 unlock
- [ ] Multiplayer quiz (2 devices)
- [ ] AI chat dengan Alsan
- [ ] Leaderboard
- [ ] Flashcard review

---

## 📊 Key Metrics to Mention

- **Progress Tracking**: Real-time 0-100%
- **Auto-unlock**: Sequential module unlocking
- **Real-time**: Socket.IO untuk multiplayer (< 100ms latency)
- **AI Response**: Groq API (fast inference)
- **Scalability**: Docker + Cloud Run ready

---

## 🚀 Tech Stack Summary

| Component | Technology |
|-----------|-----------|
| Backend API | Go 1.24 + Gin |
| AI Service | Python 3.11 + Flask |
| Database | MySQL (main) + MongoDB (chat) |
| Real-time | Socket.IO v2 |
| AI Model | Groq Llama 4 Scout 17B |
| Auth | JWT |
| Deployment | Docker + Cloud Run |

---

## 💬 Opening Statement

"Selamat pagi/siang, saya akan mempresentasikan LIDM - platform pembelajaran interaktif untuk SD kelas 4 yang fokus pada fotosintesis. Aplikasi ini menggabungkan gamifikasi, AI tutor, dan real-time multiplayer quiz untuk meningkatkan engagement dan pemahaman siswa."

---

## 🎬 Closing Statement

"LIDM tidak hanya platform pembelajaran, tapi juga ekosistem yang menggabungkan progressive learning, personalized AI assistance, dan gamification untuk menciptakan pengalaman belajar yang engaging dan efektif. Dengan teknologi Go untuk performa tinggi dan AI untuk personalisasi, kami yakin aplikasi ini dapat meningkatkan pemahaman siswa tentang fotosintesis."

---

**Tips Presentasi:**
- Mulai dengan problem statement (kesulitan belajar fotosintesis)
- Highlight unique features (auto-unlock, AI tutor, multiplayer)
- Demo flow yang smooth (prepare sebelumnya)
- Siapkan jawaban untuk pertanyaan teknis
- Emphasize user experience dan engagement

**Good luck! 🎉**





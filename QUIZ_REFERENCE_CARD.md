# Quick Reference: Seeded Quiz Questions & Answers

## Module: "Belajar Fotosintesis - Kelas 4 SD"

### SubMaterial 1: "Apa itu Fotosintesis?" - Prequizzes

**Prequiz 1:**
- **Question:** "Apa yang membuat daun berwarna hijau?"
- **Options:** A) Air, B) Klorofil, C) Tanah, D) Udara
- **Correct Answer:** B
- **Explanation:** "Klorofil adalah zat hijau yang membuat daun berwarna hijau"

**Prequiz 2:**
- **Question:** "Dimana fotosintesis terjadi?"
- **Options:** A) Akar, B) Daun, C) Batang, D) Bunga
- **Correct Answer:** B
- **Explanation:** "Fotosintesis terjadi di daun karena mengandung klorofil"

**Prequiz 3:**
- **Question:** "Apa yang dibutuhkan untuk fotosintesis?"
- **Options:** A) Sinar matahari, air, dan udara, B) Tanah dan pupuk, C) Air saja, D) Udara saja
- **Correct Answer:** A
- **Explanation:** "Fotosintesis membutuhkan sinar matahari, air, dan karbon dioksida dari udara"

### SubMaterial 2: "Mengapa Fotosintesis Penting?" - Prequizzes

**Prequiz 4:**
- **Question:** "Mengapa fotosintesis penting bagi kita?"
- **Options:** A) Menghasilkan oksigen untuk bernapas, B) Membuat tumbuhan cantik, C) Tidak penting, D) Membuat udara kotor
- **Correct Answer:** A
- **Explanation:** "Fotosintesis menghasilkan oksigen yang kita butuhkan untuk bernapas!"

**Prequiz 5:**
- **Question:** "Apa yang dihasilkan dari fotosintesis?"
- **Options:** A) Air dan tanah, B) Glukosa dan oksigen, C) Karbon dioksida, D) Sampah
- **Correct Answer:** B
- **Explanation:** "Fotosintesis menghasilkan glukosa (makanan tumbuhan) dan oksigen!"

**Prequiz 6:**
- **Question:** "Kapan tumbuhan melakukan fotosintesis?"
- **Options:** A) Saat ada sinar matahari, B) Malam hari, C) Saat hujan, D) Tidak pernah
- **Correct Answer:** A
- **Explanation:** "Tumbuhan melakukan fotosintesis saat ada sinar matahari karena matahari adalah sumber energinya!"

### SubMaterial 3: "Laboratorium Virtual AR" - Prequizzes

**Prequiz 7:**
- **Question:** "Apa yang ingin kamu lihat di dalam daun?"
- **Options:** A) Warna hijau saja, B) Proses fotosintesis, C) Tidak ada yang menarik, D) Serangga kecil
- **Correct Answer:** B
- **Explanation:** "Dengan AR, kita bisa melihat bagaimana fotosintesis terjadi di dalam daun!"

**Prequiz 8:**
- **Question:** "Menurutmu, di bagian mana fotosintesis terjadi?"
- **Options:** A) Di seluruh tumbuhan, B) Hanya di daun, C) Di akar saja, D) Di batang saja
- **Correct Answer:** B
- **Explanation:** "Fotosintesis terutama terjadi di daun, mari kita lihat di AR Lab!"

**Prequiz 9:**
- **Question:** "Apa yang kamu bayangkan ada di dalam sel daun?"
- **Options:** A) Ruang kosong, B) Kloroplas hijau, C) Air saja, D) Tidak tahu
- **Correct Answer:** B
- **Explanation:** "Kloroplas adalah tempat fotosintesis terjadi, kita akan melihatnya di AR!"

## Video Materials - Video Quizzes

### Video Material 1: "Video Pengenalan Fotosintesis"

**Video Quiz 1:**
- **Question:** "Berdasarkan video, apa yang dibutuhkan untuk fotosintesis?"
- **Timestamp:** 30-45 seconds
- **Options:** A) Air saja, B) Sinar matahari, air, dan CO2, C) Tanah dan pupuk, D) Oksigen dan nitrogen
- **Correct Answer:** B
- **Explanation:** "Fotosintesis membutuhkan sinar matahari, air, dan karbon dioksida"

**Video Quiz 2:**
- **Question:** "Apa hasil dari fotosintesis menurut video?"
- **Timestamp:** 120-135 seconds
- **Options:** A) Air dan tanah, B) Glukosa dan oksigen, C) Karbon dioksida, D) Protein
- **Correct Answer:** B
- **Explanation:** "Fotosintesis menghasilkan glukosa dan oksigen yang penting untuk kehidupan"

### Video Material 2: "Video Proses Fotosintesis"

**Video Quiz 3:**
- **Question:** "Berdasarkan video, di bagian mana fotosintesis terjadi?"
- **Timestamp:** 45-60 seconds
- **Options:** A) Di akar, B) Di batang, C) Di kloroplas dalam daun, D) Di bunga
- **Correct Answer:** C
- **Explanation:** "Benar! Fotosintesis terjadi di kloroplas yang terdapat dalam sel-sel daun!"

**Video Quiz 4:**
- **Question:** "Apa hasil utama dari proses fotosintesis yang dijelaskan dalam video?"
- **Timestamp:** 120-135 seconds
- **Options:** A) Air dan tanah, B) Glukosa dan oksigen, C) Karbon dioksida dan nitrogen, D) Protein dan lemak
- **Correct Answer:** B
- **Explanation:** "Tepat! Fotosintesis menghasilkan glukosa (makanan tumbuhan) dan oksigen yang kita hirup!"

**Video Quiz 5:**
- **Question:** "Menurut video, mengapa fotosintesis penting bagi kehidupan?"
- **Timestamp:** 200-215 seconds
- **Options:** A) Menghasilkan oksigen untuk bernapas, B) Membuat tumbuhan terlihat cantik, C) Menghasilkan suara, D) Membuat udara panas
- **Correct Answer:** A
- **Explanation:** "Benar! Fotosintesis menghasilkan oksigen yang sangat penting untuk semua makhluk hidup bernapas!"

## Quick Test Commands

### Test Prequizzes:
```bash
# Get SubMaterial 1 prequizzes
GET /prequiz/submaterial/1

# Answer with correct answers
POST /prequiz/submit {"prequiz_id": 1, "answer": "B"}
POST /prequiz/submit {"prequiz_id": 2, "answer": "B"}
POST /prequiz/submit {"prequiz_id": 3, "answer": "A"}
```

### Test Video Quizzes:
```bash
# Get Video Material 1 quizzes
GET /video-quiz/video/1

# Answer with correct answers
POST /video-quiz/submit {"video_quiz_id": 1, "selected_answer": "B"}
POST /video-quiz/submit {"video_quiz_id": 2, "selected_answer": "B"}
```

## Summary:
- **Total Prequizzes:** 9 (3 per SubMaterial)
- **Total Video Quizzes:** 5 (2 for Video 1, 3 for Video 2)
- **All questions are about photosynthesis for 4th grade elementary**
- **Designed to be educational and age-appropriate**

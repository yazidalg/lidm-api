# How to Answer Quizzes in LIDM API

## Prerequisites
- You must be logged in and have a valid JWT token
- Use the token in Authorization header: `Bearer <your-token>`

## 1. Answering Prequizzes in Sub Materials

### Step 1: Get Prequizzes for a Sub Material
```bash
GET http://localhost:3000/prequiz/submaterial/{sub_material_id}
Authorization: Bearer <your-token>
```

**Example:**
```bash
GET http://localhost:3000/prequiz/submaterial/1
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Response Example:**
```json
{
  "success": true,
  "message": "Prequizzes retrieved successfully",
  "data": [
    {
      "id": 1,
      "sub_material_id": 1,
      "question": "Apa yang membuat daun berwarna hijau?",
      "option_a": "Air",
      "option_b": "Klorofil",
      "option_c": "Tanah",
      "option_d": "Udara",
      "correct_answer": "B",
      "explanation": "Klorofil adalah zat hijau yang membuat daun berwarna hijau"
    },
    {
      "id": 2,
      "sub_material_id": 1,
      "question": "Dimana fotosintesis terjadi?",
      "option_a": "Akar",
      "option_b": "Daun",
      "option_c": "Batang",
      "option_d": "Bunga",
      "correct_answer": "B",
      "explanation": "Fotosintesis terjadi di daun karena mengandung klorofil"
    }
  ]
}
```

### Step 2: Submit Prequiz Answer
```bash
POST http://localhost:3000/prequiz/submit
Authorization: Bearer <your-token>
Content-Type: application/json

{
  "prequiz_id": 1,
  "answer": "B"
}
```

**Response Example:**
```json
{
  "success": true,
  "message": "Answer submitted successfully",
  "data": {
    "id": 123,
    "prequiz_id": 1,
    "user_id": 5,
    "answer": "B",
    "is_correct": true,
    "answered_at": 1692360000
  },
  "is_correct": true,
  "correct_answer": "B",
  "explanation": "Klorofil adalah zat hijau yang membuat daun berwarna hijau"
}
```

## 2. Answering Video Quizzes

### Step 1: Get Video Quizzes for a Video Material
```bash
GET http://localhost:3000/video-quiz/video/{video_material_id}
Authorization: Bearer <your-token>
```

**Example:**
```bash
GET http://localhost:3000/video-quiz/video/1
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Response Example:**
```json
{
  "success": true,
  "message": "Video quizzes retrieved successfully",
  "data": [
    {
      "id": 6,
      "video_material_id": 1,
      "question": "Berdasarkan video, apa yang dibutuhkan untuk fotosintesis?",
      "timestamp_start": 30,
      "timestamp_end": 45,
      "option_a": "Air saja",
      "option_b": "Sinar matahari, air, dan CO2",
      "option_c": "Tanah dan pupuk",
      "option_d": "Oksigen dan nitrogen",
      "correct_answer": "B",
      "explanation": "Fotosintesis membutuhkan sinar matahari, air, dan karbon dioksida",
      "order": 1
    },
    {
      "id": 7,
      "video_material_id": 1,
      "question": "Apa hasil dari fotosintesis?",
      "timestamp_start": 120,
      "timestamp_end": 135,
      "option_a": "Air dan tanah",
      "option_b": "Glukosa dan oksigen",
      "option_c": "Karbon dioksida",
      "option_d": "Protein",
      "correct_answer": "B",
      "explanation": "Fotosintesis menghasilkan glukosa dan oksigen",
      "order": 2
    }
  ]
}
```

### Step 2: Submit Video Quiz Answer
```bash
POST http://localhost:3000/video-quiz/submit
Authorization: Bearer <your-token>
Content-Type: application/json

{
  "video_quiz_id": 6,
  "selected_answer": "B"
}
```

**Response Example:**
```json
{
  "success": true,
  "message": "Video quiz answer submitted successfully",
  "data": {
    "id": 456,
    "video_quiz_id": 6,
    "user_id": 5,
    "selected_answer": "B",
    "is_correct": true,
    "answered_at": 1692360000,
    "response_time": 8
  },
  "is_correct": true,
  "correct_answer": "B",
  "explanation": "Fotosintesis membutuhkan sinar matahari, air, dan karbon dioksida"
}
```

## 3. Complete Flow Examples with Real Data

Based on our seeded data, here are real examples you can test:

### Example 1: Answer Prequizzes for "Apa itu Fotosintesis?" (SubMaterial ID: 1)

#### Get Prequizzes:
```bash
GET http://localhost:3000/prequiz/submaterial/1
Authorization: Bearer <your-token>
```

#### Answer First Prequiz:
```bash
POST http://localhost:3000/prequiz/submit
Authorization: Bearer <your-token>
Content-Type: application/json

{
  "prequiz_id": 1,
  "answer": "B"
}
```

#### Answer Second Prequiz:
```bash
POST http://localhost:3000/prequiz/submit
Authorization: Bearer <your-token>
Content-Type: application/json

{
  "prequiz_id": 2,
  "answer": "A"
}
```

#### Answer Third Prequiz:
```bash
POST http://localhost:3000/prequiz/submit
Authorization: Bearer <your-token>
Content-Type: application/json

{
  "prequiz_id": 3,
  "answer": "A"
}
```

### Example 2: Answer Video Quizzes for "Video Pengenalan" (Video Material ID: 1)

#### Get Video Quizzes:
```bash
GET http://localhost:3000/video-quiz/video/1
Authorization: Bearer <your-token>
```

#### Answer First Video Quiz:
```bash
POST http://localhost:3000/video-quiz/submit
Authorization: Bearer <your-token>
Content-Type: application/json

{
  "video_quiz_id": 1,
  "selected_answer": "B"
}
```

#### Answer Second Video Quiz:
```bash
POST http://localhost:3000/video-quiz/submit
Authorization: Bearer <your-token>
Content-Type: application/json

{
  "video_quiz_id": 2,
  "selected_answer": "C"
}
```

## 4. Postman Collection Steps

### For Prequizzes:
1. **Login** → Copy token
2. **GET Prequizzes** → `GET /prequiz/submaterial/1` with Bearer token
3. **Submit Answer** → `POST /prequiz/submit` with Bearer token and answer

### For Video Quizzes:
1. **Login** → Copy token
2. **GET Video Quizzes** → `GET /video-quiz/video/1` with Bearer token
3. **Submit Answer** → `POST /video-quiz/submit` with Bearer token and answer

## 5. Available Sub Materials and Videos from Seeded Data

### Sub Materials with Prequizzes:
- **SubMaterial ID 1**: "Apa itu Fotosintesis?" (3 prequizzes)
- **SubMaterial ID 2**: "Mengapa Fotosintesis Penting?" (3 prequizzes)
- **SubMaterial ID 3**: "Melihat Fotosintesis dengan AR" (3 prequizzes)

### Video Materials with Quizzes:
- **Video Material ID 1**: "Video Pengenalan Fotosintesis" (2 video quizzes)
- **Video Material ID 2**: "Video Proses Fotosintesis" (3 video quizzes)

## 6. Error Handling

### Common Errors:

#### 401 Unauthorized
```json
{
  "message": "Authorization token required"
}
```
**Solution:** Add Authorization header with Bearer token

#### 404 Not Found
```json
{
  "message": "Prequiz not found"
}
```
**Solution:** Check if prequiz_id or video_quiz_id exists

#### 400 Bad Request
```json
{
  "message": "Answer already submitted for this prequiz"
}
```
**Solution:** User has already answered this quiz

## 7. Tips for Testing

1. **Use different users** to test multiple answers
2. **Try wrong answers** to see `is_correct: false`
3. **Check response times** for video quizzes
4. **Test all sub materials** (1, 2, 3) for comprehensive testing
5. **Test both video materials** (1, 2) for video quiz functionality

## 8. Dummy User Credentials for Testing

Use these accounts for testing:
- `andi.pratama@student.com` / `password123`
- `sari.dewi@student.com` / `password123`
- `budi.santoso@student.com` / `password123`
- `maya.putri@student.com` / `password123`
- `riko.firmansyah@student.com` / `password123`

Each user can answer all quizzes independently!

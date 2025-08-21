# CORRECTED API Endpoints Guide

## ✅ **WORKING ENDPOINTS** - Updated Routes

### 🔐 **Authentication (Required for all quiz endpoints)**
```bash
POST /auth/login
{
  "email": "andi.pratama@student.com",
  "password": "password123"
}
```

### 📝 **Prequizzes (SubMaterial Quizzes)**

#### ✅ Get Prequizzes by SubMaterial
```bash
GET /prequiz/submaterial/{sub_material_id}
Authorization: Bearer <token>
```
**Example:**
```bash
GET http://localhost:3000/prequiz/submaterial/1
Authorization: Bearer <your-token>
```

#### ✅ Submit Prequiz Answer
```bash
POST /prequiz/submit
Authorization: Bearer <token>
Content-Type: application/json

{
  "prequiz_id": 1,
  "answer": "B"
}
```

### 🎥 **Video Quizzes**

#### ✅ Get Video Quizzes by Video Material
```bash
GET /video-quiz/video-material/{video_material_id}
Authorization: Bearer <token>
```
**Example:**
```bash
GET http://localhost:3000/video-quiz/video-material/1
Authorization: Bearer <your-token>
```

#### ✅ Submit Video Quiz Answer
```bash
POST /video-quiz/submit
Authorization: Bearer <token>
Content-Type: application/json

{
  "video_quiz_id": 1,
  "selected_answer": "B"
}
```

## 🧪 **Test Data Available**

### SubMaterials with Prequizzes:
- **SubMaterial 1**: "Apa itu Fotosintesis?" - 3 prequizzes
- **SubMaterial 2**: "Mengapa Fotosintesis Penting?" - 3 prequizzes  
- **SubMaterial 3**: "Laboratorium Virtual AR" - 3 prequizzes

### Video Materials with Quizzes:
- **Video Material 1**: "Video Pengenalan Fotosintesis" - 2 video quizzes
- **Video Material 2**: "Video Proses Fotosintesis" - 3 video quizzes

## 📋 **Complete Testing Flow**

### Step 1: Login
```bash
curl -X POST http://localhost:3000/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "andi.pratama@student.com",
    "password": "password123"
  }'
```

### Step 2: Test Prequizzes
```bash
# Get prequizzes for SubMaterial 1
curl -X GET http://localhost:3000/prequiz/submaterial/1 \
  -H "Authorization: Bearer <your-token>"

# Submit answer (example: question 1, answer B)
curl -X POST http://localhost:3000/prequiz/submit \
  -H "Authorization: Bearer <your-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "prequiz_id": 1,
    "answer": "B"
  }'
```

### Step 3: Test Video Quizzes
```bash
# Get video quizzes for Video Material 1
curl -X GET http://localhost:3000/video-quiz/video-material/1 \
  -H "Authorization: Bearer <your-token>"

# Submit answer (example: video quiz 1, answer B)
curl -X POST http://localhost:3000/video-quiz/submit \
  -H "Authorization: Bearer <your-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "video_quiz_id": 1,
    "selected_answer": "B"
  }'
```

## 🆔 **Dummy User Accounts**
- `andi.pratama@student.com` / `password123`
- `sari.dewi@student.com` / `password123`
- `budi.santoso@student.com` / `password123`
- `maya.putri@student.com` / `password123`
- `riko.firmansyah@student.com` / `password123`

## ❌ **Common Issues & Solutions**

### Issue: 404 Not Found
**Solution:** Make sure you're using the correct endpoints:
- ✅ `/prequiz/submaterial/1` (not `/prequiz/submaterial`)
- ✅ `/video-quiz/video-material/1` (not `/video-quiz/video/1`)

### Issue: 401 Unauthorized
**Solution:** 
1. Login first to get token
2. Use `Authorization: Bearer <token>` header
3. Make sure token hasn't expired (7 days)

### Issue: Request Body Wrong Format
**Solution:**
- Prequizzes: `{"prequiz_id": 1, "answer": "B"}`
- Video Quizzes: `{"video_quiz_id": 1, "selected_answer": "B"}`

## 🎯 **Expected Responses**

### Successful Prequiz Answer:
```json
{
  "success": true,
  "message": "Prequiz answer submitted successfully",
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

### Successful Video Quiz Answer:
```json
{
  "success": true,
  "message": "Video quiz answer submitted successfully",
  "data": {
    "id": 456,
    "video_quiz_id": 1,
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

## 🚀 **Next Steps**

1. **Start the API server** if not running
2. **Test login** with dummy user
3. **Test prequizzes** with SubMaterial 1, 2, or 3
4. **Test video quizzes** with Video Material 1 or 2
5. **Try different answers** to see correct/incorrect responses

---
**📝 Note:** Both prequizzes and video quizzes are now working with the correct routes and the database has exactly 3 prequizzes per SubMaterial as requested!

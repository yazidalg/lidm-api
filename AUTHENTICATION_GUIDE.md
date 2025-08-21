# Authentication Guide untuk LIDM API

## Overview
API LIDM memerlukan authentication untuk mengakses sebagian besar endpoint, termasuk `/module/all`. Berikut cara melakukan authentication dan menggunakan token.

## 1. Login untuk Mendapatkan Token

### Option 1: Login dengan Email/Password
```bash
POST http://localhost:3000/auth/login
Content-Type: application/json

{
  "email": "andi.pratama@student.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Login successful",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "name": "Andi Pratama",
    "email": "andi.pratama@student.com",
    "role_id": 1
  }
}
```

### Option 2: Register User Baru (Jika Belum Ada)
```bash
POST http://localhost:3000/auth/register
Content-Type: application/json

{
  "name": "Test User",
  "email": "test@example.com",
  "password": "password123",
  "class": "4A",
  "role": "user"
}
```

## 2. Menggunakan Token untuk Akses API

Setelah mendapat token dari login, gunakan token tersebut di header `Authorization`:

### Mengakses Modules
```bash
GET http://localhost:3000/module/all
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### Mengakses Module Specific
```bash
GET http://localhost:3000/module/2
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

## 3. User Dummy yang Sudah Tersedia

Dari seeding yang baru saja dijalankan, tersedia 5 user dummy:

1. **Email:** `andi.pratama@student.com` | **Password:** `password123`
2. **Email:** `sari.dewi@student.com` | **Password:** `password123`
3. **Email:** `budi.santoso@student.com` | **Password:** `password123`
4. **Email:** `maya.putri@student.com` | **Password:** `password123`
5. **Email:** `riko.firmansyah@student.com` | **Password:** `password123`

## 4. Contoh Lengkap Postman

### Step 1: Login
1. **Method:** POST
2. **URL:** `http://localhost:3000/auth/login`
3. **Headers:** 
   - `Content-Type: application/json`
4. **Body (raw JSON):**
```json
{
  "email": "andi.pratama@student.com",
  "password": "password123"
}
```

### Step 2: Copy Token
Dari response login, copy nilai `token`

### Step 3: Access Modules
1. **Method:** GET
2. **URL:** `http://localhost:3000/module/all`
3. **Headers:**
   - `Authorization: Bearer <paste-token-here>`

## 5. Quiz Endpoints yang Tersedia

Setelah login, Anda bisa mengakses:

### Prequizzes
```bash
GET http://localhost:3000/prequiz/submaterial/1
Authorization: Bearer <token>
```

### Submit Prequiz Answer
```bash
POST http://localhost:3000/prequiz/submit
Authorization: Bearer <token>
Content-Type: application/json

{
  "prequiz_id": 1,
  "answer": "A"
}
```

### Video Quizzes
```bash
GET http://localhost:3000/video-quiz/video/1
Authorization: Bearer <token>
```

### Submit Video Quiz Answer
```bash
POST http://localhost:3000/video-quiz/submit
Authorization: Bearer <token>
Content-Type: application/json

{
  "video_quiz_id": 1,
  "selected_answer": "B"
}
```

## 6. Troubleshooting

### Error 401 Unauthorized
- Pastikan header `Authorization` sudah benar
- Pastikan token belum expired (valid 7 hari)
- Pastikan format: `Bearer <token>`

### Error "User not found"
- Token mungkin invalid atau user sudah dihapus
- Coba login ulang untuk mendapat token baru

### Error "Authorization token required"
- Header `Authorization` tidak ada
- Tambahkan header dengan format yang benar

## 7. Admin Access

Untuk akses admin (create/update/delete), gunakan akun admin:
```bash
POST http://localhost:3000/auth/register
Content-Type: application/json

{
  "name": "Admin User",
  "email": "admin@lidm.com",
  "password": "admin123",
  "class": "",
  "role": "admin"
}
```

Setelah register admin, login dan gunakan token untuk akses admin endpoints.

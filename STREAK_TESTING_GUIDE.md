# Panduan Curl untuk Streak System

## 1. Login untuk Mendapatkan Token

```bash
# Login dengan akun belajar
curl -X POST http://localhost:3000/auth/belajar-login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "azis@belajar.id",
    "password": "password123"
  }'
```

Response akan berisi token yang digunakan untuk endpoint lain.

## 2. Trigger Aktivitas Belajar (Otomatis Update Streak)

Streak akan otomatis update ketika user melakukan aktivitas belajar berikut:

### a. Melihat Pelajaran
```bash
curl -X GET "http://localhost:3000/lesson/all" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json"
```

### b. Melihat Modul
```bash
curl -X GET "http://localhost:3000/module/all" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json"
```

### c. Mengakses Quiz (jika ada endpoint)
```bash
curl -X GET "http://localhost:3000/quiz/list" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json"
```

## 3. Melihat Informasi Streak

### a. Cek Streak Pribadi
```bash
curl -X GET "http://localhost:3000/user-activity/my-streak" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json"
```

**Contoh Response:**
```json
{
  "message": "Informasi streak berhasil diambil",
  "data": {
    "current_streak": 3,
    "max_streak": 5,
    "streak_status": "Streak 3 hari, tetap semangat!"
  }
}
```

### b. Melihat Most Active User dengan Streak
```bash
curl -X GET "http://localhost:3000/user-activity/most-active-detailed" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json"
```

**Contoh Response dengan Streak:**
```json
{
  "message": "Pengguna paling aktif dengan detail berhasil diambil",
  "data": {
    "user_id": 11,
    "username": "azis@belajar.id",
    "total_activities": 15,
    "total_learning_minutes": 75,
    "total_learning_hours": 1.25,
    "current_streak": 3,
    "max_streak": 5,
    "activity_breakdown": {
      "lihat_pelajaran": 8,
      "lihat_modul": 4,
      "masuk": 3
    },
    "last_activity": {
      "ID": 25,
      "activity_type": "lihat_pelajaran",
      "description": "Melihat pelajaran",
      "created_at": "2025-08-04T10:30:00Z"
    },
    "time_since_last_activity": "2 menit yang lalu"
  }
}
```

## 4. Script Test Streak

Buat script untuk test streak secara otomatis:

```bash
#!/bin/bash

# Login
TOKEN=$(curl -s -X POST http://localhost:3000/auth/belajar-login \
  -H "Content-Type: application/json" \
  -d '{"email": "azis@belajar.id", "password": "password123"}' | \
  grep -o '"token":"[^"]*"' | cut -d'"' -f4)

echo "Token: $TOKEN"

# Trigger beberapa aktivitas belajar
echo "1. Melihat semua pelajaran..."
curl -s -X GET "http://localhost:3000/lesson/all" \
  -H "Authorization: Bearer $TOKEN" > /dev/null

echo "2. Melihat semua modul..."
curl -s -X GET "http://localhost:3000/module/all" \
  -H "Authorization: Bearer $TOKEN" > /dev/null

echo "3. Cek streak setelah aktivitas..."
curl -X GET "http://localhost:3000/user-activity/my-streak" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" | jq '.'
```

## 5. Cara Kerja Streak System

### Logika Streak:
- **Hari Pertama**: Streak dimulai = 1
- **Hari Berturut-turut**: Streak bertambah +1
- **Lewat 1 hari tanpa aktivitas**: Streak reset ke 0
- **Aktivitas di hari yang sama**: Streak tidak berubah

### Aktivitas yang Menambah Streak:
- ✅ `lihat_pelajaran` (lesson view)
- ✅ `selesai_pelajaran` (lesson complete)
- ✅ `lihat_modul` (module view)
- ✅ `selesai_modul` (module complete)
- ✅ `gabung_kuis` (quiz join)
- ✅ `selesai_kuis` (quiz complete)
- ✅ `jawab_kuis` (quiz answer)
- ❌ `masuk` (login) - tidak menambah streak
- ❌ `keluar` (logout) - tidak menambah streak

### Status Streak:
- 0 hari: "Belum ada streak, ayo mulai belajar!"
- 1 hari: "Streak baru dimulai, pertahankan!"
- 2-6 hari: "Streak X hari, tetap semangat!"
- 7-29 hari: "Streak X hari, luar biasa!"
- 30-99 hari: "Streak X hari, amazing!"
- 100+ hari: "Streak X hari, legendary!"

## 6. Testing Streak System

### Test Case 1: Streak Baru
```bash
# User baru atau belum pernah aktif
# Melakukan 1 aktivitas belajar → streak = 1
```

### Test Case 2: Streak Berturut-turut
```bash
# Hari 1: aktivitas → streak = 1
# Hari 2: aktivitas → streak = 2
# Hari 3: aktivitas → streak = 3
```

### Test Case 3: Putus Streak
```bash
# Hari 1: aktivitas → streak = 1
# Hari 2: tidak ada aktivitas
# Hari 4: aktivitas → streak = 1 (reset)
```

## 7. Database Migration

Sebelum testing, jalankan migration untuk menambahkan field streak:

```sql
-- Jalankan script ini di database
ALTER TABLE users ADD COLUMN current_streak INTEGER DEFAULT 0;
ALTER TABLE users ADD COLUMN max_streak INTEGER DEFAULT 0;
ALTER TABLE users ADD COLUMN last_active_date TIMESTAMP;
```

Atau gunakan file migration yang sudah dibuat: `add_user_streak_migration.sql`

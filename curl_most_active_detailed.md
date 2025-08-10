# Curl untuk Endpoint Pengguna Paling Aktif dengan Detail

## 1. Login Dulu untuk Mendapatkan Token

```bash
# Login dengan akun belajar
curl -X POST http://localhost:3000/auth/belajar-login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "azis@belajar.id",
    "password": "password123"
  }'
```

## 2. Test Endpoint Pengguna Paling Aktif dengan Detail

```bash
# Ganti TOKEN dengan token yang didapat dari login
curl -X GET "http://localhost:3000/user-activity/most-active-detailed" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" | jq '.'
```

## Contoh Response yang Diharapkan:

```json
{
  "message": "Pengguna paling aktif dengan detail berhasil diambil",
  "data": {
    "user_id": 1,
    "username": "azis@belajar.id",
    "total_activities": 25,
    "total_learning_minutes": 120,
    "total_learning_hours": 2.0,
    "activity_breakdown": {
      "lihat_pelajaran": 8,
      "selesai_pelajaran": 5,
      "lihat_modul": 3,
      "gabung_kuis": 4,
      "selesai_kuis": 3,
      "masuk": 2
    },
    "last_activity": {
      "id": 45,
      "user_id": 1,
      "activity_type": "selesai_pelajaran",
      "description": "Menyelesaikan pelajaran",
      "meta_data": "{\"lesson_id\":\"5\",\"path\":\"/lesson/5/complete\",\"method\":\"POST\"}",
      "ip_address": "127.0.0.1",
      "user_agent": "curl/7.84.0",
      "created_at": "2025-08-04T10:30:00Z"
    },
    "time_since_last_activity": "5 menit yang lalu"
  }
}
```

## Penjelasan Data:

- **user_id**: ID pengguna paling aktif
- **username**: Nama pengguna/email
- **total_activities**: Total jumlah aktivitas
- **total_learning_minutes**: Estimasi total waktu belajar dalam menit
- **total_learning_hours**: Estimasi total waktu belajar dalam jam
- **activity_breakdown**: Rincian jenis-jenis aktivitas
- **last_activity**: Detail aktivitas terakhir
- **time_since_last_activity**: Waktu sejak aktivitas terakhir

## Estimasi Waktu Belajar:

- **Lihat/Selesai Pelajaran**: 5 menit per aktivitas
- **Lihat/Selesai Modul**: 10 menit per aktivitas  
- **Gabung/Selesai Kuis**: 3 menit per aktivitas
- **Jawab Kuis**: 1 menit per aktivitas

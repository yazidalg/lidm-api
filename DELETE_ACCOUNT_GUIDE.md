# Delete Account Guide

## Endpoint untuk Hapus Akun User

### DELETE /user/delete-account

Endpoint ini digunakan untuk menghapus akun user **tanpa memerlukan authentication token**. User hanya perlu memberikan email dan password untuk verifikasi.

**Authentication:** Not Required (Public endpoint)

**Request Type:** `application/json`

---

## ⚠️ Cascade Delete - Data yang Dihapus

Ketika akun dihapus, **semua data terkait user juga akan dihapus secara otomatis**:

1. ✅ **Leaderboard Entry** - Statistik dan posisi user di leaderboard
2. ✅ **Quiz Sessions** - Semua session quiz yang pernah diikuti
3. ✅ **Participant Records** - History partisipasi dalam quiz
4. ✅ **Module Progress** - Progress belajar di semua modul
5. ✅ **User Activities** - Log aktivitas user
6. ✅ **Flashcard Progress** - Progress review flashcard
7. ✅ **User Account** - Data user utama

**⚠️ PENTING:** Penghapusan bersifat **PERMANEN** dan **TIDAK DAPAT DIKEMBALIKAN!**

---

## Request Format

```json
{
  "email": "user@example.com",
  "password": "user_password"
}
```

### Field Requirements
- `email` (string, required) - Email akun yang akan dihapus
- `password` (string, required) - Password akun untuk verifikasi

---

## Contoh Penggunaan

### Delete Account
```bash
curl -X DELETE http://localhost:8080/user/delete-account \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "mypassword123"
  }'
```

**Response Success (200 OK):**
```json
{
  "message": "Account deleted successfully"
}
```

---

## Response Codes

- **200 OK** - Account berhasil dihapus beserta semua data terkait
- **400 Bad Request** - Request tidak valid (email/password kosong)
- **401 Unauthorized** - Password salah
- **404 Not Found** - Email tidak ditemukan
- **500 Internal Server Error** - Gagal menghapus data

---

## Error Responses

### Email tidak ditemukan
```json
{
  "message": "User not found"
}
```

### Password salah
```json
{
  "message": "Invalid email or password"
}
```

### Request tidak valid
```json
{
  "message": "Invalid request",
  "error": "Key: 'DeleteAccountRequest.Email' Error:Field validation for 'Email' failed on the 'required' tag"
}
```

### Gagal menghapus
```json
{
  "message": "Failed to delete account",
  "error": "database error message"
}
```

---

## Keamanan & Validasi

### 1. Verifikasi Password
- Endpoint ini **tidak memerlukan JWT token**
- User harus memberikan **email dan password yang benar**
- Password di-verify menggunakan bcrypt

### 2. Transaction Safety
- Semua penghapusan dilakukan dalam **database transaction**
- Jika ada error di tengah proses, **semua perubahan di-rollback**
- Memastikan data consistency

### 3. Public Access
- Endpoint ini bisa diakses **tanpa login**
- Cocok untuk user yang ingin hapus akun dari halaman public

---

## Testing

### Jalankan test script lengkap:
```bash
./test_delete_account_cleanup.sh
```

Script ini akan:
1. ✅ Create test account
2. ✅ Login dan get token
3. ✅ Check leaderboard sebelum delete
4. ✅ Delete account
5. ✅ Verify account terhapus
6. ✅ Verify data leaderboard terhapus

### Manual test dengan existing account:
```bash
# Delete account
curl -X DELETE http://localhost:8080/user/delete-account \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "testpassword"
  }'

# Verify tidak bisa login lagi
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "testpassword"
  }'
```

---

## Database Transaction Flow

```
1. BEGIN TRANSACTION
2. DELETE leaderboards WHERE user_id = ?
3. DELETE quiz_sessions WHERE user_id = ?
4. DELETE participants WHERE user_id = ?
5. DELETE module_progresses WHERE user_id = ?
6. DELETE user_activities WHERE user_id = ?
7. DELETE user_flashcard_progresses WHERE user_id = ?
8. DELETE users WHERE id = ?
9. COMMIT TRANSACTION
```

Jika salah satu step gagal, semua perubahan akan di-rollback.

---

## Use Cases

### 1. User Self-Service Delete
User bisa hapus akun sendiri tanpa perlu contact admin:
```bash
# User mengisi form delete account
curl -X DELETE http://localhost:8080/user/delete-account \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "userpassword"
  }'
```

### 2. GDPR Compliance
Memenuhi requirement GDPR untuk "right to be forgotten":
- User data dihapus permanen
- Tidak ada trace di leaderboard
- Semua aktivitas terhapus

### 3. Account Reset
User yang ingin mulai dari awal bisa:
1. Delete account lama
2. Register dengan email yang sama
3. Mulai dengan data fresh

---

## Notes

1. **Email Reuse**: Setelah account dihapus, email tersebut bisa digunakan untuk register account baru.

2. **No Soft Delete**: Penghapusan adalah hard delete, data benar-benar dihapus dari database.

3. **Atomic Operation**: Semua penghapusan terjadi dalam satu transaction, ensuring data consistency.

4. **Profile Photo**: Jika user punya profile photo yang di-upload, file photo tidak otomatis terhapus dari disk. Pertimbangkan untuk menambahkan cleanup untuk file uploads.

---

## Related Endpoints

- `POST /auth/register` - Create new account
- `POST /auth/login` - Login to account
- `PUT /user/edit-profile` - Edit profile (requires auth)
- `GET /user/profile` - Get user profile (requires auth)

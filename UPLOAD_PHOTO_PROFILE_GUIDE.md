# Upload Photo Profile Guide

## Endpoint untuk Upload Foto Profile

### POST /user/upload-photo-profile

Upload foto profile dengan file multipart/form-data.

**Authentication:** Required (JWT Token)

**Request Type:** `multipart/form-data`

**Form Fields:**
- `photo_profile` (file, required) - File gambar profile
  - Max size: 5MB
  - Allowed types: .jpg, .jpeg, .png, .webp

---

## Contoh Penggunaan dengan cURL

### 1. Upload Foto Profile

```bash
curl -X POST http://localhost:8080/user/upload-photo-profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -F "photo_profile=@/path/to/your/photo.jpg"
```

**Response Success (200 OK):**
```json
{
  "message": "Photo profile uploaded successfully",
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "photo_profile": "uploads/profiles/1732534567_photo.jpg"
  }
}
```

---

## Contoh dengan File Berbeda

### Upload foto PNG
```bash
curl -X POST http://localhost:8080/user/upload-photo-profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -F "photo_profile=@/Users/john/Pictures/avatar.png"
```

### Upload foto WebP
```bash
curl -X POST http://localhost:8080/user/upload-photo-profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -F "photo_profile=@/Users/john/Pictures/profile.webp"
```

---

## Response Codes

- **200 OK** - Photo profile berhasil diupload
- **400 Bad Request** - File tidak valid atau melebihi ukuran maksimal
- **401 Unauthorized** - Token tidak valid atau tidak ada
- **500 Internal Server Error** - Gagal menyimpan ke database

---

## Error Responses

### File terlalu besar
```json
{
  "message": "Failed to upload photo profile",
  "error": "file size exceeds limit of 5242880 bytes"
}
```

### File type tidak diperbolehkan
```json
{
  "message": "Failed to upload photo profile",
  "error": "file type .gif not allowed. Allowed types: [.jpg .jpeg .png .webp]"
}
```

### Tidak ada file
```json
{
  "message": "Failed to upload photo profile",
  "error": "failed to get file: http: no such file"
}
```

### Token tidak valid
```json
{
  "message": "User not authenticated"
}
```

---

## Catatan Penting

1. **Automatic Cleanup**: Jika user sudah memiliki foto profile sebelumnya, file lama akan otomatis dihapus ketika upload foto baru.

2. **File Storage**: File akan disimpan di folder `uploads/profiles/` dengan nama file yang unique (menggunakan timestamp).

3. **Database Update**: Path foto profile akan otomatis disimpan di database di field `profile_picture`.

4. **Rollback on Failure**: Jika gagal update database, file yang sudah diupload akan otomatis dihapus.

5. **Security**: 
   - Hanya user yang sudah login yang bisa upload foto profile
   - Hanya bisa upload foto profile untuk akun sendiri
   - File type dan size sudah dibatasi untuk keamanan

---

## Testing dengan Postman

1. Buat request baru dengan method `POST`
2. URL: `http://localhost:8080/user/upload-photo-profile`
3. Pada tab "Authorization":
   - Type: Bearer Token
   - Token: Paste JWT token Anda
4. Pada tab "Body":
   - Pilih `form-data`
   - Tambah key: `photo_profile`
   - Ubah type dari "Text" ke "File"
   - Pilih file gambar Anda
5. Klik "Send"

---

## Akses File yang Diupload

File yang diupload bisa diakses melalui:
```
http://localhost:8080/uploads/profiles/filename.jpg
```

Pastikan aplikasi Anda sudah mengkonfigurasi static file serving untuk folder `uploads/`.

---

## Update Profile dengan Data Lain (Tanpa Upload File)

Jika ingin update nama/email tanpa upload foto baru, gunakan endpoint yang sudah ada:

```bash
curl -X PUT http://localhost:8080/user/edit-profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe Updated",
    "email": "john.updated@example.com"
  }'
```

Foto profile tidak akan berubah jika tidak menyertakan field `photo_profile`.

---

## Update Photo Profile dengan URL (Tanpa Upload File)

Anda juga bisa update photo profile langsung dengan URL tanpa upload file:

```bash
curl -X PUT http://localhost:8080/user/edit-profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "photo_profile": "https://example.com/photo.jpg"
  }'
```

**Response:**
```json
{
  "message": "Profile updated successfully",
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "photo_profile": "https://example.com/photo.jpg"
  }
}
```

---

## Dua Cara Update Photo Profile

### 1. Upload File Langsung (POST /user/upload-photo-profile)
- Upload file gambar dari komputer
- File disimpan di server
- Mendapat path file otomatis
- **Gunakan ini jika**: Ingin user upload foto dari device mereka

### 2. Kirim URL Photo (PUT /user/edit-profile)
- Kirim URL foto yang sudah ada di internet
- Tidak upload file ke server
- Langsung simpan URL-nya
- **Gunakan ini jika**: Foto sudah ada di cloud storage/CDN

---

## Testing Script

Jalankan test script untuk mencoba edit profile dengan photo:

```bash
# Set token dulu
export TOKEN="your_jwt_token_here"

# Jalankan test
./test_edit_profile_photo.sh
```

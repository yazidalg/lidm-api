# Cara Mendapatkan Google Client ID - Step by Step

## 🚀 Langkah 1: Buka Google Cloud Console

1. Buka browser dan pergi ke: **https://console.cloud.google.com/**
2. Login dengan akun Google Anda

## 🏗️ Langkah 2: Buat Project Baru

1. **Klik "Select a project"** di bagian atas halaman
2. **Klik "NEW PROJECT"** 
3. **Isi form project:**
   - Project name: `LIDM-Auth` (atau nama yang Anda inginkan)
   - Organization: (biarkan default)
4. **Klik "CREATE"**
5. **Tunggu** sampai project selesai dibuat (1-2 menit)
6. **Pastikan project Anda yang baru sudah selected** di bagian atas

## 🔧 Langkah 3: Enable Google Identity API

1. Di sidebar kiri, **klik "APIs & Services"**
2. **Klik "Library"**
3. Di search box, **ketik "Google Identity"**
4. **Klik "Google Identity Services"**
5. **Klik "ENABLE"**
6. Tunggu sampai enabled (beberapa detik)

## 🛡️ Langkah 4: Setup OAuth Consent Screen

1. Di sidebar kiri, **klik "APIs & Services" > "OAuth consent screen"**
2. **Pilih "External"** (karena ini untuk testing)
3. **Klik "CREATE"**
4. **Isi OAuth consent screen:**

### App Information:
- **App name:** `LIDM Learning App`
- **User support email:** [pilih email Anda dari dropdown]
- **App logo:** [skip dulu]

### App domain (optional - bisa skip):
- Application home page: `http://localhost:3000`
- Application privacy policy link: [skip]
- Application terms of service link: [skip]

### Developer contact information:
- **Email addresses:** [masukkan email Anda]

5. **Klik "SAVE AND CONTINUE"**

### Scopes (Langkah 2):
6. **Klik "SAVE AND CONTINUE"** (skip scopes untuk sekarang)

### Test users (Langkah 3):
7. **Klik "ADD USERS"**
8. **Masukkan email Anda** (untuk testing)
9. **Klik "ADD"**
10. **Klik "SAVE AND CONTINUE"**

### Summary (Langkah 4):
11. **Review informasi** dan **klik "BACK TO DASHBOARD"**

## 🔑 Langkah 5: Buat OAuth 2.0 Client ID

1. Di sidebar kiri, **klik "APIs & Services" > "Credentials"**
2. **Klik "CREATE CREDENTIALS"**
3. **Pilih "OAuth client ID"**

### Configure OAuth Client:
4. **Application type:** Pilih sesuai kebutuhan:
   - **Web application** (untuk React/Vue/Angular web app)
   - **Android** (untuk Android app)
   - **iOS** (untuk iOS app)

### Untuk Web Application:
5. **Name:** `LIDM Web Client`
6. **Authorized JavaScript origins:**
   ```
   http://localhost:3000
   http://localhost:3001
   http://localhost:5173
   http://127.0.0.1:3000
   ```
7. **Authorized redirect URIs:**
   ```
   http://localhost:3000
   http://localhost:3000/auth/callback
   http://localhost:3001
   http://localhost:5173
   ```

### Untuk Android:
5. **Name:** `LIDM Android Client`
6. **Package name:** `com.lidm.app` (sesuaikan dengan package app Anda)

### Untuk iOS:
5. **Name:** `LIDM iOS Client`
6. **Bundle ID:** `com.lidm.app` (sesuaikan dengan bundle ID app Anda)

8. **Klik "CREATE"**

## 🎉 Langkah 6: Copy Client ID

1. **Popup akan muncul** dengan Client ID dan Client Secret
2. **COPY Client ID** yang berbentuk seperti:
   ```
   123456789-abcdefghijklmnop.apps.googleusercontent.com
   ```
3. **Simpan** di tempat yang aman
4. **Klik "OK"**

## 📝 Langkah 7: Update File .env

1. **Buka file .env** di project Anda
2. **Update GOOGLE_CLIENT_ID:**
   ```env
   GOOGLE_CLIENT_ID=123456789-abcdefghijklmnop.apps.googleusercontent.com
   ```
3. **Save file**

## 🧪 Langkah 8: Test Setup

### Option 1: Menggunakan Google OAuth Playground
1. **Buka:** https://developers.google.com/oauthplayground/
2. **Klik gear icon** di pojok kanan atas
3. **Centang "Use your own OAuth credentials"**
4. **Masukkan Client ID** yang sudah Anda copy
5. **Masukkan Client Secret** (jika diperlukan)
6. **Klik "Close"**
7. **Di Step 1:** pilih `https://www.googleapis.com/auth/userinfo.profile`
8. **Klik "Authorize APIs"**
9. **Login dengan Google**
10. **Klik "Exchange authorization code for tokens"**
11. **Copy id_token** yang muncul

### Option 2: Test dengan Script
```bash
# Jalankan script test yang sudah dibuat
./test_google_auth.sh [ID_TOKEN_DARI_PLAYGROUND]
```

## 🔍 Troubleshooting

### Error: "redirect_uri_mismatch"
- **Solusi:** Pastikan URL redirect di Google Console sama dengan yang digunakan di aplikasi

### Error: "unauthorized_client"
- **Solusi:** Pastikan Client ID benar dan OAuth consent screen sudah di-setup

### Error: "access_denied"
- **Solusi:** Pastikan user sudah ditambahkan sebagai test user

### Error: "invalid_client"
- **Solusi:** Pastikan Client ID format benar dan project sudah enabled Google Identity API

## 📱 Format Client ID yang Benar

Client ID Google selalu berbentuk:
```
[NUMBERS]-[RANDOM_STRING].apps.googleusercontent.com
```

Contoh:
```
123456789012-abc123def456ghi789jkl012mno345pqr.apps.googleusercontent.com
```

## 🛠️ Untuk Development vs Production

### Development:
- Gunakan `http://localhost:3000`
- User type: External
- Verification status: Testing

### Production:
- Gunakan domain HTTPS yang sebenarnya
- User type: External
- Verification status: In production (setelah review Google)

## 📚 Resources

- [Google Cloud Console](https://console.cloud.google.com/)
- [OAuth 2.0 Playground](https://developers.google.com/oauthplayground/)
- [Google Identity Documentation](https://developers.google.com/identity)
- [OAuth 2.0 Scopes](https://developers.google.com/identity/protocols/oauth2/scopes)

---

**💡 Tips:**
- Simpan Client ID dan Client Secret di password manager
- Gunakan environment variables yang berbeda untuk dev dan prod
- Jangan commit credentials ke Git
- Review OAuth consent screen secara berkala

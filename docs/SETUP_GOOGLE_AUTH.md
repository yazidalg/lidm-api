# Setup Google Client ID dan JWT Secret

## 1. Mendapatkan Google Client ID

### Langkah 1: Buat Project di Google Cloud Console
1. Buka [Google Cloud Console](https://console.cloud.google.com/)
2. Login dengan akun Google Anda
3. Klik "Select a project" atau "New Project"
4. Klik "NEW PROJECT"
5. Masukkan nama project (misal: "LIDM Auth")
6. Klik "CREATE"

### Langkah 2: Enable Google+ API atau Identity API
1. Di sidebar kiri, pilih "APIs & Services" > "Library"
2. Cari "Google+ API" atau "Google Identity"
3. Klik pada API tersebut
4. Klik "ENABLE"

### Langkah 3: Buat OAuth 2.0 Credentials
1. Di sidebar kiri, pilih "APIs & Services" > "Credentials"
2. Klik "CREATE CREDENTIALS" > "OAuth client ID"
3. Jika belum ada OAuth consent screen, Anda akan diminta membuatnya:
   - Pilih "External" (untuk testing)
   - Isi informasi aplikasi:
     - App name: "LIDM Learning App"
     - User support email: email Anda
     - Developer contact: email Anda
   - Klik "SAVE AND CONTINUE"
   - Skip scopes (klik "SAVE AND CONTINUE")
   - Skip test users (klik "SAVE AND CONTINUE")

4. Setelah OAuth consent screen selesai, kembali ke Credentials
5. Klik "CREATE CREDENTIALS" > "OAuth client ID"
6. Pilih "Application type":
   - **Web application** (untuk React/Vue/Angular)
   - **Android** (untuk Android app)
   - **iOS** (untuk iOS app)

### Untuk Web Application:
- Name: "LIDM Web Client"
- Authorized JavaScript origins:
  ```
  http://localhost:3000
  http://localhost:3001
  https://yourdomain.com
  ```
- Authorized redirect URIs:
  ```
  http://localhost:3000
  http://localhost:3001
  https://yourdomain.com
  ```

### Untuk Mobile (Android/iOS):
- Name: "LIDM Mobile Client"
- Package name (Android): com.lidm.app
- Bundle ID (iOS): com.lidm.app

7. Klik "CREATE"
8. **Copy Client ID** yang muncul

## 2. Mendapatkan JWT Secret

JWT Secret adalah string random yang Anda buat sendiri untuk signing JWT tokens.

### Cara Generate JWT Secret:

#### Opsi 1: Manual (Simple)
```bash
# Buat string random 32-64 karakter
# Contoh:
lidm_super_secret_jwt_key_2024_very_secure_random_string
```

#### Opsi 2: Using OpenSSL
```bash
openssl rand -base64 32
```

#### Opsi 3: Using Node.js
```bash
node -e "console.log(require('crypto').randomBytes(32).toString('hex'))"
```

#### Opsi 4: Online Generator
- Buka https://generate-secret.vercel.app/32
- Copy hasil yang di-generate

## 3. Update File .env

Setelah mendapatkan Google Client ID dan JWT Secret, update file `.env`:

```env
# Database
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=yourpassword
DB_NAME=lidm

# JWT Secret (32-64 karakter random string)
JWT_SECRET=your_super_secret_jwt_key_here_32_chars_min

# Google OAuth
GOOGLE_CLIENT_ID=123456789-abcdefghijklmnop.apps.googleusercontent.com

# Email (if using email verification)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASS=your-app-password

# Server
PORT=3000
```

## 4. Testing Google Auth

### Test dengan Postman/curl:
1. Dapatkan Google ID Token dari client (web/mobile)
2. Test endpoint:

```bash
curl -X POST http://localhost:3000/auth/google \
  -H "Content-Type: application/json" \
  -d '{"id_token": "GOOGLE_ID_TOKEN_HERE"}'
```

### Mendapatkan Test ID Token:

#### Untuk Development/Testing:
1. Buka [Google OAuth 2.0 Playground](https://developers.google.com/oauthplayground/)
2. Di sebelah kiri, pilih "Google OAuth2 API v2"
3. Select "https://www.googleapis.com/auth/userinfo.profile"
4. Klik "Authorize APIs"
5. Login dengan Google
6. Klik "Exchange authorization code for tokens"
7. Copy `id_token` yang muncul

## 5. Client Integration Examples

### JavaScript (Web)
```html
<!-- Include Google API -->
<script src="https://apis.google.com/js/api:client.js"></script>

<script>
function initGoogleAuth() {
    gapi.load('auth2', function() {
        gapi.auth2.init({
            client_id: 'YOUR_GOOGLE_CLIENT_ID'
        });
    });
}

function signInWithGoogle() {
    const authInstance = gapi.auth2.getAuthInstance();
    authInstance.signIn().then(function(googleUser) {
        const idToken = googleUser.getAuthResponse().id_token;
        
        // Send to your backend
        fetch('/auth/google', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                id_token: idToken
            })
        })
        .then(response => response.json())
        .then(data => {
            localStorage.setItem('token', data.token);
            console.log('Login successful:', data);
        });
    });
}

// Initialize when page loads
initGoogleAuth();
</script>

<!-- Login button -->
<button onclick="signInWithGoogle()">Sign in with Google</button>
```

### React
```jsx
npm install @google-cloud/local-auth googleapis

// components/GoogleLogin.jsx
import { useEffect } from 'react';

function GoogleLogin() {
    useEffect(() => {
        // Load Google API
        const script = document.createElement('script');
        script.src = 'https://apis.google.com/js/api:client.js';
        script.onload = initGoogleAuth;
        document.body.appendChild(script);
    }, []);

    const initGoogleAuth = () => {
        window.gapi.load('auth2', function() {
            window.gapi.auth2.init({
                client_id: process.env.REACT_APP_GOOGLE_CLIENT_ID
            });
        });
    };

    const handleGoogleLogin = async () => {
        const authInstance = window.gapi.auth2.getAuthInstance();
        const googleUser = await authInstance.signIn();
        const idToken = googleUser.getAuthResponse().id_token;
        
        try {
            const response = await fetch('/auth/google', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ id_token: idToken })
            });
            
            const data = await response.json();
            localStorage.setItem('token', data.token);
            console.log('Login successful:', data);
        } catch (error) {
            console.error('Login failed:', error);
        }
    };

    return (
        <button onClick={handleGoogleLogin}>
            Sign in with Google
        </button>
    );
}

export default GoogleLogin;
```

## 6. Security Notes

### JWT Secret:
- **Jangan commit ke git**
- Minimal 32 karakter
- Gunakan karakter random
- Berbeda untuk setiap environment (dev, staging, prod)

### Google Client ID:
- Aman untuk di-expose di client-side
- Tetapi tetap simpan di environment variable
- Gunakan berbeda untuk dev dan production

## 7. Troubleshooting

### Error: "Invalid token"
- Pastikan Google Client ID benar
- Cek apakah domain/origin sudah terdaftar di Google Console
- Pastikan ID token belum expired

### Error: "JWT secret not configured"
- Pastikan JWT_SECRET ada di file .env
- Restart server setelah update .env

### Error: "Failed to verify token"
- Cek koneksi internet
- Pastikan Google ID token valid
- Cek format request body

## 8. Environment Variables Template

```env
# Copy ini ke file .env Anda
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=
DB_NAME=lidm

# Generate JWT secret dengan: openssl rand -base64 32
JWT_SECRET=

# Dapatkan dari Google Cloud Console > Credentials
GOOGLE_CLIENT_ID=

PORT=3000
```

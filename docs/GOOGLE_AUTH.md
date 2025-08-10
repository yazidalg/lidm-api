# Google Authentication API Documentation

## Overview
API ini memungkinkan login menggunakan Google OAuth2. Client mengirim ID token yang didapat dari Google Sign-In, dan server akan memverifikasi token tersebut untuk autentikasi.

## Setup

### 1. Environment Variables
Tambahkan ke file `.env`:
```
GOOGLE_CLIENT_ID=your_google_client_id_here
JWT_SECRET=your_jwt_secret_here
```

### 2. Google Client ID
- Buat project di [Google Cloud Console](https://console.cloud.google.com/)
- Enable Google+ API
- Buat OAuth 2.0 Client ID untuk aplikasi web/mobile
- Gunakan Client ID tersebut di environment variable

## API Endpoints

### Google Login
**POST** `/auth/google`

Melakukan autentikasi menggunakan Google ID token.

#### Request Body
```json
{
  "id_token": "google_id_token_from_client"
}
```

#### Response Success (200)
```json
{
  "message": "Login successful",
  "token": "jwt_token_here",
  "user": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "class": "",
    "is_verified": true,
    "profile_picture": "https://lh3.googleusercontent.com/...",
    "point": 0,
    "total_xp": 0,
    "role_id": 1,
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T00:00:00Z"
  }
}
```

#### Response Error (400)
```json
{
  "message": "Invalid request",
  "details": "Error details here"
}
```

#### Response Error (401)
```json
{
  "message": "Invalid Google token",
  "details": "Token verification failed"
}
```

## Client Implementation

### JavaScript (Web)
```javascript
// 1. Include Google API script
<script src="https://apis.google.com/js/api:client.js"></script>

// 2. Initialize Google Sign-In
gapi.load('auth2', function() {
    gapi.auth2.init({
        client_id: 'YOUR_GOOGLE_CLIENT_ID'
    });
});

// 3. Handle sign-in
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
        })
        .catch(error => {
            console.error('Login failed:', error);
        });
    });
}
```

### React Native
```javascript
import { GoogleSignin } from '@react-native-google-signin/google-signin';

// 1. Configure Google Sign-In
GoogleSignin.configure({
    webClientId: 'YOUR_GOOGLE_CLIENT_ID',
});

// 2. Sign in function
const signInWithGoogle = async () => {
    try {
        await GoogleSignin.hasPlayServices();
        const userInfo = await GoogleSignin.signIn();
        const idToken = userInfo.idToken;
        
        // Send to your backend
        const response = await fetch('YOUR_API_URL/auth/google', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                id_token: idToken
            })
        });
        
        const data = await response.json();
        // Store JWT token
        await AsyncStorage.setItem('token', data.token);
        
    } catch (error) {
        console.error('Login failed:', error);
    }
};
```

### Flutter
```dart
import 'package:google_sign_in/google_sign_in.dart';

final GoogleSignIn _googleSignIn = GoogleSignIn(
    clientId: 'YOUR_GOOGLE_CLIENT_ID',
);

Future<void> signInWithGoogle() async {
    try {
        final GoogleSignInAccount? googleUser = await _googleSignIn.signIn();
        final GoogleSignInAuthentication googleAuth = 
            await googleUser!.authentication;
        
        final String? idToken = googleAuth.idToken;
        
        // Send to your backend
        final response = await http.post(
            Uri.parse('YOUR_API_URL/auth/google'),
            headers: {'Content-Type': 'application/json'},
            body: jsonEncode({'id_token': idToken}),
        );
        
        final data = jsonDecode(response.body);
        // Store JWT token
        await storage.write(key: 'token', value: data['token']);
        
    } catch (error) {
        print('Login failed: $error');
    }
}
```

## Security Features

1. **Token Verification**: Server memverifikasi ID token dengan Google
2. **Audience Validation**: Memastikan token dibuat untuk aplikasi ini
3. **Email Verification**: Hanya email yang sudah diverifikasi Google yang diterima
4. **Issuer Validation**: Memastikan token berasal dari Google
5. **Automatic User Creation**: User baru dibuat otomatis jika belum ada

## Flow

1. Client melakukan Google Sign-In
2. Google mengembalikan ID token
3. Client mengirim ID token ke `/auth/google`
4. Server memverifikasi token dengan Google
5. Jika valid, server:
   - Mencari user berdasarkan email
   - Jika tidak ada, buat user baru dengan status verified
   - Generate JWT token
   - Return JWT token dan data user
6. Client menyimpan JWT token untuk request selanjutnya

## Error Handling

- **400**: Request body tidak valid
- **401**: Google token tidak valid atau expired
- **500**: Error internal server (database, JWT generation, dll)

## Notes

- User yang login via Google otomatis ter-verifikasi
- Tidak ada password untuk user Google (field kosong)
- Profile picture dari Google disimpan di database
- Class field mungkin kosong saat pertama login (bisa diupdate nanti)

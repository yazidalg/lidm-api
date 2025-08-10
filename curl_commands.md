# Manual cURL commands untuk test activity tracking

## 1. Logout dulu (clear session)
curl -X POST http://localhost:8080/auth/logout \
  -H "Content-Type: application/json"

## 2. Login dengan akun reguler (jika ada)
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "Email": "user@example.com",
    "Password": "password123"
  }'

## 3. Login dengan Google (butuh valid Google ID token)
curl -X POST http://localhost:8080/auth/google \
  -H "Content-Type: application/json" \
  -d '{
    "id_token": "YOUR_GOOGLE_ID_TOKEN_HERE"
  }'

## 4. Login dengan akun belajar (butuh valid Google ID token dengan domain @belajar.id)
curl -X POST http://localhost:8080/auth/belajar-login \
  -H "Content-Type: application/json" \
  -d '{
    "id_token": "YOUR_BELAJAR_ID_TOKEN_HERE"
  }'

## 5. Setelah login, gunakan token untuk test endpoints (ganti YOUR_TOKEN)
export TOKEN="YOUR_JWT_TOKEN_FROM_LOGIN_RESPONSE"

## 6. Get user profile (akan ditrack sebagai activity)
curl -X GET http://localhost:8080/user/profile \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json"

## 7. Get modules (akan ditrack sebagai module view activity)
curl -X GET http://localhost:8080/module/all \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json"

## 8. Get lessons (akan ditrack sebagai lesson view activity)
curl -X GET http://localhost:8080/lesson/all \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json"

## 9. Get specific module (akan ditrack dengan module_id)
curl -X GET http://localhost:8080/module/1 \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json"

## 10. Get specific lesson (akan ditrack dengan lesson_id)
curl -X GET http://localhost:8080/lesson/1 \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json"

## 11. Check user activities (jika endpoint tersedia)
curl -X GET http://localhost:8080/user-activity/my-activities \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json"

## 12. Logout (akan ditrack sebagai logout activity)
curl -X POST http://localhost:8080/auth/logout \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json"

# Setelah semua test di atas, cek database table user_activities untuk melihat logs!

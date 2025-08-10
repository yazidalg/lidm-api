# Curl commands untuk test Most Active Users

## Pastikan server sudah running di localhost:3000

# 1. Test server
curl -X GET http://localhost:3000/

# 2. Login dulu untuk mendapatkan admin token (ganti dengan kredensial admin yang valid)
curl -X POST http://localhost:3000/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "Email": "admin@example.com",
    "Password": "admin123"
  }'

# 3. Setelah dapat token, ganti YOUR_ADMIN_TOKEN dengan token yang didapat
export ADMIN_TOKEN="YOUR_ADMIN_TOKEN_HERE"

# 4. Get Most Active Users (Admin only) - Ini yang kamu mau lihat!
curl -X GET "http://localhost:3000/user-activity/most-active?limit=10" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" | jq .

# 5. Get Recent Activities  
curl -X GET "http://localhost:3000/user-activity/recent?limit=20" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" | jq .

# 6. Get Activity Stats
curl -X GET "http://localhost:3000/user-activity/stats" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" | jq .

# 7. Get Dashboard (includes most active users data)
curl -X GET "http://localhost:3000/dashboard/" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" | jq .

# 8. Get user's own activities (tidak perlu admin)
curl -X GET "http://localhost:3000/user-activity/my-activities" \
  -H "Authorization: Bearer $ANY_USER_TOKEN" \
  -H "Content-Type: application/json" | jq .

# Output format untuk most active users akan seperti:
# {
#   "message": "Most active users retrieved successfully",
#   "data": [
#     {
#       "user_id": 5,
#       "email": "azisa6980@gmail.com", 
#       "name": "azisa6980",
#       "activity_count": 15,
#       "last_activity": "2025-08-04T08:30:00Z"
#     },
#     {
#       "user_id": 3,
#       "email": "arinzaaa@gmail.com",
#       "name": "arinza", 
#       "activity_count": 8,
#       "last_activity": "2025-08-04T07:15:00Z"
#     }
#   ]
# }

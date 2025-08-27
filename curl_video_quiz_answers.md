# Video Quiz Answer Commands
# Pastikan untuk mengganti {JWT_TOKEN} dengan token autentikasi yang valid

# Base URL
BASE_URL="http://localhost:8080"

# Headers untuk autentikasi
HEADERS="-H \"Content-Type: application/json\" -H \"Authorization: Bearer {JWT_TOKEN}\""

# ============================================
# VIDEO QUIZ COMMANDS FOR MODULE 2 (Quiz ID: 1)
# ============================================

# Menjawab Quiz 1 Module 2 (Timestamp 02:07)
# Question: "Apa fungsi utama akar pada tumbuhan?"
# Correct Answer: A - Menyerap air dan nutrisi dari tanah

# Jawaban BENAR (A)
curl -X POST "$BASE_URL/api/video-quiz/submit" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {JWT_TOKEN}" \
  -d '{
    "video_quiz_id": 1,
    "selected_answer": "A",
    "response_time": 15
  }'

# Jawaban SALAH untuk testing (B)
curl -X POST "$BASE_URL/api/video-quiz/submit" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {JWT_TOKEN}" \
  -d '{
    "video_quiz_id": 1,
    "selected_answer": "B",
    "response_time": 20
  }'

# ============================================
# VIDEO QUIZ COMMANDS FOR MODULE 4
# ============================================

# Quiz 1 Module 4 (Quiz ID: 2, Timestamp 03:47)
# Question: "Apa saja yang dibutuhkan tumbuhan untuk melakukan fotosintesis?"
# Correct Answer: A - Air, Karbondioksida, Matahari, Klorofil

curl -X POST "$BASE_URL/api/video-quiz/submit" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {JWT_TOKEN}" \
  -d '{
    "video_quiz_id": 2,
    "selected_answer": "A",
    "response_time": 12
  }'

# Quiz 2 Module 4 (Quiz ID: 3, Timestamp 04:08)
# Question: "Apa hasil utama dari proses fotosintesis?"
# Correct Answer: B - Oksigen dan Karbohidrat (makanan/glukosa)

curl -X POST "$BASE_URL/api/video-quiz/submit" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {JWT_TOKEN}" \
  -d '{
    "video_quiz_id": 3,
    "selected_answer": "B",
    "response_time": 18
  }'

# Quiz 3 Module 4 (Quiz ID: 4, Timestamp 04:44)
# Question: "Mengapa fotosintesis penting bagi makhluk hidup di Bumi?"
# Correct Answer: B - Karena menghasilkan oksigen dan makanan untuk rantai makanan

curl -X POST "$BASE_URL/api/video-quiz/submit" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {JWT_TOKEN}" \
  -d '{
    "video_quiz_id": 4,
    "selected_answer": "B",
    "response_time": 25
  }'

# ============================================
# GET USER VIDEO QUIZ ANSWERS
# ============================================

# Mendapatkan semua jawaban video quiz user
curl -X GET "$BASE_URL/api/video-quiz/user-answers" \
  -H "Authorization: Bearer {JWT_TOKEN}"

# Mendapatkan jawaban video quiz user untuk video material tertentu
# Module 2 Video Material ID: 3
curl -X GET "$BASE_URL/api/video-quiz/user-answers/3" \
  -H "Authorization: Bearer {JWT_TOKEN}"

# Module 4 Video Material ID: 5
curl -X GET "$BASE_URL/api/video-quiz/user-answers/5" \
  -H "Authorization: Bearer {JWT_TOKEN}"

# ============================================
# EXAMPLE RESPONSES
# ============================================

# Expected Response for Correct Answer:
# {
#   "message": "Video quiz answer submitted successfully",
#   "data": {
#     "video_quiz_id": 1,
#     "selected_answer": "A",
#     "is_correct": true,
#     "answered_at": "2025-08-26T11:30:00Z",
#     "response_time": 15,
#     "question": "Apa fungsi utama akar pada tumbuhan?",
#     "correct_answer": "A",
#     "explanation": "Akar memiliki fungsi utama untuk menyerap air dan nutrisi dari tanah yang diperlukan untuk pertumbuhan dan proses fotosintesis tumbuhan.",
#     "options": {
#       "option_a": "Menyerap air dan nutrisi dari tanah",
#       "option_b": "Menghasilkan cahaya",
#       "option_c": "Menyimpan oksigen",
#       "option_d": "Menghindari hama"
#     }
#   }
# }

# Expected Response for Wrong Answer:
# {
#   "message": "Video quiz answer submitted successfully",
#   "data": {
#     "video_quiz_id": 1,
#     "selected_answer": "B",
#     "is_correct": false,
#     "answered_at": "2025-08-26T11:30:00Z",
#     "response_time": 20,
#     "question": "Apa fungsi utama akar pada tumbuhan?",
#     "correct_answer": "A",
#     "explanation": "Akar memiliki fungsi utama untuk menyerap air dan nutrisi dari tanah yang diperlukan untuk pertumbuhan dan proses fotosintesis tumbuhan.",
#     "options": {
#       "option_a": "Menyerap air dan nutrisi dari tanah",
#       "option_b": "Menghasilkan cahaya",
#       "option_c": "Menyimpan oksigen",
#       "option_d": "Menghindari hama"
#     }
#   }
# }

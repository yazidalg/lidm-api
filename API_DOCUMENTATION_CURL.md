# Complete API Documentation - CURL Commands

## Authentication Endpoints

### 1. Register User
```bash
curl -X POST http://localhost:3000/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123",
    "confirm_password": "password123"
  }'
```

### 2. Login User
```bash
curl -X POST http://localhost:3000/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

### 3. Google Login
```bash
curl -X POST http://localhost:3000/auth/google \
  -H "Content-Type: application/json" \
  -d '{
    "token": "google_oauth_token"
  }'
```

### 4. Belajar Login
```bash5. **RAG Endpoint**: `/user-activity/for-rag` does not require authentication
## Notes:

1. **Authentication**: Most endpoints require a Bearer token obtained from login endpoints
2. **Admin Endpoints**: Marked with (Admin) require admin role
3. **File Uploads**: Use `multipart/form-data` content type
4. **RAG Endpoint**: `/user-activity/for-rag` does not require authentication
5. **Base URL**: Replace `http://localhost:3000` with your actual server URL
6. **Rate Limiting**: Check with your server configuration for rate limits
7. **CORS**: Ensure proper CORS settings for web clients
8. **SubMaterials**: Modern learning content accessed through Module endpoints - see [SUBMATERIAL_DOCUMENTATION.md](./SUBMATERIAL_DOCUMENTATION.md) for complete guide

## SubMaterial Content Types:

SubMaterials within modules support multiple learning content types:

### Video Material
- Embedded video content with duration tracking
- Thumbnails and quality settings
- Progress tracking capabilities

### Quiz Questions  
- Multiple choice questions
- Correct answer validation
- Explanation text for learning
- Different difficulty levels

### AR Experiments
- Augmented Reality experiences
- 3D scene interactions
- Hands-on virtual laboratories
- Interactive learning environments

### Flashcards
- Front/back card structure
- Spaced repetition support
- Category and difficulty classification
- Memory retention tools

## Token Usage:
After successful login, use the returned token in the Authorization header:
```bash
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

## SubMaterial Access Pattern:
```bash
# Get all modules with SubMaterials
curl -X GET http://localhost:3000/module/all \
  -H "Authorization: Bearer YOUR_TOKEN"

# Get specific module with SubMaterials  
curl -X GET http://localhost:3000/module/1 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

For detailed SubMaterial structure, content types, and future endpoint plans, see [SUBMATERIAL_DOCUMENTATION.md](./SUBMATERIAL_DOCUMENTATION.md).: Replace `http://localhost:3000` with your actual server URL
7. **Rate Limiting**: Check with your server configuration for rate limits
8. **CORS**: Ensure proper CORS settings for web clients
9. **SubMaterials**: Modern learning content structure within modules. SubMaterials support multiple content types:
   - **Video Material**: Video pembelajaran dengan durasi dan URL
   - **Quiz Questions**: Pertanyaan kuis interaktif dengan berbagai jenis (multiple choice, etc.)
   - **AR Experiments**: Pengalaman Augmented Reality untuk pembelajaran hands-on
   - **Flashcards**: Kartu belajar untuk menghafal dan review
   
   SubMaterials diakses melalui Module endpoints (`/module/all` dan `/module/:id`) karena merupakan bagian dari struktur module. Lihat [SUBMATERIAL_DOCUMENTATION.md](./SUBMATERIAL_DOCUMENTATION.md) untuk panduan lengkap.
10. **Content Structure**: Modules contain both Lessons (legacy) and SubMaterials (current structure)l -X POST http://localhost:3000/auth/belajar-login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "azis@belajar.id",
    "password": "password123"
  }'
```

### 5. Verify Email
```bash
curl -X GET http://localhost:3000/auth/verify/verification_token_here
```

### 6. Logout
```bash
curl -X POST http://localhost:3000/auth/logout \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Forgot Password Endpoints

### 7. Request Password Reset
```bash
curl -X POST http://localhost:3000/forgot-password \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com"
  }'
```

### 8. Reset Password
```bash
curl -X POST http://localhost:3000/forgot-password/reset \
  -H "Content-Type: application/json" \
  -d '{
    "token": "reset_token",
    "new_password": "newpassword123"
  }'
```

## User Endpoints

### 9. Get All Users (Admin)
```bash
curl -X GET http://localhost:3000/users \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 10. Get User by ID (Admin)
```bash
curl -X GET http://localhost:3000/users/1 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 11. Get My Profile
```bash
curl -X GET http://localhost:3000/users/me \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 12. Update My Profile
```bash
curl -X PUT http://localhost:3000/users/me \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "newusername",
    "email": "newemail@example.com"
  }'
```

### 13. Update User (Admin)
```bash
curl -X PUT http://localhost:3000/users/1 \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "updateduser",
    "email": "updated@example.com"
  }'
```

### 14. Delete User (Admin)
```bash
curl -X DELETE http://localhost:3000/users/1 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Module Endpoints

**Note**: Modules contain SubMaterials (modern learning content structure). See [SUBMATERIAL_DOCUMENTATION.md](./SUBMATERIAL_DOCUMENTATION.md) for complete SubMaterial guide.

### 15. Get All Modules (includes SubMaterials)
```bash
curl -X GET http://localhost:3000/module/all \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Response includes SubMaterials with different content types:**
```json
{
  "message": "Modules retrieved successfully",
  "data": [
    {
      "ID": 1,
      "title": "Introduction to Programming",
      "description": "Basic programming concepts",
      "icon": "/uploads/icons/module1.png",
      "sub_materials": [
        {
          "ID": 1,
          "title": "Welcome Video",
          "description": "Introduction video for the module",
          "order": 1,
          "video_material": {
            "url": "https://example.com/video.mp4",
            "duration": 300,
            "thumbnail": "https://example.com/thumb.jpg"
          }
        },
        {
          "ID": 2,
          "title": "Knowledge Check",
          "description": "Test your understanding",
          "order": 2,
          "quiz_questions": [
            {
              "question": "What is a variable?",
              "question_type": "multiple_choice",
              "options": ["Storage location", "Function", "Loop", "Condition"],
              "correct_answer": "Storage location"
            }
          ]
        },
        {
          "ID": 3,
          "title": "AR Lab Experience",
          "description": "Virtual laboratory",
          "order": 3,
          "ar_experiment": {
            "name": "Virtual Chemistry Lab",
            "description": "Interactive chemistry experiments",
            "scene_url": "https://example.com/ar-scene.json"
          }
        },
        {
          "ID": 4,
          "title": "Key Terms",
          "description": "Important vocabulary",
          "order": 4,
          "flashcards": [
            {
              "front": "API",
              "back": "Application Programming Interface"
            },
            {
              "front": "JSON",
              "back": "JavaScript Object Notation"
            }
          ]
        }
      ]
    }
  ]
}
```

### 16. Get Module by ID (includes SubMaterials)
```bash
curl -X GET http://localhost:3000/module/1 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Returns detailed module with all SubMaterials and their complete content structure.**

### 17. Create Module (Admin)
```bash
curl -X POST http://localhost:3000/module \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Module Title",
    "description": "Module Description",
    "offset_x": 100.5,
    "offset_y": 200.5
  }'
```

### 18. Update Module (Admin)
```bash
curl -X PUT http://localhost:3000/module/1 \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated Module Title",
    "description": "Updated Description",
    "offset_x": 150.5,
    "offset_y": 250.5
  }'
```

### 19. Delete Module (Admin)
```bash
curl -X DELETE http://localhost:3000/module/1 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 20. Upload Module Icon (Admin)
```bash
curl -X POST http://localhost:3000/module/1/icon \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "icon=@/path/to/icon.png"
```

## Lesson Endpoints

### 21. Get All Lessons
```bash
curl -X GET http://localhost:3000/lesson/all \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 22. Get Lesson by ID
```bash
curl -X GET http://localhost:3000/lesson/1 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 23. Create Lesson (Admin)
```bash
curl -X POST http://localhost:3000/lesson \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "module_id": 1,
    "title": "Lesson Title",
    "content": "Lesson content here",
    "sort_order": 1
  }'
```

### 24. Update Lesson (Admin)
```bash
curl -X PUT http://localhost:3000/lesson/1 \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "module_id": 1,
    "title": "Updated Lesson Title",
    "content": "Updated lesson content",
    "sort_order": 2
  }'
```

### 25. Delete Lesson (Admin)
```bash
curl -X DELETE http://localhost:3000/lesson/1 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Question Endpoints

### 26. Get All Questions (Admin)
```bash
curl -X GET http://localhost:3000/question/all \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 27. Get Question by ID (Admin)
```bash
curl -X GET http://localhost:3000/question/1 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 28. Create Question (Admin)
```bash
curl -X POST http://localhost:3000/question \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "lesson_id": 1,
    "question_text": "What is the answer?",
    "question_type": "multiple_choice"
  }'
```

### 29. Update Question (Admin)
```bash
curl -X PUT http://localhost:3000/question/1 \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "lesson_id": 1,
    "question_text": "Updated question?",
    "question_type": "multiple_choice"
  }'
```

### 30. Delete Question (Admin)
```bash
curl -X DELETE http://localhost:3000/question/1 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Answer Endpoints

### 31. Get Answers for Question (Admin)
```bash
curl -X GET http://localhost:3000/answer/question/1 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 32. Create Answer (Admin)
```bash
curl -X POST http://localhost:3000/answer \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "question_id": 1,
    "answer_text": "Answer option",
    "is_correct": true
  }'
```

### 33. Update Answer (Admin)
```bash
curl -X PUT http://localhost:3000/answer/1 \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "question_id": 1,
    "answer_text": "Updated answer",
    "is_correct": false
  }'
```

### 34. Delete Answer (Admin)
```bash
curl -X DELETE http://localhost:3000/answer/1 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Quiz Endpoints

### 35. Get All Quizzes
```bash
curl -X GET http://localhost:3000/quiz/all \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 36. Get Quiz by ID
```bash
curl -X GET http://localhost:3000/quiz/1 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 37. Create Quiz (Admin)
```bash
curl -X POST http://localhost:3000/quiz \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Quiz Title",
    "description": "Quiz Description",
    "lesson_id": 1
  }'
```

### 38. Update Quiz (Admin)
```bash
curl -X PUT http://localhost:3000/quiz/1 \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated Quiz Title",
    "description": "Updated Description",
    "lesson_id": 1
  }'
```

### 39. Delete Quiz (Admin)
```bash
curl -X DELETE http://localhost:3000/quiz/1 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Quiz Session Endpoints

### 40. Create Quiz Session
```bash
curl -X POST http://localhost:3000/quiz-session \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "quiz_id": 1
  }'
```

### 41. Get Quiz Session
```bash
curl -X GET http://localhost:3000/quiz-session/1 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 42. Submit Quiz Answer
```bash
curl -X POST http://localhost:3000/quiz-session/1/answer \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "question_id": 1,
    "answer_id": 1
  }'
```

### 43. Complete Quiz Session
```bash
curl -X POST http://localhost:3000/quiz-session/1/complete \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 44. Get My Quiz Sessions
```bash
curl -X GET http://localhost:3000/quiz-session/my-sessions \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Progress Endpoints

### 45. Get My Progress
```bash
curl -X GET http://localhost:3000/progress/my-progress \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 46. Mark Lesson as Completed
```bash
curl -X POST http://localhost:3000/progress/complete \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "lesson_id": 1
  }'
```

### 47. Get User Progress (Admin)
```bash
curl -X GET http://localhost:3000/progress/user/1 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Prequiz Endpoints

### 48. Create Prequiz Room (Admin)
```bash
curl -X POST http://localhost:3000/prequiz \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "quiz_id": 1,
    "room_name": "Room 1",
    "max_participants": 10
  }'
```

### 49. Get All Prequiz Rooms
```bash
curl -X GET http://localhost:3000/prequiz/all \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 50. Get Prequiz Room by ID
```bash
curl -X GET http://localhost:3000/prequiz/1 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 51. Join Prequiz Room
```bash
curl -X POST http://localhost:3000/prequiz/1/join \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 52. Start Prequiz (Admin)
```bash
curl -X POST http://localhost:3000/prequiz/1/start \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 53. Update Prequiz Room (Admin)
```bash
curl -X PUT http://localhost:3000/prequiz/1 \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "room_name": "Updated Room",
    "max_participants": 15
  }'
```

### 54. Delete Prequiz Room (Admin)
```bash
curl -X DELETE http://localhost:3000/prequiz/1 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Participant Endpoints

### 55. Get All Participants (Admin)
```bash
curl -X GET http://localhost:3000/participant/all \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 56. Get Participant by ID (Admin)
```bash
curl -X GET http://localhost:3000/participant/1 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 57. Create Participant (Admin)
```bash
curl -X POST http://localhost:3000/participant \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "prequiz_id": 1
  }'
```

### 58. Update Participant (Admin)
```bash
curl -X PUT http://localhost:3000/participant/1 \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "score": 85
  }'
```

### 59. Delete Participant (Admin)
```bash
curl -X DELETE http://localhost:3000/participant/1 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Role Endpoints

### 60. Get All Roles (Admin)
```bash
curl -X GET http://localhost:3000/role/all \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 61. Get Role by ID (Admin)
```bash
curl -X GET http://localhost:3000/role/1 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 62. Create Role (Admin)
```bash
curl -X POST http://localhost:3000/role \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "moderator",
    "description": "Moderator role"
  }'
```

### 63. Update Role (Admin)
```bash
curl -X PUT http://localhost:3000/role/1 \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "updated_role",
    "description": "Updated description"
  }'
```

### 64. Delete Role (Admin)
```bash
curl -X DELETE http://localhost:3000/role/1 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## User Activity Endpoints

### 65. Get My Activities
```bash
curl -X GET http://localhost:3000/user-activity/my-activities?limit=20 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 66. Get My Last Activity
```bash
curl -X GET http://localhost:3000/user-activity/my-last \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 67. Get My Streak
```bash
curl -X GET http://localhost:3000/user-activity/my-streak \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 68. Get Recent Activities
```bash
curl -X GET http://localhost:3000/user-activity/recent?limit=50 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 69. Get Most Active Users
```bash
curl -X GET http://localhost:3000/user-activity/most-active?limit=10 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 70. Get Most Active Users Detailed
```bash
curl -X GET http://localhost:3000/user-activity/most-active-detailed \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 71. Get Activities for RAG (No Auth Required)
```bash
curl -X GET http://localhost:3000/user-activity/for-rag?limit=100
```

### 72. Get Activities for RAG with User Filter (No Auth Required)
```bash
curl -X GET http://localhost:3000/user-activity/for-rag?limit=50&user_id=11
```

### 73. Get User Activities (Admin)
```bash
curl -X GET http://localhost:3000/user-activity/users/1?limit=20 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 74. Get Activity Stats (Admin)
```bash
curl -X GET http://localhost:3000/user-activity/stats \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 75. Log Activity Manually (Admin)
```bash
curl -X POST http://localhost:3000/user-activity/log \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "activity_type": "lihat_pelajaran",
    "description": "Manual log",
    "metadata": {
      "source": "manual",
      "admin_logged": true
    }
  }'
```

## Dashboard Endpoints

### 76. Get Dashboard Stats (Admin)
```bash
curl -X GET http://localhost:3000/dashboard/stats \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 77. Get User Analytics (Admin)
```bash
curl -X GET http://localhost:3000/dashboard/user-analytics \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## File Upload Endpoints

### 78. Upload Files
```bash
curl -X POST http://localhost:3000/uploads \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "file=@/path/to/file.jpg"
```

## Flashcard Endpoints

### 79. Get All Flashcards
```bash
curl -X GET http://localhost:3000/flashcard/all \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 80. Get Flashcard by ID
```bash
curl -X GET http://localhost:3000/flashcard/1 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 81. Create Flashcard (Admin)
```bash
curl -X POST http://localhost:3000/flashcard \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "question": "What is Go?",
    "answer": "Programming language",
    "lesson_id": 1
  }'
```

### 82. Update Flashcard (Admin)
```bash
curl -X PUT http://localhost:3000/flashcard/1 \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "question": "Updated question?",
    "answer": "Updated answer",
    "lesson_id": 1
  }'
```

### 83. Delete Flashcard (Admin)
```bash
curl -X DELETE http://localhost:3000/flashcard/1 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Common Headers and Parameters

### Required Headers:
- `Content-Type: application/json` (for POST/PUT requests with JSON body)
- `Authorization: Bearer YOUR_TOKEN` (for authenticated endpoints)

### Common Query Parameters:
- `limit`: Number of results to return (default varies by endpoint)
- `page`: Page number for pagination (where applicable)
- `user_id`: Filter by specific user ID (where applicable)

### Response Format:
All endpoints return JSON with this general structure:
```json
{
  "message": "Success message",
  "data": {
    // Response data here
  }
}
```

### Error Response Format:
```json
{
  "error": "Error message"
}
```

## Important Notes:

1. **Authentication**: Most endpoints require a Bearer token obtained from login endpoints
2. **Admin Endpoints**: Marked with (Admin) require admin role
3. **File Uploads**: Use `multipart/form-data` content type
4. **RAG Endpoint**: `/user-activity/for-rag` does not require authentication
5. **Base URL**: Replace `http://localhost:3000` with your actual server URL
6. **Rate Limiting**: Check with your server configuration for rate limits
7. **CORS**: Ensure proper CORS settings for web clients

## Token Usage:
After successful login, use the returned token in the Authorization header:
```bash
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

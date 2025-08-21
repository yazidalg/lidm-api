# Quiz Answer Documentation

## Overview
This document explains how to submit answers for both **Prequizzes** and **Video Quizzes** in the LIDM API.

## Changes Made
- ✅ **Reduced prequizzes from 10 to 3 per submodule** in the seed data
- ✅ **Quiz answering functionality already exists** - both Prequizzes and Video Quizzes

## API Endpoints

### 1. Submit Prequiz Answer

**Endpoint:** `POST /prequiz/submit`

**Headers:**
```
Authorization: Bearer {your_jwt_token}
Content-Type: application/json
```

**Request Body:**
```json
{
  "prequiz_id": 1,
  "selected_answer": "A"
}
```

**Response Example:**
```json
{
  "message": "Prequiz answer submitted successfully",
  "data": {
    "ID": 1,
    "prequiz_id": 1,
    "user_id": 123,
    "answer": "A",
    "is_correct": true,
    "answered_at": 1692345678
  }
}
```

### 2. Submit Video Quiz Answer

**Endpoint:** `POST /video-quiz/submit`

**Headers:**
```
Authorization: Bearer {your_jwt_token}
Content-Type: application/json
```

**Request Body:**
```json
{
  "video_quiz_id": 1,
  "selected_answer": "B",
  "response_time": 15
}
```

**Response Example:**
```json
{
  "message": "Video quiz answer submitted successfully",
  "data": {
    "ID": 1,
    "video_quiz_id": 1,
    "user_id": 123,
    "selected_answer": "B",
    "is_correct": false,
    "answered_at": 1692345678,
    "response_time": 15
  }
}
```

## Get User Answers

### Get Prequiz Answers
**Endpoint:** `GET /prequiz/user-answers`

**Headers:**
```
Authorization: Bearer {your_jwt_token}
```

### Get Video Quiz Answers
**Endpoint:** `GET /video-quiz/user-answers/{video_material_id}`

**Headers:**
```
Authorization: Bearer {your_jwt_token}
```

### Get All Video Quiz Answers
**Endpoint:** `GET /video-quiz/user-answers`

**Headers:**
```
Authorization: Bearer {your_jwt_token}
```

## Testing with cURL

### Test Prequiz Answer Submission
```bash
curl -X POST http://localhost:8080/prequiz/submit \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "prequiz_id": 1,
    "selected_answer": "A"
  }'
```

### Test Video Quiz Answer Submission
```bash
curl -X POST http://localhost:8080/video-quiz/submit \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "video_quiz_id": 1,
    "selected_answer": "B",
    "response_time": 12
  }'
```

## Answer Options
Valid answer options are: `"A"`, `"B"`, `"C"`, `"D"`

## Error Handling

### Common Errors:
1. **401 Unauthorized** - Missing or invalid JWT token
2. **400 Bad Request** - Invalid request body or missing required fields
3. **409 Conflict** - User has already answered this quiz
4. **404 Not Found** - Quiz not found

### Error Response Example:
```json
{
  "message": "Failed to submit prequiz answer",
  "error": "user has already answered this prequiz"
}
```

## Database Changes Summary

### Before:
- SubMaterial 1: 10 prequizzes
- SubMaterial 2: 10 prequizzes  
- SubMaterial 3: 10 prequizzes

### After:
- SubMaterial 1: 3 prequizzes
- SubMaterial 2: 3 prequizzes
- SubMaterial 3: 3 prequizzes

## How to Re-seed Database
To apply the changes, run:
```bash
# From the project root directory
go run cmd/seed-fotosintesis/main.go
```

This will clear existing data and create new data with only 3 prequizzes per submodule.

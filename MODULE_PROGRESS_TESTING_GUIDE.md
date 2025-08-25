# Module Progress Testing Guide

## Overview
Sistem module progress yang baru mengimplementasikan unlock/lock mechanism untuk modul. Modul 1 selalu unlock, dan menyelesaikan satu modul akan membuka modul berikutnya.

## Architecture Changes

### Database Schema
- **ModuleProgress table**: Tracks unlock status dan progress per user
  - `user_id` + `module_id` unique constraint
  - `is_unlocked` boolean (default true untuk module 1)
  - `is_completed` boolean
  - `progress_percentage` calculated otomatis
  - `unlocked_at`, `completed_at` timestamps

### API Changes
- **GET /module/all**: Sekarang mengembalikan data dengan unlock status dan progress
- Auto-progress update saat user selesai prequiz atau video quiz

### Auto-Progress Logic
1. User jawab prequiz/video quiz dengan benar
2. System calculate progress berdasarkan completed items dalam module
3. Jika progress 100%, module di-mark sebagai completed
4. Auto-unlock module berikutnya (sequential: 1→2→3→...)

## Testing Steps

### 1. Database Setup
```bash
# Run migration to create ModuleProgress table
go run cmd/main.go

# Initialize module progress for existing users
go run cmd/init-module-progress/main.go
```

### 2. Test Module Access Control

#### 2.1 Test Module List API
```bash
# Get modules with unlock status (replace JWT token)
curl -X GET "http://localhost:8080/module/all" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

Expected response:
- Module 1: `is_unlocked: true`
- Module 2+: `is_unlocked: false` (initially)
- All modules show `progress_percentage: 0` (initially)

#### 2.2 Test New User Registration
```bash
# Register new user
curl -X POST "http://localhost:8080/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test User",
    "email": "testuser@example.com",
    "password": "password123",
    "role_id": 2
  }'

# Login to get token
curl -X POST "http://localhost:8080/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "testuser@example.com",
    "password": "password123"
  }'

# Check modules for new user
curl -X GET "http://localhost:8080/module/all" \
  -H "Authorization: Bearer NEW_USER_JWT_TOKEN"
```

Expected: Only module 1 unlocked untuk user baru

### 3. Test Progress Updates

#### 3.1 Complete Prequiz
```bash
# Answer prequiz questions in module 1
curl -X POST "http://localhost:8080/prequiz/answer" \
  -H "Authorization: Bearer JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "prequiz_id": 1,
    "selected_answer": "correct_answer",
    "response_time": 5000
  }'
```

#### 3.2 Complete Video Quiz
```bash
# Answer video quiz in module 1
curl -X POST "http://localhost:8080/video-quiz/answer" \
  -H "Authorization: Bearer JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "video_quiz_id": 1,
    "selected_answer": "correct_answer",
    "response_time": 3000
  }'
```

#### 3.3 Check Progress Update
```bash
# Check modules after answering quizzes
curl -X GET "http://localhost:8080/module/all" \
  -H "Authorization: Bearer JWT_TOKEN"
```

Expected:
- Module 1 progress percentage increased
- When all prequizzes/video quizzes in module 1 completed → progress = 100%
- Module 1 marked as completed
- Module 2 automatically unlocked

### 4. Test Sequential Unlock

Complete all items in module 1:
1. Answer semua prequiz di module 1
2. Answer semua video quiz di module 1
3. Check bahwa module 2 ter-unlock
4. Ulangi untuk module 2 → check module 3 unlock

### 5. Test Edge Cases

#### 5.1 User dengan Existing Progress
```bash
# Check user yang sudah ada progress tidak di-reset
# Run init script multiple times
go run cmd/init-module-progress/main.go
```

#### 5.2 Concurrent Access
Test multiple users completing modules simultaneously

#### 5.3 Invalid Module Access
Try accessing locked modules (should be restricted di frontend tapi backend harus handle)

## Database Verification

### Check ModuleProgress Records
```sql
-- Check module progress for specific user
SELECT mp.*, m.title as module_title, u.name as user_name
FROM module_progress mp
JOIN modules m ON mp.module_id = m.id
JOIN users u ON mp.user_id = u.id
WHERE mp.user_id = 1;

-- Check unlock sequence
SELECT user_id, module_id, is_unlocked, is_completed, progress_percentage, unlocked_at, completed_at
FROM module_progress
ORDER BY user_id, module_id;

-- Check users without progress
SELECT u.id, u.name
FROM users u
LEFT JOIN module_progress mp ON u.id = mp.user_id
WHERE mp.user_id IS NULL;
```

## Expected Behavior

### Module 1 (Always Unlocked)
- ✓ New users: Module 1 auto-unlocked
- ✓ Existing users: Module 1 unlocked after init script
- ✓ Progress calculated from completed prequizzes + video quizzes

### Module 2+ (Sequential Unlock)
- ✓ Initially locked untuk semua users
- ✓ Unlocked only when previous module completed 100%
- ✓ Auto-unlock triggered by progress update

### Progress Calculation
- ✓ Formula: (completed_prequizzes + completed_video_quizzes) / total_items * 100
- ✓ Updated real-time saat user answer quiz
- ✓ Module completed when progress = 100%

### API Responses
- ✓ `/module/all` includes unlock status dan progress
- ✓ Locked modules tidak show detailed content (security)
- ✓ Progress percentage accurate

## Troubleshooting

### Common Issues
1. **ModuleProgress table missing**: Run migration
2. **Existing users no progress**: Run init script
3. **Progress not updating**: Check service wiring in build.go
4. **Modules not unlocking**: Check auto-unlock logic in repository
5. **Wrong progress calculation**: Verify prequiz/video quiz relationships

### Debug Commands
```bash
# Check service wiring
go build ./cmd/main.go

# Check logs for progress updates
tail -f logs/app.log

# Manual progress check
curl -X GET "http://localhost:8080/module/all" -H "Authorization: Bearer TOKEN" | jq
```

## Migration Notes
- ✓ No breaking changes to existing APIs
- ✓ Backward compatible dengan existing data
- ✓ Initialization script safe to run multiple times
- ✓ New users automatically get proper progress setup

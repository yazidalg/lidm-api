# Auto-Unlock Module Feature

## Overview

The auto-unlock module feature automatically unlocks the next module when a user completes all prequizzes and video quizzes in the current module. This provides a seamless learning progression experience.

## How It Works

### 1. Module Completion Detection
When a user submits an answer for a prequiz or video quiz, the system:
- Updates the module progress
- Checks if all quizzes in the current module are completed
- If completed, automatically unlocks the next module

### 2. Completion Criteria
A module is considered completed when:
- **All prequizzes** in the module have been answered by the user
- **All video quizzes** in the module have been answered by the user (if video material exists)

### 3. Auto-Unlock Process
When module completion is detected:
1. Current module is marked as `is_complete: true`
2. Current module progress is set to `100%`
3. `completed_at` timestamp is recorded
4. Next module (by ID order) is automatically unlocked
5. Next module gets `is_unlocked: true` status

## Implementation Details

### Service Methods

#### `ModuleProgressService.CheckAndUnlockNextModule(userID, currentModuleID)`
- Checks if current module is completed
- Finds next module in sequence
- Unlocks next module if current is completed

#### `ModuleProgressService.isModuleCompleted(userID, moduleID)`
- Validates all prequizzes are answered
- Validates all video quizzes are answered (if applicable)
- Returns true only if 100% completion is achieved

### Integration Points

The auto-unlock logic is triggered in:

1. **Prequiz Answer Submission** (`prequiz_service.go`)
   ```go
   // After successful answer submission
   s.moduleProgressService.CheckAndUnlockNextModule(userID, prequiz.ModuleID)
   ```

2. **Video Quiz Answer Submission** (`video_quiz_service.go`)
   ```go
   // After successful answer submission  
   s.moduleProgressService.CheckAndUnlockNextModule(userID, videoQuiz.VideoMaterial.ModuleID)
   ```

### Database Changes

The system uses existing tables:
- `module_progresses` - tracks unlock status and completion
- `prequiz_user_answers` - stores prequiz answers
- `video_quiz_user_answers` - stores video quiz answers

## API Response Changes

### Module List Endpoint (`GET /module/all`)
Now includes:
```json
{
  "data": [
    {
      "id": 1,
      "title": "Module 1",
      "is_unlocked": true,
      "is_complete": true,
      "progress": 100.0,
      "completed_at": "2025-01-01T10:00:00Z"
    },
    {
      "id": 2, 
      "title": "Module 2",
      "is_unlocked": true,  // ← Automatically unlocked
      "is_complete": false,
      "progress": 0.0
    }
  ]
}
```

### Module Detail Endpoint (`GET /module/:id`)
Returns same comprehensive structure with unlock status and progress data.

## Testing

### Manual Testing
Use the provided test script:
```bash
./test_auto_unlock.sh
```

### Test Scenarios

1. **Single Module Completion**
   - Answer all prequizzes in Module 1
   - Verify Module 2 becomes unlocked
   - Verify Module 1 shows 100% progress

2. **Sequential Unlocking** 
   - Complete Module 1 → Module 2 unlocks
   - Complete Module 2 → Module 3 unlocks
   - And so on...

3. **Partial Completion**
   - Answer some (but not all) prequizzes
   - Verify next module remains locked
   - Verify progress is calculated correctly

### API Test Examples

```bash
# Check initial status
curl -H "Authorization: Bearer $TOKEN" http://localhost:3000/module/all

# Answer a prequiz
curl -X POST -H "Authorization: Bearer $TOKEN" \
     -H "Content-Type: application/json" \
     -d '{"prequiz_id": 1, "selected_answer": "A"}' \
     http://localhost:3000/prequiz/submit

# Check updated status
curl -H "Authorization: Bearer $TOKEN" http://localhost:3000/module/all
```

## Benefits

1. **Seamless UX**: No manual intervention needed to unlock content
2. **Progress Validation**: Ensures users complete content before advancing  
3. **Automatic Tracking**: Real-time progress updates and completion status
4. **Scalable**: Works for any number of modules in sequence

## Configuration

The auto-unlock feature is enabled by default and requires no additional configuration. Module order is determined by database ID sequence.

## Error Handling

- If module progress service is unavailable, quiz answers are still recorded
- Unlock failures are logged but don't affect answer submission
- System gracefully handles missing modules or invalid sequences

## Performance Considerations

- Auto-unlock checks run asynchronously for video quizzes (using goroutines)
- Completion validation queries are optimized with proper indexing
- Progress calculations use efficient database queries with joins

# Position Tracking API

When you call the `GET /profile` endpoint, you'll now get two additional fields:

## Response Example:
```json
{
  "message": "User found",
  "data": {
    "id": 1,
    "email": "user@example.com",
    "name": "John Doe",
    "total_xp": 1250,
    "lives": 5,
    // ... all other user fields
    "position_type": "stable",      // This field tracks position type  
    "position_change": 0            // This field tracks position change amount
  },
  "type": "increasing",             // "increasing" or "decreasing"
  "position_change": 2              // 1, 2, 3, etc. (how many positions changed)
}
```

## Field Descriptions:

1. **`type`**: Shows whether the user's leaderboard position is:
   - `"increasing"` - User moved up in ranking
   - `"decreasing"` - User moved down in ranking  
   - `"stable"` - No position change

2. **`position_change`**: Shows how many positions the user moved:
   - `1` - Moved 1 position
   - `2` - Moved 2 positions  
   - `3` - Moved 3 positions
   - etc.

## Database Changes:
Run the migration script `add_position_tracking.sql` to add the required fields to your users table.
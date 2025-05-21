# LIDM Backend

Backend service for LIDM application.

## Socket API Documentation for Quiz Matchmaking

The backend provides a WebSocket-based API for real-time quiz matchmaking and gameplay. This document describes how to integrate with this API from a mobile application.

### Connection

Connect to the WebSocket endpoint with user parameters:

```
ws://{server_url}/ws/matchmaking?user_id={user_id}&username={username}
```

Parameters:
- `user_id`: The unique identifier for the user
- `username`: The display name of the user

Example:
```
ws://example.com/ws/matchmaking?user_id=123&username=Player1
```

### Message Format

All messages follow this JSON format:

```json
{
  "type": "event_type",
  "payload": {
    // Event-specific data
  }
}
```

### Client-to-Server Events

1. **Find Match**
   ```json
   {
     "type": "find_match",
     "payload": {}
   }
   ```

2. **Cancel Matchmaking**
   ```json
   {
     "type": "cancel_matchmaking",
     "payload": {}
   }
   ```

3. **Submit Answer**
   ```json
   {
     "type": "answer_question",
     "payload": {
       "question_id": 123,
       "answer": "A"  // Can be "A", "B", "C", or "D"
     }
   }
   ```

### Server-to-Client Events

1. **Connected Confirmation**
   ```json
   {
     "type": "connected",
     "payload": {
       "message": "Connected to quiz matchmaking",
       "user_id": 123
     }
   }
   ```

2. **Matchmaking Status**
   ```json
   {
     "type": "matchmaking_status",
     "payload": {
       "status": "queued",
       "message": "Waiting for an opponent"
     }
   }
   ```

3. **Match Started**
   ```json
   {
     "type": "match_started",
     "payload": {
       "session_id": "session_12345",
       "opponent": {
         "user_id": 456,
         "username": "Player2"
       },
       "quiz_id": 789
     }
   }
   ```

4. **New Question**
   ```json
   {
     "type": "new_question",
     "payload": {
       "question_id": 123,
       "text": "What is the capital of France?",
       "options": {
         "A": "London",
         "B": "Paris",
         "C": "Berlin",
         "D": "Madrid"
       },
       "read_time": 5,
       "answer_time": 10,
       "question_number": 1,
       "total_questions": 5
     }
   }
   ```

5. **Start Answering**
   ```json
   {
     "type": "start_answering",
     "payload": {
       "question_id": 123,
       "time_remaining": 10
     }
   }
   ```

6. **Answer Result**
   ```json
   {
     "type": "answer_result",
     "payload": {
       "question_id": 123,
       "is_correct": true,
       "score": 10
     }
   }
   ```

7. **Question Ended**
   ```json
   {
     "type": "question_ended",
     "payload": {
       "question_id": 123,
       "correct_answer": "B",
       "explanation": "Paris is the capital of France",
       "scores": {
         "123": 10,
         "456": 0
       }
     }
   }
   ```

8. **Quiz Completed**
   ```json
   {
     "type": "quiz_completed",
     "payload": {
       "final_scores": {
         "123": 30,
         "456": 20
       },
       "winner": {
         "user_id": 123,
         "username": "Player1",
         "score": 30
       }
     }
   }
   ```

9. **Opponent Disconnected**
   ```json
   {
     "type": "opponent_disconnected",
     "payload": {
       "message": "Your opponent has disconnected"
     }
   }
   ```

10. **Error**
    ```json
    {
      "type": "error",
      "payload": {
        "message": "Error message"
      }
    }
    ```

### Game Flow

1. User connects to the WebSocket
2. User sends a "find_match" message
3. Server responds with "matchmaking_status" (queued)
4. When opponent is found, server sends "match_started" to both players
5. Server starts sending "new_question" messages with questions
6. After read time, server sends "start_answering" message
7. Players submit answers with "answer_question" message
8. Server responds with "answer_result" to each player
9. After answer time, server sends "question_ended" with correct answer and scores
10. Steps 5-9 repeat for each question
11. After all questions, server sends "quiz_completed" with final results
12. Session ends

### Scoring

- Players receive 10 points for each correct answer
- No points are awarded for incorrect answers
- The player with the highest score at the end is declared the winner
- In case of a tie, the result is a draw

### Implementation Notes for Mobile Developers

1. Use a WebSocket library for your platform (e.g., OkHttp for Android, Starscream for iOS)
2. Maintain the connection and handle reconnects
3. Implement timeouts for each phase of the quiz
4. Update the UI in real-time based on the server events
5. Display a loading indicator during matchmaking
6. Show a countdown during read_time and answer_time
7. Animate the transition between questions
8. Handle disconnection gracefully

### Testing

You can test the socket connection using the web interface at:
```
http://{server_url}/ws-test
```

This page provides a simple UI to test the WebSocket functionality. 
# Socket.IO Multiplayer Usage Guide

Endpoint base: `/socket.io`
Transport: WebSocket (library will fallback automatically if needed)

## Events

Client -> Server:
- `create_lobby` [module_id:number, host_user_id:number, question_count?]
- `join_lobby` [invite_code:string, user_id:number]
- `join_quiz` [quiz_id:number, user_id:number] (direct start if already have quiz id)
- `submit_answer` [quiz_id:number, question_id:number, option:string, user_id:number]

Server -> Client:
- `lobby_created` { quiz_id, invite_code, module_id }
- `lobby_joined` { quiz_id, user_id }
- `lobby_full` { invite_code }
- `lobby_not_found` { invite_code }
- `user_join` { user_id, room }
- `start_quiz` { quiz_id, total_questions }
- `question` Question JSON (id, content, options, etc.)
- `answer_result` { user_id, question_id, is_correct, score }
- `question_ended` { question_id, correct_answer }
- `quiz_completed` { scores: {user_id:score,...}, winner }
- `user_leave` { user_id }

Single player extra:
- `lives_exhausted`

## Flow (Realtime Lobby + Multiplayer)
1. Host emit `create_lobby`.
2. Server balas `lobby_created` => tampilkan invite_code ke host.
3. Opponent emit `join_lobby` dengan invite_code.
4. Server balas `lobby_joined`; jika pemain ke-2 masuk, status quiz diubah ke in_progress dan otomatis start (server fetch pertanyaan & emit `start_quiz` + `question`).
5. Jawab soal pakai `submit_answer`.
6. Siklus soal seperti sebelumnya.

Jika sudah punya quiz_id (misal dari REST lama) bisa langsung gunakan `join_quiz`.

## Time Handling
- Per question timer = (answer_time + read_time) seconds from Question model; defaults to 15s if zero.

## Example (JavaScript client)
```js
import { io } from 'socket.io-client';
const socket = io('http://localhost:3000', { path: '/socket.io' });

socket.on('connect', ()=> console.log('connected'));

socket.emit('join_quiz', quizId, userId);

socket.on('start_quiz', data => console.log('Start', data));

socket.on('question', q => {
  console.log('Question', q);
  // choose option
  const option = q.option_a; // example
  socket.emit('submit_answer', q.quiz_id || quizId, q.id, option, userId);
});

socket.on('answer_result', r => console.log('Answer Result', r));
socket.on('question_ended', d => console.log('Ended', d));
socket.on('quiz_completed', d => console.log('Completed', d));
```

## Example (JavaScript lobby)
```js
socket.emit('create_lobby', moduleId, hostUserId, 5);
socket.on('lobby_created', d => console.log('Invite', d.invite_code));
// opponent side
socket.emit('join_lobby', inviteCode, otherUserId);
```

## Flutter (dart socket_io_client)
```dart
import 'package:socket_io_client/socket_io_client.dart' as IO;

final socket = IO.io('http://localhost:3000', IO.OptionBuilder()
  .setPath('/socket.io')
  .setTransports(['websocket'])
  .build());

socket.onConnect((_) {
  socket.emit('join_quiz', [quizId, userId]);
});

socket.on('question', (q) {
  // choose option
  final option = q['option_a'];
  socket.emit('submit_answer', [quizId, q['id'], option, userId]);
});
```

## Example (Flutter lobby)
```dart
socket.emit('create_lobby', [moduleId, hostUserId, 5]);
// listen lobby_created
socket.emit('join_lobby', [inviteCode, userId]);
```

## Notes / Limitations
- Auth not yet enforced on Socket.IO handshake (pass JWT next iteration).
- Random EXP boost not yet applied in multiplayer path.
- Lives mechanic only for single player.
- Room removed only when server restarts (no explicit cleanup after finish yet).

## Next Steps (Planned)
- JWT verification in connection query or auth event.
- Random EXP boost parity for multiplayer, adjustable via env.
- Automatic room cleanup & reconnection handling.
- Migration of remaining legacy WebSocket features.

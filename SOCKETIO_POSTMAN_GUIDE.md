# Socket.IO Multiplayer Testing Guide (Postman & CLI)

Dokumen ini memandu end-to-end test flow multiplayer quiz via Socket.IO: dari auth (JWT), membuat lobby, join, sampai quiz berjalan.

## 1. Prasyarat
- Server jalan (default: `:3000`).
- Minimal 2 user terdaftar & terverifikasi (atau pakai admin + 1 user lain).
- Module sudah ada (misal module id untuk quiz). Jika tidak, buat via endpoint `/module/create` (role admin).

## 2. Dapatkan JWT Token (REST)
Login (contoh):
```
POST /auth/login
{
  "email": "user1@example.com",
  "password": "password123"
}
```
Response (potong): `{ "token": "<JWT>" }`
Lakukan untuk kedua user (user1, user2). Simpan kedua token.

Saat ini Socket.IO server belum memverifikasi token di handshake (TO DO). Sementara, tetap pakai token untuk REST create data. (Jika nanti auth handshake diterapkan, tambahkan query param `?token=<JWT>` di URL Socket.IO.)

## 3. Postman: Membuat WebSocket (Socket.IO)
Postman (versi terbaru) mendukung WebSocket; pilih New > WebSocket Request.
Masukkan URL:
```
ws://localhost:3000/socket.io/?EIO=4&transport=websocket
```
Catatan: Postman kadang belum handshake penuh versi Socket.IO; bila gagal, gunakan alternatif:
- Client JS kecil (Node)
- Flutter app
- `socket.io-client` via REPL

Contoh Node quick test (opsional):
```js
import { io } from 'socket.io-client';
const socket = io('http://localhost:3000', { path: '/socket.io', transports:['websocket'] });
```

## 4. Event Referensi Multiplayer
Client -> Server:
- `create_lobby` [module_id:number, host_user_id:number, question_count?]
- `join_lobby` [invite_code:string, user_id:number]
- `join_quiz` [quiz_id:number, user_id:number] (jika sudah punya quiz id)
- `submit_answer` [quiz_id:number, question_id:number, option:string, user_id:number]

Server -> Client:
- `lobby_created` { quiz_id, invite_code, module_id }
- `lobby_joined` { quiz_id, user_id }
- `lobby_full` { invite_code }
- `lobby_not_found` { invite_code }
- `user_join` { user_id, room }
- `start_quiz` { quiz_id, total_questions }
- `question` { id, question, options:{optionA..}, answer_time, read_time, correct_answer (tidak sebelum selesai), ... }
- `answer_result` { user_id, question_id, is_correct, score }
- `question_ended` { question_id, correct_answer }
- `quiz_completed` { scores:{ user_id:score }, winner }
- `user_leave` { user_id }

Single Player tambahan:
- `lives_exhausted`

## 5. Alur Lengkap Multiplayer (Step-by-Step)
1. Host connect socket.
2. Host emit `create_lobby` dengan module id & host user id:
   - Payload (array): `[12, 5, 10]` artinya module_id=12, host_user_id=5, question_count=10.
3. Terima `lobby_created` -> catat `quiz_id` & `invite_code`.
4. Player kedua connect socket.
5. Player kedua emit `join_lobby` dengan invite code: `["AB12CD34", 7]` (user_id=7).
6. Host & player kedua sama-sama menerima `lobby_joined` & `user_join`.
7. Ketika pemain kedua bergabung, server otomatis set status quiz `in_progress`, fetch pertanyaan, emit `start_quiz` lalu `question` pertama.
8. Untuk menjawab, tiap pemain emit `submit_answer`:
   - `[quiz_id, question_id, option, user_id]` contoh: `[21, 101,"A",5]`.
9. Server broadcast `answer_result` untuk tiap jawaban.
10. Setelah semua jawab atau timer habis, dapat `question_ended` dan lanjut ke soal berikutnya.
11. Terakhir, dapat `quiz_completed` berisi skor & winner.

## 6. Format Payload di Postman (WebSocket Tab)
Karena Socket.IO memakai framing sendiri, Postman kadang perlu mode Socket.IO (jika tersedia). Bila tidak, gunakan Node/Flutter client.
Jika Postman mendukung event, masukkan JSON sesuai library UI-nya. Bila tidak, berikut contoh Node script lebih pasti.

### Contoh Node Script (Multiplayer)
```js
import { io } from 'socket.io-client';
const hostUserId = 5;
const moduleId = 12;
const opponentUserId = 7;

const host = io('http://localhost:3000', { path:'/socket.io', transports:['websocket'] });
const opponent = io('http://localhost:3000', { path:'/socket.io', transports:['websocket'] });

host.on('connect', ()=> {
  console.log('Host connected');
  host.emit('create_lobby', moduleId, hostUserId, 10);
});

host.on('lobby_created', d => {
  console.log('Lobby created', d);
  opponent.emit('join_lobby', d.invite_code, opponentUserId);
});

opponent.on('lobby_joined', d => console.log('Opponent joined quiz', d));

host.on('start_quiz', d => console.log('Start quiz', d));

function answerLogic(socket, userId) {
  socket.on('question', q => {
    // simple always choose OptionA
    socket.emit('submit_answer', q.quiz_id || d.quiz_id, q.id, 'A', userId);
  });
}
answerLogic(host, hostUserId);
answerLogic(opponent, opponentUserId);

host.on('answer_result', r => console.log('[HOST] result', r));
opponent.on('answer_result', r => console.log('[OPP] result', r));

host.on('quiz_completed', d => console.log('Completed', d));
```

## 7. Testing Manual (Urutan Event yang Diharapkan)
| Langkah | Host Menerima | Opponent Menerima |
|--------|---------------|-------------------|
| Create Lobby | lobby_created | - |
| Opponent Join | user_join, lobby_joined | lobby_joined, user_join |
| Start | start_quiz, question | start_quiz, question |
| Jawab Soal | answer_result (tiap pemain) | answer_result |
| Selesai Soal | question_ended | question_ended |
| Akhir Quiz | quiz_completed | quiz_completed |

## 8. Troubleshooting
| Masalah | Penyebab | Solusi |
|--------|----------|--------|
| Tidak ada respon setelah create_lobby | Module ID salah | Verifikasi module id via GET /module/all |
| join_lobby balas lobby_not_found | Invite code salah / quiz bukan pending | Pastikan status pending & code benar |
| quiz tidak start setelah join kedua | Player kedua tidak tercatat | Cek log server, pastikan event join_lobby terkirim |
| Postman gagal connect | Versi Postman belum dukung Socket.IO penuh | Pakai node script / Flutter sementara |
| answer_result tidak muncul untuk pemain lain | Salah parameter quiz_id / question_id | Log argumen emit untuk validasi |

## 9. Rencana Pengembangan (Next)
- Validasi JWT di handshake (query token).
- Random EXP boost untuk multiplayer.
- Cleanup room setelah selesai quiz.
- Rejoin logic & reconnection safe.

## 10. Ringkas Payload Contoh
```
create_lobby -> [12, 5, 10]
(lobby_created) <- {"quiz_id":21,"invite_code":"AB12CD34","module_id":12}
join_lobby -> ["AB12CD34",7]
(lobby_joined) <- {"quiz_id":21,"user_id":7}
(start_quiz) <- {"quiz_id":21,"total_questions":10}
(question) <- {"id":101,"question":"Sample question ...","options":{"optionA":"A","optionB":"B","optionC":"C","optionD":"D"},"answer_time":5,"read_time":5}
submit_answer -> [21,101,"A",5]
(answer_result) <- {"user_id":5,"question_id":101,"is_correct":true,"score":10}
(question_ended) <- {"question_id":101,"correct_answer":"A"}
(quiz_completed) <- {"scores":{"5":80,"7":70},"winner":"User Host"}
```

---
Jika ingin template koleksi Postman (JSON) bisa ditambahkan; beri tahu kalau perlu.

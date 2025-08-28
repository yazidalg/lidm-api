# Leaderboard API Documentation

## Overview
API untuk menampilkan leaderboard berdasarkan point dari solo quiz dan matchmaking quiz dengan format yang diinginkan:

```json
{
  "juara1": {}, // user dan data point
  "juara2": {}, // user dan data point  
  "juara3": {}, // user dan data point
  "leaderboard": [] // juara 4 seterusnya
}
```

## Endpoints

### 1. Get Leaderboard
**GET** `/leaderboard`

Menampilkan leaderboard dengan top 3 dan sisanya.

#### Query Parameters
- `module_id` (optional): Filter berdasarkan module tertentu
- `quiz_type` (optional): Filter berdasarkan tipe quiz ("solo" atau "matchmaking")

#### Example Request
```bash
# Get all leaderboard
curl -X GET "http://localhost:3000/leaderboard" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Get leaderboard untuk module 1
curl -X GET "http://localhost:3000/leaderboard?module_id=1" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Get leaderboard untuk solo quiz
curl -X GET "http://localhost:3000/leaderboard?quiz_type=solo" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Get leaderboard untuk matchmaking quiz di module 1
curl -X GET "http://localhost:3000/leaderboard?module_id=1&quiz_type=matchmaking" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Example Response
```json
{
  "juara1": {
    "user": {
      "id": 1,
      "name": "Ahmad",
      "email": "ahmad@example.com",
      "avatar": "avatar.jpg"
    },
    "score": 250,
    "rank": 1,
    "is_current_user": false
  },
  "juara2": {
    "user": {
      "id": 2,
      "name": "Siti",
      "email": "siti@example.com",
      "avatar": "avatar2.jpg"
    },
    "score": 230,
    "rank": 2,
    "is_current_user": true
  },
  "juara3": {
    "user": {
      "id": 3,
      "name": "Budi",
      "email": "budi@example.com",
      "avatar": "avatar3.jpg"
    },
    "score": 200,
    "rank": 3,
    "is_current_user": false
  },
  "leaderboard": [
    {
      "user": {
        "id": 4,
        "name": "Citra",
        "email": "citra@example.com",
        "avatar": "avatar4.jpg"
      },
      "score": 180,
      "rank": 4,
      "is_current_user": false
    },
    {
      "user": {
        "id": 5,
        "name": "Deni",
        "email": "deni@example.com",
        "avatar": "avatar5.jpg"
      },
      "score": 150,
      "rank": 5,
      "is_current_user": false
    }
  ]
}
```

### 2. Get User Rank
**GET** `/leaderboard/user/{user_id}`

Mendapatkan ranking dan score spesifik user.

#### Path Parameters
- `user_id`: ID user yang ingin dicari rankingnya

#### Query Parameters
- `module_id` (optional): Filter berdasarkan module tertentu
- `quiz_type` (optional): Filter berdasarkan tipe quiz ("solo" atau "matchmaking")

#### Example Request
```bash
# Get rank user ID 5
curl -X GET "http://localhost:3000/leaderboard/user/5" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Get rank user ID 5 untuk module 1
curl -X GET "http://localhost:3000/leaderboard/user/5?module_id=1" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Example Response
```json
{
  "user": {
    "id": 5,
    "name": "Deni",
    "email": "deni@example.com",
    "avatar": "avatar5.jpg"
  },
  "score": 150,
  "rank": 5,
  "is_current_user": true
}
```

## Data Source

Leaderboard mengambil data dari:
- **All Users**: Semua user di sistem akan muncul di leaderboard
- **Score Calculation**: Total score dihitung dari table `participants`
- **Field Score**: `total_score` dari semua quiz yang `is_finished = true`
- **Default Score**: User yang belum pernah ikut quiz akan memiliki score 0
- **Grouping**: Per user (total score dari semua quiz yang sudah diselesaikan)
- **Sorting**: Descending berdasarkan total score

## Leaderboard Logic

1. **All Users Included**: Semua user di sistem akan muncul di leaderboard, tidak hanya yang pernah ikut quiz
2. **Score Calculation**: 
   - User yang sudah ikut quiz: Total dari semua `total_score` di table `participants`
   - User yang belum pernah ikut quiz: Score = 0
3. **Ranking**: Berdasarkan total score (descending)
4. **Tie Breaking**: User dengan score sama akan mendapat rank yang sama

## Authentication

Semua endpoint memerlukan JWT token yang valid dalam header `Authorization: Bearer <token>`.

## Error Responses

### 401 Unauthorized
```json
{
  "error": "Unauthorized",
  "message": "Invalid or missing token"
}
```

### 400 Bad Request
```json
{
  "error": "Invalid user ID"
}
```

### 500 Internal Server Error
```json
{
  "error": "Failed to get leaderboard",
  "message": "Database connection error"
}
```

## Implementation Notes

1. **All Users Displayed**: Semua user di sistem akan ditampilkan di leaderboard, termasuk yang belum pernah ikut quiz
2. **Scoring System**: Menggunakan `total_score` dari table `participants` untuk quiz yang sudah selesai
3. **Zero Score Handling**: User yang belum pernah ikut quiz akan mendapat score 0
4. **Quiz Types**: 
   - `solo`: Quiz single player
   - `matchmaking`: Quiz multiplayer/head-to-head
5. **Module Filtering**: Dapat filter berdasarkan module ID tertentu
6. **Ranking Logic**: 
   - Rank 1-3 masuk ke `juara1`, `juara2`, `juara3`
   - Rank 4+ masuk ke array `leaderboard`
7. **User Aggregation**: Total score user dijumlahkan dari semua quiz yang pernah diselesaikan
8. **Current User Identification**: Field `is_current_user` menandakan user yang sedang login
   - `true`: User yang sedang request API ini
   - `false`: User lain di leaderboard
   - Berguna untuk highlight posisi user di frontend

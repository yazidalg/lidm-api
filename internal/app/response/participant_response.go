package response

import "time"

type ParticipantResponse struct {
	ID         uint      `json:"id"`
	QuizID     uint      `json:"quiz_id"`
	UserID     uint      `json:"user_id"`
	TotalScore int       `json:"total_score"`
	CreatedAt  time.Time `json:"created_at"`
	User       User      `json:"user"`
	Quiz       Quiz      `json:"quiz"`
}

type Quiz struct {
	ID     uint   `json:"id"`
	Status string `json:"status"`
}

type User struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Class string `json:"class"`
}

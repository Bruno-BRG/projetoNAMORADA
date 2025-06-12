package models

import (
	"time"
)

type Question struct {
	ID            int       `json:"id" db:"id"`
	Title         string    `json:"title" db:"title"`
	Content       string    `json:"content" db:"content"`
	Options       []string  `json:"options"` // JSON array
	CorrectAnswer int       `json:"correct_answer" db:"correct_answer"`
	Reward        string    `json:"reward" db:"reward"`
	ScheduledAt   time.Time `json:"scheduled_at" db:"scheduled_at"`
	IsActive      bool      `json:"is_active" db:"is_active"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

type UserResponse struct {
	ID         int       `json:"id" db:"id"`
	QuestionID int       `json:"question_id" db:"question_id"`
	Answer     int       `json:"answer" db:"answer"`
	IsCorrect  bool      `json:"is_correct" db:"is_correct"`
	AnsweredAt time.Time `json:"answered_at" db:"answered_at"`
}

type User struct {
	ID        int       `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	IsAdmin   bool      `json:"is_admin" db:"is_admin"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type QuizSession struct {
	UserID        int
	CurrentQuiz   *Question
	CanAnswer     bool
	NextQuizTime  *time.Time
	TotalCorrect  int
	TotalAnswered int
}

type QuizProgress struct {
	Total    int `json:"total"`
	Answered int `json:"answered"`
	Correct  int `json:"correct"`
}

type QuizStatus struct {
	HasCurrent    bool         `json:"has_current"`
	Current       *Question    `json:"current,omitempty"`
	HasNext       bool         `json:"has_next"`
	NextAt        time.Time    `json:"next_at,omitempty"`
	TimeUntilNext string       `json:"time_until_next,omitempty"`
	Progress      QuizProgress `json:"progress"`
}

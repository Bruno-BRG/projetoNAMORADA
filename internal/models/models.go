package models

import (
	"time"
)

// User representa um usuário do sistema (você e sua namorada)
type User struct {
	ID          int        `json:"id" db:"id"`
	Username    string     `json:"username" db:"username"`
	Password    string     `json:"-" db:"password_hash"` // nunca exposer a senha no JSON
	Role        string     `json:"role" db:"role"`       // "admin" ou "user"
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	LastLoginAt *time.Time `json:"last_login_at" db:"last_login_at"`
}

// Quiz representa um quiz completo
type Quiz struct {
	ID          int        `json:"id" db:"id"`
	Title       string     `json:"title" db:"title"`
	Description string     `json:"description" db:"description"`
	ScheduledAt time.Time  `json:"scheduled_at" db:"scheduled_at"`
	IsActive    bool       `json:"is_active" db:"is_active"`
	Reward      string     `json:"reward" db:"reward"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	Questions   []Question `json:"questions,omitempty"`
}

// Question representa uma pergunta do quiz
type Question struct {
	ID      int      `json:"id" db:"id"`
	QuizID  int      `json:"quiz_id" db:"quiz_id"`
	Text    string   `json:"text" db:"text"`
	Options []string `json:"options"` // JSON array no banco
	Correct int      `json:"correct" db:"correct_option"`
	Order   int      `json:"order" db:"question_order"`
}

// UserAnswer representa a resposta de um usuário
type UserAnswer struct {
	ID         int       `json:"id" db:"id"`
	UserID     int       `json:"user_id" db:"user_id"`
	QuizID     int       `json:"quiz_id" db:"quiz_id"`
	QuestionID int       `json:"question_id" db:"question_id"`
	Answer     int       `json:"answer" db:"selected_option"`
	IsCorrect  bool      `json:"is_correct" db:"is_correct"`
	AnsweredAt time.Time `json:"answered_at" db:"answered_at"`
}

// QuizAttempt representa uma tentativa completa de quiz
type QuizAttempt struct {
	ID             int       `json:"id" db:"id"`
	UserID         int       `json:"user_id" db:"user_id"`
	QuizID         int       `json:"quiz_id" db:"quiz_id"`
	Score          int       `json:"score" db:"score"`
	TotalQuestions int       `json:"total_questions" db:"total_questions"`
	CompletedAt    time.Time `json:"completed_at" db:"completed_at"`
	IPAddress      string    `json:"ip_address" db:"ip_address"`
}

// RateLimit para controle de tentativas
type RateLimit struct {
	ID          int       `json:"id" db:"id"`
	IPAddress   string    `json:"ip_address" db:"ip_address"`
	QuizID      int       `json:"quiz_id" db:"quiz_id"`
	Attempts    int       `json:"attempts" db:"attempts"`
	LastAttempt time.Time `json:"last_attempt" db:"last_attempt"`
	ResetAt     time.Time `json:"reset_at" db:"reset_at"`
}

// Session para controle de autenticação
type Session struct {
	ID        string    `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	IPAddress string    `json:"ip_address" db:"ip_address"`
	UserAgent string    `json:"user_agent" db:"user_agent"`
}

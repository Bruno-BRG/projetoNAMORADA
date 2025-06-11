package models

import (
	"time"
)

// User representa um usuário do sistema
type User struct {
	ID       int       `json:"id" db:"id"`
	Username string    `json:"username" db:"username"`
	Password string    `json:"-" db:"password"` // Nunca expor senha no JSON
	Role     string    `json:"role" db:"role"`  // "admin" ou "visitor"
	CreateAt time.Time `json:"created_at" db:"created_at"`
}

// Question representa uma pergunta do quiz
type Question struct {
	ID            int       `json:"id" db:"id"`
	Title         string    `json:"title" db:"title"`
	Description   string    `json:"description" db:"description"`
	Options       []string  `json:"options"` // JSON array no banco
	CorrectAnswer string    `json:"correct_answer" db:"correct_answer"`
	Reward        string    `json:"reward" db:"reward"`
	ScheduledAt   time.Time `json:"scheduled_at" db:"scheduled_at"`
	IsActive      bool      `json:"is_active" db:"is_active"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

// Answer representa uma resposta do usuário
type Answer struct {
	ID         int       `json:"id" db:"id"`
	UserID     int       `json:"user_id" db:"user_id"`
	QuestionID int       `json:"question_id" db:"question_id"`
	Answer     string    `json:"answer" db:"answer"`
	IsCorrect  bool      `json:"is_correct" db:"is_correct"`
	AnsweredAt time.Time `json:"answered_at" db:"answered_at"`
}

// QuizSession representa uma sessão de quiz
type QuizSession struct {
	ID             int        `json:"id" db:"id"`
	UserID         int        `json:"user_id" db:"user_id"`
	TotalQuestions int        `json:"total_questions" db:"total_questions"`
	CorrectAnswers int        `json:"correct_answers" db:"correct_answers"`
	CompletedAt    *time.Time `json:"completed_at" db:"completed_at"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
}

// QuestionWithStatus inclui o status de disponibilidade da pergunta
type QuestionWithStatus struct {
	Question
	IsAvailable bool `json:"is_available"`
	IsAnswered  bool `json:"is_answered"`
}

// LoginRequest representa o payload de login
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// CreateQuestionRequest representa o payload para criar uma pergunta
type CreateQuestionRequest struct {
	Title         string   `json:"title" binding:"required"`
	Description   string   `json:"description"`
	Options       []string `json:"options" binding:"required,min=2"`
	CorrectAnswer string   `json:"correct_answer" binding:"required"`
	Reward        string   `json:"reward" binding:"required"`
	ScheduledAt   string   `json:"scheduled_at" binding:"required"` // ISO string, será parseado
}

// AnswerRequest representa o payload para responder uma pergunta
type AnswerRequest struct {
	QuestionID int    `json:"question_id" binding:"required"`
	Answer     string `json:"answer" binding:"required"`
}

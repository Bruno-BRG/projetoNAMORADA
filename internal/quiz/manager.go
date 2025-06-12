package quiz

import (
	"database/sql"
	"time"
	"valentine-quiz/internal/models"
)

type QuizManager struct {
	db *sql.DB
}

func NewQuizManager(db *sql.DB) *QuizManager {
	return &QuizManager{db: db}
}

// GetAvailableQuiz retorna o quiz atual disponível baseado no horário
func (qm *QuizManager) GetAvailableQuiz() (*models.Question, error) {
	now := time.Now()
	var question models.Question

	err := qm.db.QueryRow(`
		SELECT id, title, content, options, correct_answer, reward, scheduled_at 
		FROM questions 
		WHERE scheduled_at <= ? AND is_active = 1
		ORDER BY scheduled_at DESC 
		LIMIT 1
	`, now).Scan(&question.ID, &question.Title, &question.Content,
		&question.Options, &question.CorrectAnswer, &question.Reward, &question.ScheduledAt)

	if err != nil {
		return nil, err
	}

	return &question, nil
}

// GetNextQuiz retorna o próximo quiz agendado
func (qm *QuizManager) GetNextQuiz() (*models.Question, error) {
	now := time.Now()
	var question models.Question

	err := qm.db.QueryRow(`
		SELECT id, title, content, options, correct_answer, reward, scheduled_at 
		FROM questions 
		WHERE scheduled_at > ? AND is_active = 1
		ORDER BY scheduled_at ASC 
		LIMIT 1
	`, now).Scan(&question.ID, &question.Title, &question.Content,
		&question.Options, &question.CorrectAnswer, &question.Reward, &question.ScheduledAt)

	if err != nil {
		return nil, err
	}

	return &question, nil
}

// HasUserAnswered verifica se o usuário já respondeu a pergunta
func (qm *QuizManager) HasUserAnswered(questionID int, userID string) bool {
	var count int
	qm.db.QueryRow(`
		SELECT COUNT(*) FROM user_responses 
		WHERE question_id = ? AND user_id = ?
	`, questionID, userID).Scan(&count)

	return count > 0
}

// SubmitAnswer registra a resposta do usuário
func (qm *QuizManager) SubmitAnswer(questionID int, userID string, answer string) error {
	_, err := qm.db.Exec(`
		INSERT INTO user_responses (question_id, user_id, answer, answered_at)
		VALUES (?, ?, ?, ?)
	`, questionID, userID, answer, time.Now())

	return err
}

// GetTimeUntilNext calcula o tempo até o próximo quiz
func (qm *QuizManager) GetTimeUntilNext() (*time.Duration, error) {
	nextQuiz, err := qm.GetNextQuiz()
	if err != nil {
		return nil, err
	}

	duration := time.Until(nextQuiz.ScheduledAt)
	return &duration, nil
}

// GetQuizProgress retorna o progresso do usuário
func (qm *QuizManager) GetQuizProgress(userID string) (models.QuizProgress, error) {
	var progress models.QuizProgress

	// Total de quizzes
	qm.db.QueryRow("SELECT COUNT(*) FROM questions WHERE is_active = 1").Scan(&progress.Total)

	// Quizzes respondidos corretamente
	qm.db.QueryRow(`
		SELECT COUNT(*) FROM user_responses ur
		JOIN questions q ON ur.question_id = q.id
		WHERE ur.user_id = ? AND ur.answer = q.correct_answer
	`, userID).Scan(&progress.Correct)

	// Quizzes respondidos (total)
	qm.db.QueryRow(`
		SELECT COUNT(*) FROM user_responses 
		WHERE user_id = ?
	`, userID).Scan(&progress.Answered)

	return progress, nil
}

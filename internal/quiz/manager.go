package quiz

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
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
	now := time.Now().UTC() // FIX: Usar UTC para consistência
	var question models.Question
	var optionsJSON string

	err := qm.db.QueryRow(`
		SELECT id, title, content, options, correct_answer, reward, scheduled_at 
		FROM questions 
		WHERE datetime(scheduled_at) <= datetime(?)
		ORDER BY scheduled_at DESC 
		LIMIT 1
	`, now.Format("2006-01-02T15:04:05Z")).Scan(&question.ID, &question.Title, &question.Content,
		&optionsJSON, &question.CorrectAnswer, &question.Reward, &question.ScheduledAt)

	if err != nil {
		return nil, err
	}

	// Parse JSON options
	if optionsJSON != "" {
		json.Unmarshal([]byte(optionsJSON), &question.Options)
	}

	return &question, nil
}

// GetNextQuiz retorna o próximo quiz agendado
func (qm *QuizManager) GetNextQuiz() (*models.Question, error) {
	now := time.Now().UTC() // FIX: Usar UTC para consistência
	var question models.Question
	var optionsJSON string

	err := qm.db.QueryRow(`
		SELECT id, title, content, options, correct_answer, reward, scheduled_at 
		FROM questions 
		WHERE datetime(scheduled_at) > datetime(?)
		ORDER BY scheduled_at ASC 
		LIMIT 1
	`, now.Format("2006-01-02T15:04:05Z")).Scan(&question.ID, &question.Title, &question.Content,
		&optionsJSON, &question.CorrectAnswer, &question.Reward, &question.ScheduledAt)

	if err != nil {
		return nil, err
	}

	// Parse JSON options
	if optionsJSON != "" {
		json.Unmarshal([]byte(optionsJSON), &question.Options)
	}

	return &question, nil
}

// HasUserAnswered verifica se o usuário já respondeu a pergunta
func (qm *QuizManager) HasUserAnswered(questionID int, userID string) bool {
	var count int
	userIDInt := getUserIDInt(userID)
	qm.db.QueryRow(`
		SELECT COUNT(*) FROM user_responses 
		WHERE question_id = ? AND user_id = ?
	`, questionID, userIDInt).Scan(&count)

	return count > 0
}

// SubmitAnswer registra a resposta do usuário
func (qm *QuizManager) SubmitAnswer(questionID int, userID string, answer string) error {
	// Converter answer string para int (índice da opção)
	answerIndex, err := strconv.Atoi(answer)
	if err != nil {
		return fmt.Errorf("resposta inválida: %v", err)
	}

	// Buscar a pergunta para verificar resposta correta
	var correctAnswer int
	err = qm.db.QueryRow(`
		SELECT correct_answer FROM questions WHERE id = ?
	`, questionID).Scan(&correctAnswer)
	if err != nil {
		return fmt.Errorf("pergunta não encontrada: %v", err)
	}

	// Verificar se a resposta está correta
	isCorrect := answerIndex == correctAnswer

	// Buscar ou criar user_id (por enquanto vamos usar hash do username)
	userIDInt := getUserIDInt(userID)

	// Inserir resposta
	_, err = qm.db.Exec(`
		INSERT INTO user_responses (user_id, question_id, answer, is_correct, answered_at)
		VALUES (?, ?, ?, ?, ?)
	`, userIDInt, questionID, answerIndex, isCorrect, time.Now())

	return err
}

// Helper para converter username em ID numérico consistente
func getUserIDInt(username string) int {
	// Por simplicidade, vamos usar um hash simples do username
	// Em produção, você deveria ter uma tabela de usuários adequada
	hash := 0
	for _, char := range username {
		hash = hash*31 + int(char)
	}
	// Garantir que seja positivo
	if hash < 0 {
		hash = -hash
	}
	return hash % 1000000 // Limitar a 6 dígitos
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
	userIDInt := getUserIDInt(userID)

	// Total de quizzes
	qm.db.QueryRow("SELECT COUNT(*) FROM questions WHERE is_active = 1").Scan(&progress.Total)

	// Quizzes respondidos corretamente
	qm.db.QueryRow(`
		SELECT COUNT(*) FROM user_responses 
		WHERE user_id = ? AND is_correct = 1
	`, userIDInt).Scan(&progress.Correct)

	// Quizzes respondidos (total)
	qm.db.QueryRow(`
		SELECT COUNT(*) FROM user_responses 
		WHERE user_id = ?
	`, userIDInt).Scan(&progress.Answered)

	return progress, nil
}

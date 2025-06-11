package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"namorada-quiz/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	conn *sql.DB
}

// NewDB cria uma nova conexão com o banco
func NewDB(dataSourceName string) (*DB, error) {
	conn, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db := &DB{conn: conn}

	// Executar migrations
	if err := db.migrate(); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}

// migrate executa as migrations do banco
func (db *DB) migrate() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			role TEXT NOT NULL DEFAULT 'visitor',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS questions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			description TEXT,
			options TEXT NOT NULL, -- JSON array
			correct_answer TEXT NOT NULL,
			reward TEXT NOT NULL,
			scheduled_at DATETIME NOT NULL,
			is_active BOOLEAN DEFAULT TRUE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS answers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			question_id INTEGER NOT NULL,
			answer TEXT NOT NULL,
			is_correct BOOLEAN NOT NULL,
			answered_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (question_id) REFERENCES questions(id),
			UNIQUE(user_id, question_id)
		)`,
		`CREATE TABLE IF NOT EXISTS quiz_sessions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			total_questions INTEGER DEFAULT 0,
			correct_answers INTEGER DEFAULT 0,
			completed_at DATETIME,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,
	}

	for _, query := range queries {
		if _, err := db.conn.Exec(query); err != nil {
			return fmt.Errorf("failed to execute migration: %w", err)
		}
	}

	return nil
}

// Close fecha a conexão com o banco
func (db *DB) Close() error {
	return db.conn.Close()
}

// User operations
func (db *DB) CreateUser(user *models.User) error {
	query := `INSERT INTO users (username, password, role) VALUES (?, ?, ?)`
	result, err := db.conn.Exec(query, user.Username, user.Password, user.Role)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get user ID: %w", err)
	}

	user.ID = int(id)
	return nil
}

func (db *DB) GetUserByUsername(username string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, password, role, created_at FROM users WHERE username = ?`

	err := db.conn.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.Password, &user.Role, &user.CreateAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// Question operations
func (db *DB) CreateQuestion(question *models.Question) error {
	optionsJSON, err := json.Marshal(question.Options)
	if err != nil {
		return fmt.Errorf("failed to marshal options: %w", err)
	}

	query := `INSERT INTO questions (title, description, options, correct_answer, reward, scheduled_at) 
			  VALUES (?, ?, ?, ?, ?, ?)`

	result, err := db.conn.Exec(query,
		question.Title,
		question.Description,
		string(optionsJSON),
		question.CorrectAnswer,
		question.Reward,
		question.ScheduledAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create question: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get question ID: %w", err)
	}

	question.ID = int(id)
	return nil
}

func (db *DB) GetAvailableQuestions(userID int) ([]models.QuestionWithStatus, error) {
	now := time.Now()

	query := `
		SELECT q.id, q.title, q.description, q.options, q.correct_answer, q.reward, 
			   q.scheduled_at, q.is_active, q.created_at,
			   CASE WHEN a.id IS NOT NULL THEN 1 ELSE 0 END as is_answered
		FROM questions q
		LEFT JOIN answers a ON q.id = a.question_id AND a.user_id = ?
		WHERE q.is_active = 1
		ORDER BY q.scheduled_at ASC
	`

	rows, err := db.conn.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get questions: %w", err)
	}
	defer rows.Close()

	var questions []models.QuestionWithStatus
	for rows.Next() {
		var q models.QuestionWithStatus
		var optionsJSON string
		var isAnswered int

		err := rows.Scan(
			&q.ID, &q.Title, &q.Description, &optionsJSON, &q.CorrectAnswer,
			&q.Reward, &q.ScheduledAt, &q.IsActive, &q.CreatedAt, &isAnswered,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan question: %w", err)
		}

		if err := json.Unmarshal([]byte(optionsJSON), &q.Options); err != nil {
			return nil, fmt.Errorf("failed to unmarshal options: %w", err)
		}

		q.IsAvailable = q.ScheduledAt.Before(now) || q.ScheduledAt.Equal(now)
		q.IsAnswered = isAnswered == 1

		questions = append(questions, q)
	}

	return questions, rows.Err()
}

func (db *DB) GetAllQuestions() ([]models.Question, error) {
	query := `SELECT id, title, description, options, correct_answer, reward, 
			         scheduled_at, is_active, created_at FROM questions ORDER BY scheduled_at ASC`

	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get questions: %w", err)
	}
	defer rows.Close()

	var questions []models.Question
	for rows.Next() {
		var q models.Question
		var optionsJSON string

		err := rows.Scan(
			&q.ID, &q.Title, &q.Description, &optionsJSON, &q.CorrectAnswer,
			&q.Reward, &q.ScheduledAt, &q.IsActive, &q.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan question: %w", err)
		}

		if err := json.Unmarshal([]byte(optionsJSON), &q.Options); err != nil {
			return nil, fmt.Errorf("failed to unmarshal options: %w", err)
		}

		questions = append(questions, q)
	}

	return questions, rows.Err()
}

// Answer operations
func (db *DB) CreateAnswer(answer *models.Answer) error {
	query := `INSERT INTO answers (user_id, question_id, answer, is_correct) VALUES (?, ?, ?, ?)`
	result, err := db.conn.Exec(query, answer.UserID, answer.QuestionID, answer.Answer, answer.IsCorrect)
	if err != nil {
		return fmt.Errorf("failed to create answer: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get answer ID: %w", err)
	}

	answer.ID = int(id)
	return nil
}

func (db *DB) GetQuestionByID(id int) (*models.Question, error) {
	question := &models.Question{}
	var optionsJSON string

	query := `SELECT id, title, description, options, correct_answer, reward, 
			         scheduled_at, is_active, created_at FROM questions WHERE id = ?`

	err := db.conn.QueryRow(query, id).Scan(
		&question.ID, &question.Title, &question.Description, &optionsJSON,
		&question.CorrectAnswer, &question.Reward, &question.ScheduledAt,
		&question.IsActive, &question.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get question: %w", err)
	}

	if err := json.Unmarshal([]byte(optionsJSON), &question.Options); err != nil {
		return nil, fmt.Errorf("failed to unmarshal options: %w", err)
	}

	return question, nil
}

func (db *DB) GetAnswerByUserAndQuestion(userID, questionID int) (*models.Answer, error) {
	answer := &models.Answer{}

	query := `SELECT id, user_id, question_id, answer, is_correct, answered_at 
			  FROM answers WHERE user_id = ? AND question_id = ?`

	err := db.conn.QueryRow(query, userID, questionID).Scan(
		&answer.ID, &answer.UserID, &answer.QuestionID, &answer.Answer,
		&answer.IsCorrect, &answer.AnsweredAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get answer: %w", err)
	}

	return answer, nil
}

func (db *DB) GetUserStats(userID int) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total de perguntas
	var totalQuestions int
	err := db.conn.QueryRow("SELECT COUNT(*) FROM questions WHERE is_active = TRUE").Scan(&totalQuestions)
	if err != nil {
		return nil, fmt.Errorf("failed to get total questions: %w", err)
	}

	// Perguntas respondidas
	var answeredQuestions int
	err = db.conn.QueryRow("SELECT COUNT(*) FROM answers WHERE user_id = ?", userID).Scan(&answeredQuestions)
	if err != nil {
		return nil, fmt.Errorf("failed to get answered questions: %w", err)
	}

	// Respostas corretas
	var correctAnswers int
	err = db.conn.QueryRow("SELECT COUNT(*) FROM answers WHERE user_id = ? AND is_correct = TRUE", userID).Scan(&correctAnswers)
	if err != nil {
		return nil, fmt.Errorf("failed to get correct answers: %w", err)
	}

	// Perguntas disponíveis agora
	var availableQuestions int
	err = db.conn.QueryRow(`
		SELECT COUNT(*) FROM questions q 
		WHERE q.is_active = TRUE 
		AND q.scheduled_at <= datetime('now') 
		AND NOT EXISTS (
			SELECT 1 FROM answers a 
			WHERE a.question_id = q.id AND a.user_id = ?
		)
	`, userID).Scan(&availableQuestions)
	if err != nil {
		return nil, fmt.Errorf("failed to get available questions: %w", err)
	}

	stats["total_questions"] = totalQuestions
	stats["answered_questions"] = answeredQuestions
	stats["correct_answers"] = correctAnswers
	stats["available_questions"] = availableQuestions
	stats["accuracy"] = 0.0

	if answeredQuestions > 0 {
		stats["accuracy"] = float64(correctAnswers) / float64(answeredQuestions) * 100
	}

	return stats, nil
}

func (db *DB) GetAllUsers() ([]models.User, error) {
	var users []models.User

	query := `SELECT id, username, role, created_at FROM users ORDER BY created_at DESC`
	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Username, &user.Role, &user.CreateAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	return users, nil
}

func (db *DB) GetAdminStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total de usuários
	var totalUsers int
	err := db.conn.QueryRow("SELECT COUNT(*) FROM users WHERE role = 'visitor'").Scan(&totalUsers)
	if err != nil {
		return nil, fmt.Errorf("failed to get total users: %w", err)
	}

	// Total de perguntas
	var totalQuestions int
	err = db.conn.QueryRow("SELECT COUNT(*) FROM questions").Scan(&totalQuestions)
	if err != nil {
		return nil, fmt.Errorf("failed to get total questions: %w", err)
	}

	// Perguntas ativas
	var activeQuestions int
	err = db.conn.QueryRow("SELECT COUNT(*) FROM questions WHERE is_active = TRUE").Scan(&activeQuestions)
	if err != nil {
		return nil, fmt.Errorf("failed to get active questions: %w", err)
	}

	// Total de respostas
	var totalAnswers int
	err = db.conn.QueryRow("SELECT COUNT(*) FROM answers").Scan(&totalAnswers)
	if err != nil {
		return nil, fmt.Errorf("failed to get total answers: %w", err)
	}

	// Respostas corretas
	var correctAnswers int
	err = db.conn.QueryRow("SELECT COUNT(*) FROM answers WHERE is_correct = TRUE").Scan(&correctAnswers)
	if err != nil {
		return nil, fmt.Errorf("failed to get correct answers: %w", err)
	}

	stats["total_users"] = totalUsers
	stats["total_questions"] = totalQuestions
	stats["active_questions"] = activeQuestions
	stats["total_answers"] = totalAnswers
	stats["correct_answers"] = correctAnswers
	stats["accuracy"] = 0.0

	if totalAnswers > 0 {
		stats["accuracy"] = float64(correctAnswers) / float64(totalAnswers) * 100
	}

	return stats, nil
}

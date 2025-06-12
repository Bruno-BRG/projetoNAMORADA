package database

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func Initialize() (*sql.DB, error) {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./quiz.db"
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if err := createTables(db); err != nil {
		return nil, err
	}

	if err := seedDefaultData(db); err != nil {
		return nil, err
	}

	return db, nil
}

func createTables(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			password_hash TEXT,
			is_admin BOOLEAN DEFAULT FALSE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS questions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			options TEXT NOT NULL, -- JSON array
			correct_answer INTEGER NOT NULL,
			reward TEXT NOT NULL,
			scheduled_at DATETIME NOT NULL,
			is_active BOOLEAN DEFAULT TRUE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS user_responses (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			question_id INTEGER NOT NULL,
			answer INTEGER NOT NULL,
			is_correct BOOLEAN NOT NULL,
			answered_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (question_id) REFERENCES questions(id),
			UNIQUE(user_id, question_id)
		)`,
		`CREATE INDEX IF NOT EXISTS idx_questions_scheduled ON questions(scheduled_at)`,
		`CREATE INDEX IF NOT EXISTS idx_responses_user ON user_responses(user_id)`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return err
		}
	}

	return nil
}

func seedDefaultData(db *sql.DB) error {
	// Criar usuário admin padrão se não existir
	_, err := db.Exec(`
		INSERT OR IGNORE INTO users (username, is_admin) 
		VALUES ('admin', TRUE)
	`)
	if err != nil {
		return err
	}
	// Criar usuário visitante padrão
	_, err = db.Exec(`
		INSERT OR IGNORE INTO users (username, is_admin) 
		VALUES ('momo', FALSE)
	`)

	return err
}

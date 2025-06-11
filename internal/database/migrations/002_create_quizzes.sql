-- Criação da tabela de quizzes
CREATE TABLE quizzes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    description TEXT,
    scheduled_at DATETIME NOT NULL,
    is_active BOOLEAN DEFAULT 0,
    reward TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Criação da tabela de perguntas
CREATE TABLE questions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    quiz_id INTEGER NOT NULL,
    text TEXT NOT NULL,
    options TEXT NOT NULL, -- JSON array
    correct_option INTEGER NOT NULL,
    question_order INTEGER NOT NULL,
    FOREIGN KEY (quiz_id) REFERENCES quizzes(id) ON DELETE CASCADE
);

-- Índices para performance
CREATE INDEX idx_questions_quiz_id ON questions(quiz_id);
CREATE INDEX idx_questions_order ON questions(quiz_id, question_order);

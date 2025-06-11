-- Criação da tabela de rate limiting
CREATE TABLE rate_limits (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    ip_address TEXT NOT NULL,
    quiz_id INTEGER NOT NULL,
    attempts INTEGER DEFAULT 1,
    last_attempt DATETIME DEFAULT CURRENT_TIMESTAMP,
    reset_at DATETIME NOT NULL,
    UNIQUE(ip_address, quiz_id)
);

-- Criação da tabela de sessões
CREATE TABLE sessions (
    id TEXT PRIMARY KEY,
    user_id INTEGER NOT NULL,
    expires_at DATETIME NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    ip_address TEXT,
    user_agent TEXT,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Índices para performance
CREATE INDEX idx_rate_limits_ip ON rate_limits(ip_address);
CREATE INDEX idx_rate_limits_reset ON rate_limits(reset_at);
CREATE INDEX idx_sessions_user ON sessions(user_id);
CREATE INDEX idx_sessions_expires ON sessions(expires_at);

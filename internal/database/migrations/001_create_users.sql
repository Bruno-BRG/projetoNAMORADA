-- Criação da tabela de usuários
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL CHECK (role IN ('admin', 'user')),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    last_login_at DATETIME
);

-- Inserir usuário admin padrão (senha: admin123)
-- Hash bcrypt de "admin123"
INSERT INTO users (username, password_hash, role) VALUES 
('admin', '$2a$10$N9qo8uLOickgx2ZMRZoMye.Hn5.Mw2FhA2BhJePnPTJzFGpY8b7wS', 'admin');

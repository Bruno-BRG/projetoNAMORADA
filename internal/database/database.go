package database

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"sort"
	"strings"

	_ "modernc.org/sqlite"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

type DB struct {
	*sql.DB
}

// NewConnection cria uma nova conexão com o banco SQLite
func NewConnection(dbPath string) (*DB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir banco: %w", err)
	}

	// Configurações importantes para SQLite
	if _, err := db.Exec(`
		PRAGMA foreign_keys = ON;
		PRAGMA journal_mode = WAL;
		PRAGMA synchronous = NORMAL;
		PRAGMA cache_size = 1000;
		PRAGMA temp_store = memory;
	`); err != nil {
		return nil, fmt.Errorf("erro ao configurar SQLite: %w", err)
	}

	dbConn := &DB{db}

	// Executar migrações
	if err := dbConn.runMigrations(); err != nil {
		return nil, fmt.Errorf("erro nas migrações: %w", err)
	}

	return dbConn, nil
}

func (db *DB) runMigrations() error {
	// Criar tabela de migrações se não existir
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			filename TEXT PRIMARY KEY,
			executed_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`); err != nil {
		return fmt.Errorf("erro ao criar tabela de migrações: %w", err)
	}

	// Listar arquivos de migração
	entries, err := fs.ReadDir(migrationFiles, "migrations")
	if err != nil {
		return fmt.Errorf("erro ao ler migrações: %w", err)
	}

	var sqlFiles []string
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".sql") {
			sqlFiles = append(sqlFiles, entry.Name())
		}
	}
	sort.Strings(sqlFiles)

	// Executar migrações pendentes
	for _, filename := range sqlFiles {
		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM migrations WHERE filename = ?)", filename).Scan(&exists)
		if err != nil {
			return fmt.Errorf("erro ao verificar migração %s: %w", filename, err)
		}

		if !exists {
			log.Printf("Executando migração: %s", filename)

			content, err := migrationFiles.ReadFile(filepath.Join("migrations", filename))
			if err != nil {
				return fmt.Errorf("erro ao ler migração %s: %w", filename, err)
			}

			if _, err := db.Exec(string(content)); err != nil {
				return fmt.Errorf("erro ao executar migração %s: %w", filename, err)
			}

			if _, err := db.Exec("INSERT INTO migrations (filename) VALUES (?)", filename); err != nil {
				return fmt.Errorf("erro ao registrar migração %s: %w", filename, err)
			}
		}
	}

	return nil
}

// Close fecha a conexão com o banco
func (db *DB) Close() error {
	return db.DB.Close()
}

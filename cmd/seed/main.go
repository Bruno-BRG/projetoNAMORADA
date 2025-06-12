package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Carregar .env
	_ = godotenv.Load()

	// Conectar ao banco
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./quiz.db"
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal("Erro ao conectar ao banco:", err)
	}
	defer db.Close()

	// Perguntas de exemplo
	questions := []struct {
		Title         string
		Content       string
		Options       []string
		CorrectAnswer int
		Reward        string
		ScheduledAt   time.Time
	}{
		{
			Title:         "Nosso Primeiro Encontro",
			Content:       "Onde foi nosso primeiro encontro oficial?",
			Options:       []string{"Cinema", "Restaurante", "Parque", "Casa dela"},
			CorrectAnswer: 1,
			Reward:        "Um beijo apaixonado 💋",
			ScheduledAt:   time.Now().Add(1 * time.Minute),
		},
		{
			Title:         "Nossa Música",
			Content:       "Qual é a música que mais representa nosso relacionamento?",
			Options:       []string{"Perfect - Ed Sheeran", "All of Me - John Legend", "Thinking Out Loud - Ed Sheeran", "A Thousand Years - Christina Perri"},
			CorrectAnswer: 0,
			Reward:        "Uma massagem relaxante 💆‍♀️",
			ScheduledAt:   time.Now().Add(5 * time.Minute),
		},
		{
			Title:         "Comida Favorita",
			Content:       "Qual é minha comida favorita que você prepara?",
			Options:       []string{"Lasanha", "Feijoada", "Strogonoff", "Pizza"},
			CorrectAnswer: 2,
			Reward:        "Jantar romântico em casa 🍽️❤️",
			ScheduledAt:   time.Now().Add(10 * time.Minute),
		},
	}

	// Inserir perguntas
	for _, q := range questions {
		optionsJSON, _ := json.Marshal(q.Options)
		
		_, err := db.Exec(`
			INSERT INTO questions (title, content, options, correct_answer, reward, scheduled_at, is_active)
			VALUES (?, ?, ?, ?, ?, ?, 1)
		`, q.Title, q.Content, string(optionsJSON), q.CorrectAnswer, q.Reward, q.ScheduledAt)
		
		if err != nil {
			log.Printf("Erro ao inserir pergunta '%s': %v", q.Title, err)
		} else {
			log.Printf("✅ Pergunta criada: %s (liberada às %s)", q.Title, q.ScheduledAt.Format("15:04"))
		}
	}

	log.Println("🎉 Perguntas de exemplo criadas com sucesso!")
}

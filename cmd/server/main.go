package main

import (
	"log"
	"os"

	"namorada-quiz/internal/database"
	"namorada-quiz/internal/handlers"
	"namorada-quiz/internal/middleware"
	"namorada-quiz/internal/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Configurar ambiente
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.DebugMode)
	}

	// Configurar banco de dados
	dbPath := getEnv("DATABASE_PATH", "./quiz.db")
	db, err := database.NewDB(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Criar usuário admin padrão se não existir
	if err := createDefaultAdmin(db); err != nil {
		log.Printf("Warning: Failed to create default admin: %v", err)
	}

	// Configurar handlers
	handler := handlers.NewHandler(db)

	// Configurar router
	r := gin.Default()

	// Carregar templates
	r.LoadHTMLGlob("templates/*")

	// Servir arquivos estáticos
	r.Static("/static", "./static")

	// Middleware global
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.SecurityHeaders())

	// Rotas públicas
	r.GET("/", handler.Home)
	r.GET("/login", handler.LoginPage)
	r.POST("/api/login", handler.Login)
	// Rotas protegidas
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/api/logout", handler.Logout)

		// Rotas para visitantes (namorada)
		protected.GET("/dashboard", middleware.RequireRole("visitor"), handler.Dashboard)
		protected.GET("/api/questions/available", middleware.RequireRole("visitor"), handler.GetAvailableQuestions)
		protected.POST("/api/questions/answer", middleware.RequireRole("visitor"), handler.AnswerQuestion)
		protected.GET("/api/stats", middleware.RequireRole("visitor"), handler.GetUserStats)

		// Rotas HTMX para templates parciais
		protected.GET("/api/stats/render", middleware.RequireRole("visitor"), handler.RenderStats)
		protected.GET("/api/questions/render", middleware.RequireRole("visitor"), handler.RenderQuestions)
		protected.GET("/api/questions/:id/form", middleware.RequireRole("visitor"), handler.RenderQuestionForm)

		// Rotas para admin
		admin := protected.Group("/admin")
		admin.Use(middleware.RequireRole("admin"))
		{
			admin.GET("/", handler.AdminDashboard)
			admin.GET("/questions", handler.GetAllQuestions)
			admin.POST("/questions", handler.CreateQuestion)
			admin.PUT("/questions/:id", handler.UpdateQuestion)
			admin.DELETE("/questions/:id", handler.DeleteQuestion)
			admin.GET("/users", handler.GetAllUsers)
			admin.GET("/stats", handler.GetAdminStats)
		}
	}

	// Iniciar servidor
	port := getEnv("PORT", "8080")
	log.Printf("Server starting on port %s", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func createDefaultAdmin(db *database.DB) error {
	// Verificar se já existe um admin
	user, err := db.GetUserByUsername("admin")
	if err != nil {
		return err
	}

	if user != nil {
		return nil // Admin já existe
	}

	// Criar senha hash
	password := getEnv("ADMIN_PASSWORD", "admin123")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Criar usuário admin
	admin := &models.User{
		Username: "admin",
		Password: string(hashedPassword),
		Role:     "admin",
	}

	return db.CreateUser(admin)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

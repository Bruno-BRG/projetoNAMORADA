package main

import (
	"log"
	"net/http"
	"os"
	"valentine-quiz/internal/database"
	"valentine-quiz/internal/handlers"
	"valentine-quiz/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Carregar variáveis de ambiente
	_ = godotenv.Load()

	// Inicializar banco de dados
	db, err := database.Initialize()
	if err != nil {
		log.Fatal("Erro ao inicializar banco:", err)
	}
	defer db.Close()

	// Setup do Gin
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// Middleware
	r.Use(middleware.CORS())
	r.Use(middleware.SecurityHeaders())
	// Templates e arquivos estáticos
	r.LoadHTMLGlob("web/templates/*")
	r.Static("/static", "./web/static")

	// Handlers
	h := handlers.New(db)
	// Rotas públicas
	public := r.Group("/")
	{
		public.GET("/", h.Home)
		public.GET("/login", h.LoginPage)
		public.POST("/login", h.Login)
		public.GET("/logout", h.Logout)
		public.GET("/debug", h.DebugQuizStatus) // Debug sem auth para testar
	}
	// Rotas protegidas para visitante
	visitor := r.Group("/quiz")
	visitor.Use(middleware.RequireVisitorAuth())
	{
		visitor.GET("/", h.QuizHome)
		visitor.GET("/status", h.QuizStatus)
		visitor.GET("/countdown", h.Countdown)
		visitor.GET("/current", h.CurrentQuiz)
		visitor.POST("/answer", h.SubmitAnswer)
		visitor.GET("/progress", h.Progress)
	}
	// Rotas de admin
	admin := r.Group("/admin")
	admin.Use(middleware.RequireAdminAuth())
	{
		admin.GET("/", h.AdminDashboard)
		admin.GET("/questions", h.ListQuestions)
		admin.GET("/questions/new", h.NewQuestionForm)
		admin.POST("/questions", h.CreateQuestion)
		admin.GET("/questions/:id/edit", h.EditQuestionForm)
		admin.PUT("/questions/:id", h.UpdateQuestion)
		admin.POST("/questions/:id", h.UpdateQuestion) // Permite POST também para HTML forms
		admin.DELETE("/questions/:id", h.DeleteQuestion)
		admin.GET("/responses", h.ViewResponses)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Servidor rodando na porta %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

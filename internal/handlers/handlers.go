package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"valentine-quiz/internal/auth"
	"valentine-quiz/internal/models"
	"valentine-quiz/internal/quiz"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	db          *sql.DB
	quizManager *quiz.QuizManager
}

func New(db *sql.DB) *Handler {
	return &Handler{
		db:          db,
		quizManager: quiz.NewQuizManager(db),
	}
}

// getUserID extrai o ID do usuário do contexto (simplificado para o projeto)
func (h *Handler) getUserID(c *gin.Context) string {
	token, err := c.Cookie("visitor_session")
	if err != nil || token == "" {
		return ""
	}

	// Validar token (simplificado)
	claims, err := auth.ValidateToken(token)
	if err != nil {
		return ""
	}

	return claims.Username // Usar username como ID único
}

// Página inicial
func (h *Handler) Home(c *gin.Context) {
	c.HTML(http.StatusOK, "home.html", gin.H{
		"title": "Quiz do Dia dos Namorados 💕",
	})
}

// Página de login
func (h *Handler) LoginPage(c *gin.Context) {
	isAdmin := c.Query("admin") == "1"
	c.HTML(http.StatusOK, "login.html", gin.H{
		"title":   "Login",
		"isAdmin": isAdmin,
	})
}

// Processo de login melhorado
func (h *Handler) Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	isAdmin := c.PostForm("admin") == "1"

	if auth.CheckCredentials(username, password, isAdmin) {
		token, err := auth.GenerateToken(username, isAdmin)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "login.html", gin.H{
				"error":   "Erro interno do servidor",
				"isAdmin": isAdmin,
			})
			return
		}

		if isAdmin {
			c.SetCookie("admin_session", token, 86400, "/", "", false, true)
			c.Redirect(http.StatusFound, "/admin")
		} else {
			c.SetCookie("visitor_session", token, 86400, "/", "", false, true)
			c.Redirect(http.StatusFound, "/quiz")
		}
		return
	}

	c.HTML(http.StatusUnauthorized, "login.html", gin.H{
		"error":   "Credenciais inválidas",
		"isAdmin": isAdmin,
	})
}

// Logout
func (h *Handler) Logout(c *gin.Context) {
	c.SetCookie("admin_session", "", -1, "/", "", false, true)
	c.SetCookie("visitor_session", "", -1, "/", "", false, true)
	c.Redirect(http.StatusFound, "/")
}

// Quiz Dashboard (nova abordagem HTMX)
func (h *Handler) QuizHome(c *gin.Context) {
	userID := h.getUserID(c)
	if userID == "" {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	status := h.getQuizStatus(userID)
	c.HTML(http.StatusOK, "quiz_dashboard.html", status)
}

// Status do Quiz (endpoint HTMX)
func (h *Handler) QuizStatus(c *gin.Context) {
	userID := h.getUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Não autorizado"})
		return
	}

	status := h.getQuizStatus(userID)
	c.HTML(http.StatusOK, "quiz_content.html", status)
}

// Helper para obter status do quiz
func (h *Handler) getQuizStatus(userID string) gin.H {
	// Quiz atual disponível
	currentQuiz, err := h.quizManager.GetAvailableQuiz()
	hasCurrent := err == nil

	// Verificar se já respondeu
	alreadyAnswered := false
	if hasCurrent {
		alreadyAnswered = h.quizManager.HasUserAnswered(currentQuiz.ID, userID)
	}

	// Próximo quiz
	nextQuiz, err := h.quizManager.GetNextQuiz()
	hasNext := err == nil

	// Tempo até o próximo
	timeUntilNext := ""
	if hasNext {
		duration := time.Until(nextQuiz.ScheduledAt)
		timeUntilNext = formatDuration(duration)
	}

	// Progresso
	progress, _ := h.quizManager.GetQuizProgress(userID)

	return gin.H{
		"HasCurrent":      hasCurrent,
		"Current":         currentQuiz,
		"AlreadyAnswered": alreadyAnswered,
		"HasNext":         hasNext,
		"Next":            nextQuiz,
		"TimeUntilNext":   timeUntilNext,
		"Progress":        progress,
	}
}

// Countdown endpoint para HTMX
func (h *Handler) Countdown(c *gin.Context) {
	nextQuiz, err := h.quizManager.GetNextQuiz()
	if err != nil {
		c.String(http.StatusOK, "Nenhum quiz agendado")
		return
	}

	duration := time.Until(nextQuiz.ScheduledAt)
	if duration <= 0 {
		c.String(http.StatusOK, "Disponível agora!")
		return
	}

	c.String(http.StatusOK, formatDuration(duration))
}

// Submeter resposta (HTMX endpoint)
func (h *Handler) SubmitAnswer(c *gin.Context) {
	userID := h.getUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Não autorizado"})
		return
	}

	questionID, _ := strconv.Atoi(c.PostForm("question_id"))
	answer := c.PostForm("answer")

	// Verificar se já respondeu
	if h.quizManager.HasUserAnswered(questionID, userID) {
		status := h.getQuizStatus(userID)
		c.HTML(http.StatusOK, "quiz_content.html", status)
		return
	}

	// Registrar resposta
	err := h.quizManager.SubmitAnswer(questionID, userID, answer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Retornar status atualizado
	status := h.getQuizStatus(userID)
	c.HTML(http.StatusOK, "quiz_content.html", status)
}

// Helper para formatar duração
func formatDuration(d time.Duration) string {
	if d <= 0 {
		return "Disponível agora!"
	}

	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if hours > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
	}
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

// === ADMIN HANDLERS ===

// Dashboard admin
func (h *Handler) AdminDashboard(c *gin.Context) {
	var totalQuestions, totalResponses, correctResponses int

	h.db.QueryRow("SELECT COUNT(*) FROM questions").Scan(&totalQuestions)
	h.db.QueryRow("SELECT COUNT(*) FROM user_responses").Scan(&totalResponses)
	h.db.QueryRow("SELECT COUNT(*) FROM user_responses WHERE is_correct = 1").Scan(&correctResponses)

	c.HTML(http.StatusOK, "admin_dashboard.html", gin.H{
		"title":            "Dashboard Admin",
		"totalQuestions":   totalQuestions,
		"totalResponses":   totalResponses,
		"correctResponses": correctResponses,
	})
}

// Listar perguntas
func (h *Handler) ListQuestions(c *gin.Context) {
	rows, err := h.db.Query(`
		SELECT id, title, content, options, reward, scheduled_at, is_active 
		FROM questions 
		ORDER BY scheduled_at ASC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var questions []models.Question
	for rows.Next() {
		var q models.Question
		var optionsJSON string
		err = rows.Scan(&q.ID, &q.Title, &q.Content, &optionsJSON,
			&q.Reward, &q.ScheduledAt, &q.IsActive)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Parse options JSON
		if optionsJSON != "" {
			json.Unmarshal([]byte(optionsJSON), &q.Options)
		}

		questions = append(questions, q)
	}

	c.HTML(http.StatusOK, "admin_questions.html", gin.H{
		"title":     "Gerenciar Perguntas",
		"questions": questions,
	})
}

// Form nova pergunta
func (h *Handler) NewQuestionForm(c *gin.Context) {
	c.HTML(http.StatusOK, "admin_question_form.html", gin.H{
		"title": "Nova Pergunta",
	})
}

// Criar pergunta
func (h *Handler) CreateQuestion(c *gin.Context) {
	title := c.PostForm("title")
	content := c.PostForm("content")
	option1 := c.PostForm("option1")
	option2 := c.PostForm("option2")
	option3 := c.PostForm("option3")
	option4 := c.PostForm("option4")
	correctAnswer, _ := strconv.Atoi(c.PostForm("correct_answer"))
	reward := c.PostForm("reward")
	scheduledAt := c.PostForm("scheduled_at")

	options := []string{option1, option2, option3, option4}
	optionsJSON, _ := json.Marshal(options)

	_, err := h.db.Exec(`
		INSERT INTO questions (title, content, options, correct_answer, reward, scheduled_at, is_active)
		VALUES (?, ?, ?, ?, ?, ?, 1)
	`, title, content, string(optionsJSON), correctAnswer, reward, scheduledAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "/admin/questions")
}

// Placeholders para funções restantes
func (h *Handler) CurrentQuiz(c *gin.Context) { h.QuizStatus(c) }
func (h *Handler) Progress(c *gin.Context)    { /* implementar se necessário */ }

// Editar pergunta - formulário
func (h *Handler) EditQuestionForm(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var question models.Question
	var optionsJSON string

	err = h.db.QueryRow(`
		SELECT id, title, content, options, correct_answer, reward, scheduled_at, is_active
		FROM questions WHERE id = ?
	`, id).Scan(&question.ID, &question.Title, &question.Content,
		&optionsJSON, &question.CorrectAnswer, &question.Reward,
		&question.ScheduledAt, &question.IsActive)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Pergunta não encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Parse options JSON
	if optionsJSON != "" {
		json.Unmarshal([]byte(optionsJSON), &question.Options)
	}

	c.HTML(http.StatusOK, "admin_question_edit.html", gin.H{
		"title":    "Editar Pergunta",
		"question": question,
	})
}

// Atualizar pergunta
func (h *Handler) UpdateQuestion(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	title := c.PostForm("title")
	content := c.PostForm("content")
	option1 := c.PostForm("option1")
	option2 := c.PostForm("option2")
	option3 := c.PostForm("option3")
	option4 := c.PostForm("option4")
	correctAnswer, _ := strconv.Atoi(c.PostForm("correct_answer"))
	reward := c.PostForm("reward")
	scheduledAt := c.PostForm("scheduled_at")

	options := []string{option1, option2, option3, option4}
	optionsJSON, _ := json.Marshal(options)

	result, err := h.db.Exec(`
		UPDATE questions 
		SET title = ?, content = ?, options = ?, correct_answer = ?, 
			reward = ?, scheduled_at = ?
		WHERE id = ?
	`, title, content, string(optionsJSON), correctAnswer, reward, scheduledAt, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pergunta não encontrada"})
		return
	}

	c.Redirect(http.StatusFound, "/admin/questions")
}

// Deletar pergunta (HTMX)
func (h *Handler) DeleteQuestion(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Verificar se a pergunta existe
	var exists bool
	err = h.db.QueryRow("SELECT EXISTS(SELECT 1 FROM questions WHERE id = ?)", id).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao verificar pergunta"})
		return
	}

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pergunta não encontrada"})
		return
	}

	// Deletar respostas relacionadas primeiro (integridade referencial)
	_, err = h.db.Exec("DELETE FROM user_responses WHERE question_id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar respostas"})
		return
	}

	// Deletar a pergunta
	result, err := h.db.Exec("DELETE FROM questions WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar pergunta"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pergunta não encontrada"})
		return
	}

	// Para HTMX, retornar vazio remove o elemento
	c.Status(http.StatusOK)
}

// Ver respostas
func (h *Handler) ViewResponses(c *gin.Context) {
	rows, err := h.db.Query(`
		SELECT 
			q.title,
			q.content,
			ur.answer,
			ur.is_correct,
			ur.answered_at,
			q.reward
		FROM user_responses ur
		JOIN questions q ON ur.question_id = q.id
		ORDER BY ur.answered_at DESC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	type ResponseView struct {
		QuestionTitle   string    `json:"question_title"`
		QuestionContent string    `json:"question_content"`
		Answer          int       `json:"answer"`
		IsCorrect       bool      `json:"is_correct"`
		AnsweredAt      time.Time `json:"answered_at"`
		Reward          string    `json:"reward"`
	}

	var responses []ResponseView
	for rows.Next() {
		var r ResponseView
		rows.Scan(&r.QuestionTitle, &r.QuestionContent, &r.Answer,
			&r.IsCorrect, &r.AnsweredAt, &r.Reward)
		responses = append(responses, r)
	}

	c.HTML(http.StatusOK, "admin_responses.html", gin.H{
		"title":     "Respostas dos Quiz",
		"responses": responses,
	})
}

// Debug endpoint para verificar status
func (h *Handler) DebugQuizStatus(c *gin.Context) {
	userID := h.getUserID(c)

	// Quiz atual disponível
	currentQuiz, err := h.quizManager.GetAvailableQuiz()
	hasCurrent := err == nil

	// Próximo quiz
	nextQuiz, errNext := h.quizManager.GetNextQuiz()
	hasNext := errNext == nil

	// Progresso
	progress, _ := h.quizManager.GetQuizProgress(userID)

	debugInfo := gin.H{
		"UserID":       userID,
		"HasCurrent":   hasCurrent,
		"CurrentError": err,
		"Current":      currentQuiz,
		"HasNext":      hasNext,
		"NextError":    errNext,
		"Next":         nextQuiz,
		"Progress":     progress,
		"Now":          time.Now().Format("2006-01-02 15:04:05"),
	}

	c.JSON(http.StatusOK, debugInfo)
}

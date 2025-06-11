package handlers

import (
	"fmt"
	"net/http"
	"time"

	"namorada-quiz/internal/database"
	"namorada-quiz/internal/middleware"
	"namorada-quiz/internal/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	db *database.DB
}

func NewHandler(db *database.DB) *Handler {
	return &Handler{db: db}
}

// Auth handlers
func (h *Handler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	user, err := h.db.GetUserByUsername(req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Verificar senha
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Gerar JWT
	token, err := middleware.GenerateJWT(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Definir cookie httpOnly
	c.SetCookie("auth_token", token, 24*3600, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"role":     user.Role,
		},
		"token": token,
	})
}

func (h *Handler) Logout(c *gin.Context) {
	c.SetCookie("auth_token", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// Question handlers
func (h *Handler) CreateQuestion(c *gin.Context) {
	var req models.CreateQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}

	// Parse scheduled time
	scheduledAt, err := time.Parse(time.RFC3339, req.ScheduledAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid scheduled_at format. Use ISO 8601 format"})
		return
	}

	question := &models.Question{
		Title:         req.Title,
		Description:   req.Description,
		Options:       req.Options,
		CorrectAnswer: req.CorrectAnswer,
		Reward:        req.Reward,
		ScheduledAt:   scheduledAt,
		IsActive:      true,
	}

	if err := h.db.CreateQuestion(question); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create question"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Question created successfully",
		"question": question,
	})
}

func (h *Handler) GetQuestions(c *gin.Context) {
	// Para admin: retorna todas as perguntas
	// Para visitor: retorna apenas perguntas disponíveis
	role, _ := c.Get("role")
	userID, _ := c.Get("user_id")

	if role == "admin" {
		questions, err := h.db.GetAllQuestions()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get questions"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"questions": questions})
	} else {
		questions, err := h.db.GetAvailableQuestions(userID.(int))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get questions"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"questions": questions})
	}
}

// Answer handlers
func (h *Handler) SubmitAnswer(c *gin.Context) {
	var req models.AnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	userID, _ := c.Get("user_id")

	// Verificar se a pergunta existe e está disponível
	question, err := h.db.GetQuestionByID(req.QuestionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if question == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}

	// Verificar se a pergunta está disponível (horário já passou)
	if time.Now().Before(question.ScheduledAt) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Question not available yet"})
		return
	}

	// Verificar se a resposta está correta
	isCorrect := req.Answer == question.CorrectAnswer

	answer := &models.Answer{
		UserID:     userID.(int),
		QuestionID: req.QuestionID,
		Answer:     req.Answer,
		IsCorrect:  isCorrect,
	}

	if err := h.db.CreateAnswer(answer); err != nil {
		// Se for erro de constraint (já respondeu), retornar erro específico
		c.JSON(http.StatusConflict, gin.H{"error": "Question already answered"})
		return
	}

	response := gin.H{
		"message":    "Answer submitted successfully",
		"is_correct": isCorrect,
	}

	// Se a resposta estiver correta, incluir a recompensa
	if isCorrect {
		response["reward"] = question.Reward
		response["congratulations"] = "🎉 Resposta correta! Você ganhou uma recompensa!"
	} else {
		response["feedback"] = "😔 Resposta incorreta, mas não desista!"
	}

	c.JSON(http.StatusOK, response)
}

// Dashboard handlers
func (h *Handler) GetDashboard(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	if role == "admin" {
		// Dashboard do admin
		questions, err := h.db.GetAllQuestions()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get questions"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"role":      "admin",
			"questions": questions,
			"stats": gin.H{
				"total_questions": len(questions),
				"active_questions": func() int {
					count := 0
					for _, q := range questions {
						if q.IsActive {
							count++
						}
					}
					return count
				}(),
			},
		})
	} else {
		// Dashboard do visitor
		questions, err := h.db.GetAvailableQuestions(userID.(int))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get questions"})
			return
		}

		stats := gin.H{
			"total_questions":     len(questions),
			"available_questions": 0,
			"answered_questions":  0,
			"correct_answers":     0,
		}

		for _, q := range questions {
			if q.IsAvailable {
				stats["available_questions"] = stats["available_questions"].(int) + 1
			}
			if q.IsAnswered {
				stats["answered_questions"] = stats["answered_questions"].(int) + 1
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"role":      "visitor",
			"questions": questions,
			"stats":     stats,
		})
	}
}

// Page handlers
func (h *Handler) Home(c *gin.Context) {
	// Verificar se já está logado
	token, err := c.Cookie("auth_token")
	if err == nil && token != "" {
		claims, err := middleware.ValidateJWT(token)
		if err == nil {
			// Redirecionar baseado no role
			if claims.Role == "admin" {
				c.Redirect(http.StatusFound, "/admin")
			} else {
				c.Redirect(http.StatusFound, "/dashboard")
			}
			return
		}
	}

	c.Redirect(http.StatusFound, "/login")
}

func (h *Handler) LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"Title": "Login",
	})
}

func (h *Handler) Dashboard(c *gin.Context) {
	userID := c.MustGet("userID").(int)
	username := c.MustGet("username").(string)
	role := c.MustGet("role").(string)

	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"Title": "Dashboard",
		"User": gin.H{
			"ID":       userID,
			"Username": username,
			"Role":     role,
		},
	})
}

func (h *Handler) AdminDashboard(c *gin.Context) {
	userID := c.MustGet("userID").(int)
	username := c.MustGet("username").(string)
	role := c.MustGet("role").(string)

	c.HTML(http.StatusOK, "admin.html", gin.H{
		"Title": "Admin Dashboard",
		"User": gin.H{
			"ID":       userID,
			"Username": username,
			"Role":     role,
		},
	})
}

// API handlers específicos
func (h *Handler) GetAvailableQuestions(c *gin.Context) {
	userID := c.MustGet("userID").(int)

	questions, err := h.db.GetAvailableQuestions(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get questions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"questions": questions})
}

func (h *Handler) AnswerQuestion(c *gin.Context) {
	userID := c.MustGet("userID").(int)

	var req models.AnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Verificar se a pergunta está disponível
	question, err := h.db.GetQuestionByID(req.QuestionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if question == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}

	// Verificar se já foi respondida
	existingAnswer, err := h.db.GetAnswerByUserAndQuestion(userID, req.QuestionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if existingAnswer != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Question already answered"})
		return
	}

	// Verificar se está no horário correto
	if time.Now().Before(question.ScheduledAt) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Question not available yet"})
		return
	}

	// Criar resposta
	isCorrect := req.Answer == question.CorrectAnswer
	answer := &models.Answer{
		UserID:     userID,
		QuestionID: req.QuestionID,
		Answer:     req.Answer,
		IsCorrect:  isCorrect,
	}

	if err := h.db.CreateAnswer(answer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save answer"})
		return
	}

	response := gin.H{
		"correct": isCorrect,
		"message": "Answer submitted successfully",
	}

	if isCorrect {
		response["reward"] = question.Reward
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) GetUserStats(c *gin.Context) {
	userID := c.MustGet("userID").(int)

	stats, err := h.db.GetUserStats(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (h *Handler) GetAllQuestions(c *gin.Context) {
	questions, err := h.db.GetAllQuestions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get questions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"questions": questions})
}

func (h *Handler) UpdateQuestion(c *gin.Context) {
	// TODO: Implementar atualização de pergunta
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented yet"})
}

func (h *Handler) DeleteQuestion(c *gin.Context) {
	// TODO: Implementar exclusão de pergunta
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented yet"})
}

func (h *Handler) GetAllUsers(c *gin.Context) {
	users, err := h.db.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

func (h *Handler) GetAdminStats(c *gin.Context) {
	stats, err := h.db.GetAdminStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// Handlers para templates parciais (HTMX)
func (h *Handler) RenderStats(c *gin.Context) {
	userID := c.MustGet("userID").(int)

	stats, err := h.db.GetUserStats(userID)
	if err != nil {
		c.String(http.StatusInternalServerError, "Erro ao carregar estatísticas")
		return
	}

	c.HTML(http.StatusOK, "partials/stats.html", gin.H{
		"stats": stats,
	})
}

func (h *Handler) RenderQuestions(c *gin.Context) {
	userID := c.MustGet("userID").(int)

	questions, err := h.db.GetAvailableQuestions(userID)
	if err != nil {
		c.String(http.StatusInternalServerError, "Erro ao carregar perguntas")
		return
	}

	c.HTML(http.StatusOK, "partials/questions.html", gin.H{
		"questions": questions,
	})
}

func (h *Handler) RenderQuestionForm(c *gin.Context) {
	questionIDStr := c.Param("id")
	questionID := 0

	if _, err := fmt.Sscanf(questionIDStr, "%d", &questionID); err != nil {
		c.String(http.StatusBadRequest, "ID da pergunta inválido")
		return
	}

	userID := c.MustGet("userID").(int)

	// Verificar se a pergunta existe e está disponível
	question, err := h.db.GetQuestionByID(questionID)
	if err != nil {
		c.String(http.StatusInternalServerError, "Erro ao carregar pergunta")
		return
	}

	if question == nil {
		c.String(http.StatusNotFound, "Pergunta não encontrada")
		return
	}

	// Verificar se já foi respondida
	existingAnswer, err := h.db.GetAnswerByUserAndQuestion(userID, questionID)
	if err != nil {
		c.String(http.StatusInternalServerError, "Erro no banco de dados")
		return
	}

	if existingAnswer != nil {
		c.String(http.StatusBadRequest, "Pergunta já foi respondida")
		return
	}

	c.HTML(http.StatusOK, "partials/question-form.html", gin.H{
		"question": question,
	})
}

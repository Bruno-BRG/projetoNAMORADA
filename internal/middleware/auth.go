package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/chari/projetoNAMORADA/internal/database"
	"github.com/chari/projetoNAMORADA/internal/models"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware verifica se o usuário está autenticado
func AuthMiddleware(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, err := c.Cookie("session_id")
		if err != nil {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		var session models.Session
		var user models.User

		err = db.QueryRow(`
			SELECT s.id, s.user_id, s.expires_at, u.id, u.username, u.role 
			FROM sessions s 
			JOIN users u ON s.user_id = u.id 
			WHERE s.id = ? AND s.expires_at > datetime('now')
		`, sessionID).Scan(
			&session.ID, &session.UserID, &session.ExpiresAt,
			&user.ID, &user.Username, &user.Role,
		)

		if err != nil {
			c.SetCookie("session_id", "", -1, "/", "", false, true)
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		// Renovar sessão se estiver próxima do vencimento
		if session.ExpiresAt.Sub(time.Now()) < 24*time.Hour {
			newExpiry := time.Now().Add(7 * 24 * time.Hour)
			db.Exec("UPDATE sessions SET expires_at = ? WHERE id = ?", newExpiry, sessionID)
		}

		// Adicionar usuário ao contexto
		c.Set("user", &user)
		c.Set("session", &session)
		c.Next()
	}
}

// AdminMiddleware verifica se o usuário é admin
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
			c.Abort()
			return
		}

		if user.(*models.User).Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Acesso negado"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RateLimitMiddleware controla tentativas por IP
func RateLimitMiddleware(db *database.DB, maxAttempts int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := getClientIP(c)
		quizIDStr := c.Param("id")

		if quizIDStr == "" {
			c.Next()
			return
		}

		// Verificar rate limit
		var attempts int
		var lastAttempt, resetAt time.Time

		err := db.QueryRow(`
			SELECT attempts, last_attempt, reset_at 
			FROM rate_limits 
			WHERE ip_address = ? AND quiz_id = ?
		`, ip, quizIDStr).Scan(&attempts, &lastAttempt, &resetAt)

		if err == nil {
			// Rate limit existe
			if time.Now().Before(resetAt) && attempts >= maxAttempts {
				c.JSON(http.StatusTooManyRequests, gin.H{
					"error":    "Muitas tentativas. Tente novamente em alguns minutos.",
					"reset_at": resetAt,
				})
				c.Abort()
				return
			}

			// Reset se o window expirou
			if time.Now().After(resetAt) {
				attempts = 0
			}
		}

		c.Set("current_attempts", attempts)
		c.Set("client_ip", ip)
		c.Next()
	}
}

// getClientIP obtém o IP real do cliente considerando proxies
func getClientIP(c *gin.Context) string {
	// Verificar cabeçalhos de proxy (Cloudflare, etc.)
	if ip := c.GetHeader("CF-Connecting-IP"); ip != "" {
		return ip
	}
	if ip := c.GetHeader("X-Forwarded-For"); ip != "" {
		return strings.Split(ip, ",")[0]
	}
	if ip := c.GetHeader("X-Real-IP"); ip != "" {
		return ip
	}
	return c.ClientIP()
}

package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
}

func SecurityHeaders() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Next()
	})
}

func RequireAdminAuth() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// Verificar session de admin
		session, err := c.Cookie("admin_session")
		if err != nil || session == "" {
			c.Redirect(http.StatusFound, "/login?admin=1")
			c.Abort()
			return
		}

		// Aqui você validaria o token de sessão
		// Por simplicidade, estou usando cookie básico
		// Em produção, use JWT ou sessões mais seguras

		c.Next()
	})
}

func RequireVisitorAuth() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// Verificar session de visitante
		session, err := c.Cookie("visitor_session")
		if err != nil || session == "" {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		c.Next()
	})
}

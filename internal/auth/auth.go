package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Username string `json:"username"`
	IsAdmin  bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

var jwtSecret = []byte(getJWTSecret())

func getJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default_secret_change_in_production"
	}
	return secret
}

// GenerateToken cria um JWT token para o usuário
func GenerateToken(username string, isAdmin bool) (string, error) {
	claims := &Claims{
		Username: username,
		IsAdmin:  isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ValidateToken valida um JWT token
func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("token inválido")
	}

	return claims, nil
}

// CheckCredentials verifica as credenciais do usuário
func CheckCredentials(username, password string, isAdmin bool) bool {
	if isAdmin {
		adminUser := os.Getenv("ADMIN_USERNAME")
		adminPass := os.Getenv("ADMIN_PASSWORD")
		if adminUser == "" {
			adminUser = "admin"
		}
		if adminPass == "" {
			adminPass = "admin123"
		}
		return username == adminUser && password == adminPass
	} else {
		visitorUser := os.Getenv("VISITOR_USERNAME")
		visitorPass := os.Getenv("VISITOR_PASSWORD")
		if visitorUser == "" {
			visitorUser = "momo"
		}
		if visitorPass == "" {
			visitorPass = "momo3006"
		}
		return username == visitorUser && password == visitorPass
	}
}

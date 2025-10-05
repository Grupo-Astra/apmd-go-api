package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims define as reivindicações personalizadas para o nosso token.
//
// Ele embute as reivindicações registradas padrão (como expiração) e adiciona o ID do usuário.
type JWTClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateToken cria e assina um novo token JWT para um ID de usuário específico.
//
// O token terá uma validade de 2 horas.
func GenerateToken(userID uint) (string, error) {
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		return "", errors.New("JWT_SECRET_KEY não definida nas variáveis de ambiente")
	}

	expirationTime := time.Now().Add(2 * time.Hour)

	claims := &JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

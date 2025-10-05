// Package middleware fornece handlers Gin para processar
// requisições antes de atingirem os handlers principais
// da aplicação, como para autenticação.
package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/Grupo-Astra/apmd-go-api/auth"
	"github.com/Grupo-Astra/apmd-go-api/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWTAuthMiddleware cira um handler Gin que protege as rotas.
//
// Ele verifica a presença e validade de um token JWT no
// cabeçalho Authorization.
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.LogWarn(fmt.Sprintf(
				"Acesso negado: Cabeçalho de autorização ausente. IP: %s, Rota: %s",
				c.ClientIP(), c.Request.URL.Path,
			))
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "Cabeçalho de autorização não encontrado"},
			)
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || strings.ToLower(headerParts[0]) != "bearer" {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "Formato do token de autorização inválido"},
			)
			return
		}

		tokenString := headerParts[1]
		secretKey := os.Getenv("JWT_SECRET_KEY")
		if secretKey == "" {
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				gin.H{"error": "Chave JWT não configurada no servidor"},
			)
			return
		}

		claims := &auth.JWTClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("método de assinatura inesperado")
			}
			return []byte(secretKey), nil
		})
		if err != nil {
			errMsg := "Token inválido"
			if errors.Is(err, jwt.ErrTokenExpired) {
				errMsg = "Token expirado"
			}

			utils.LogWarn(fmt.Sprintf(
				"Acesso negado: %s. IP: %s, Rota: %s",
				errMsg, c.ClientIP(), c.Request.URL.Path,
			))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": errMsg})
			return
		}

		if !token.Valid {
			utils.LogWarn(fmt.Sprintf(
				"Acesso negado: Token marcado como inválido. IP: %s, Rota: %s",
				c.ClientIP(), c.Request.URL.Path,
			))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			return
		}

		utils.LogInfo(fmt.Sprintf(
			"Token validado para o usuário ID: %d. Acessando rota: %s",
			claims.UserID, c.Request.URL.Path,
		))
		c.Set("userID", claims.UserID)
		c.Next()
	}
}

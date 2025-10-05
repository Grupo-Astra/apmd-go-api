// Package handlers é responsável por processar as requisições
// HTTP e retornar as respostas.
package handlers

import (
	"net/http"

	"github.com/Grupo-Astra/apmd-go-api/auth"
	"github.com/Grupo-Astra/apmd-go-api/models"
	"github.com/Grupo-Astra/apmd-go-api/repositories"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AuthRequest define a estrutura esperada para uma requisição
// de login e registro.
type AuthRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// AuthHandler gerencia a lógica para as rotas de autenticação.
type AuthHandler struct {
	repo repositories.UserRepositoryInterface
}

// NewAuthHandler cria uma nova instância de AuthHandler com suas dependências.
func NewAuthHandler(repo repositories.UserRepositoryInterface) *AuthHandler {
	return &AuthHandler{repo: repo}
}

// Register processa a requisição de registro de um novo usuário.
func (h *AuthHandler) Register(c *gin.Context) {
	var req AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Requisição inválida"})
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar senha"})
		return
	}

	newUser := models.User{
		Username: req.Username,
		Password: hashedPassword,
	}

	if err := h.repo.Create(&newUser); err != nil {
		// TODO: Adicionar tratamento específico para usuário duplicado
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Não foi possível criar o usuário"},
		)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Usuário criado com sucesso"})
}

// Login processa a requisição de login e retorna um token JWT
// se as credenciais forem válidas.
func (h *AuthHandler) Login(c *gin.Context) {
	var req AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Requisição inválida"})
		return
	}

	user, err := h.repo.FindByUsername(req.Username)
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciais inválidas"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar usuário"})
		return
	}

	if !auth.CheckPasswordHash(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciais inválidas"})
		return
	}

	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Não foi possível gerar o token"},
		)
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

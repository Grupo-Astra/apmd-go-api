// Package handlers é responsável por processar as requisições
// HTTP e retornar as respostas.
package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Grupo-Astra/apmd-go-api/auth"
	"github.com/Grupo-Astra/apmd-go-api/models"
	"github.com/Grupo-Astra/apmd-go-api/repositories"
	"github.com/Grupo-Astra/apmd-go-api/utils"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
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

// ResetPasswordRequest define a estrutura para a requisição de reset de senha.
type ResetPasswordRequest struct {
	Username    string `json:"username" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
}

// MeResponse define a estrutura de resposta para o endpoint /me.
// Expõe apenas dados seguros e relevantes sobre o usuário autenticado.
type MeResponse struct {
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"createdAt"`
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
		var pqErr *pq.Error
		const PG_UNIQUE_CONSTRAINT_VIOLATION string = "23505"
		if errors.As(err, &pqErr) && pqErr.Code == pq.ErrorCode(PG_UNIQUE_CONSTRAINT_VIOLATION) {
			utils.LogWarn(fmt.Sprintf("Falha ao registrar usuário duplicado: %s", req.Username))
			c.JSON(http.StatusConflict, gin.H{"error": "Nome de usuário já existe"})
			return
		}

		utils.LogError("Erro de banco de dados ao criar usuário: " + err.Error())
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Não foi possível criar o usuário"},
		)
		return
	}

	utils.LogSuccess(fmt.Sprintf("Novo usuário registrado: %s", newUser.Username))
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
	if err != nil {
		utils.LogWarn(fmt.Sprintf(
			"Tentativa de login falhou para usuário inexistente: %s",
			req.Username,
		))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciais inválidas"})
		return
	}

	if !auth.CheckPasswordHash(req.Password, user.Password) {
		utils.LogWarn(fmt.Sprintf(
			"Tentativa de login com senha incorreta para o usuário: %s",
			user.Username,
		))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciais inválidas"})
		return
	}

	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		utils.LogError(
			"Erro ao gerar token JWT para o usuário " + user.Username + ": " + err.Error(),
		)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Não foi possível gerar o token"},
		)
		return
	}

	utils.LogInfo(fmt.Sprintf(
		"Usuário '%s' logado com sucesso.",
		user.Username,
	))
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// ResetPassword processa uma requisição pública para redefinir a senha de um usuário.
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Requisição inválida. 'username' e 'newPassword' são obrigatórios."})
		return
	}

	user, err := h.repo.FindByUsername(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.LogWarn(fmt.Sprintf("Tentativa de reset de senha para usuário inexistente: %s", req.Username))
			c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
			return
		}
		utils.LogError("Erro ao buscar usuário para reset de senha: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar requisição"})
		return
	}

	hashedPassword, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		utils.LogError("Erro ao gerar hash para reset de senha: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar nova senha"})
		return
	}

	if err := h.repo.UpdatePassword(user.ID, hashedPassword); err != nil {
		utils.LogError("Erro ao atualizar senha no banco de dados: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Não foi possível atualizar a senha"})
		return
	}

	utils.LogSuccess(fmt.Sprintf("Senha resetada com sucesso para o usuário: %s", user.Username))
	c.JSON(http.StatusOK, gin.H{"message": "Senha atualizada com sucesso."})
}

// GetMe retorna os detalhes do usuário atualmente autenticado.
//
// Ele extrai o ID do usuário do contexto Gin (injetado pelo middleware JWT).
func (h *AuthHandler) GetMe(c *gin.Context) {
	userIDUntyped, exists := c.Get("userID")
	if !exists {
		utils.LogError("Falha ao extrair 'userID' do contexto no endpoint GetMe. Isso não deveria acontecer após o middleware.")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Erro de autenticação"})
		return
	}

	userID, ok := userIDUntyped.(uint)
	if !ok {
		utils.LogError("'userID' no contexto não é do tipo uint.")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Erro interno do servidor"})
		return
	}

	user, err := h.repo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.LogWarn(fmt.Sprintf("Usuário autenticado (ID: %d) não encontrado no banco de dados.", userID))
			c.JSON(http.StatusNotFound, gin.H{"error": "Usuário autenticado não encontrado"})
			return
		}
		utils.LogError("Erro de banco de dados ao buscar usuário por ID: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar dados do usuário"})
		return
	}

	response := MeResponse{
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
	}

	c.JSON(http.StatusOK, response)
}

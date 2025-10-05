package handlers

import (
	"net/http"

	"github.com/Grupo-Astra/apmd-go-api/database"
	"github.com/Grupo-Astra/apmd-go-api/repositories"
	"github.com/Grupo-Astra/apmd-go-api/utils"
	"github.com/gin-gonic/gin"
)

// DatabaseAdminHandler gerencia as rotas para operações
// administrativas no banco de dados.
type DatabaseAdminHandler struct {
	sensorRepo repositories.SensorRepositoryInterface
	userRepo   repositories.UserRepositoryInterface
}

// NewDatabaseAdminHandler cria uma nova instância do
// DatabaseAdminHandler com suas dependências.
func NewDatabaseAdminHandler(
	sensorRepo repositories.SensorRepositoryInterface,
	userRepo repositories.UserRepositoryInterface,
) *DatabaseAdminHandler {
	return &DatabaseAdminHandler{
		sensorRepo: sensorRepo,
		userRepo:   userRepo,
	}
}

// ResetAndSeedDatabase limpa todos os dados das tabelas e executa o seeder de sensores.
func (h *DatabaseAdminHandler) ResetAndSeedDatabase(c *gin.Context) {
	utils.LogSection("Iniciando Reset do Banco de Dados")

	if err := h.sensorRepo.ClearSensorData(); err != nil {
		utils.LogError("Erro ao limpar dados de sensores: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao limpar dados de sensores"})
		return
	}

	if err := h.userRepo.ClearAll(); err != nil {
		utils.LogError("Erro ao limpar dados de usuários: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao limpar dados de usuários"})
		return
	}
	utils.LogInfo("Todas as tabelas foram limpas.")

	database.SeedSensors(h.sensorRepo)

	utils.LogSuccess("Reset do Banco de Dados concluído com sucesso.")
	c.JSON(http.StatusOK, gin.H{"message": "Banco de dados resetado e populado com sucesso."})
}

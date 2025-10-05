package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Grupo-Astra/apmd-go-api/database"
	"github.com/Grupo-Astra/apmd-go-api/models"
	"github.com/Grupo-Astra/apmd-go-api/repositories"
	"github.com/gin-gonic/gin"
)

type SensorHandler struct {
	repo repositories.SensorRepositoryInterface
}

func NewSensorHandler(repo repositories.SensorRepositoryInterface) *SensorHandler {
	return &SensorHandler{repo: repo}
}

func (h *SensorHandler) CreateSensor(c *gin.Context) {
	var sensor models.Sensor

	if err := c.ShouldBindJSON(&sensor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	history := models.SensorHistory{
		Value:     sensor.CurrentValue,
		Status:    sensor.CurrentStatus,
		Timestamp: time.Now(),
	}

	if err := h.repo.Create(&sensor, &history); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar sensor"})
		return
	}

	c.JSON(http.StatusCreated, sensor)
}

func (h *SensorHandler) GetAllSensors(c *gin.Context) {
	sensors, err := h.repo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar sensores"})
		return
	}
	c.JSON(http.StatusOK, sensors)
}

func (h *SensorHandler) GetSensorByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
	}

	sensor, err := h.repo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sensor não encontrado"})
		return
	}

	c.JSON(http.StatusOK, sensor)
}

func (h *SensorHandler) ResetAndSeedDatabase(c *gin.Context) {
	if err := h.repo.ClearAllData(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao limpar o banco de dados", "details": err.Error()})
		return
	}

	database.SeedSensors(h.repo)

	c.JSON(http.StatusOK, gin.H{"message": "Banco de dados resetado e populado com sucesso."})
}

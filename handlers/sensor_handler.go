package handlers

import (
	"net/http"
	"time"

	"github.com/Grupo-Astra/apmd-go-api/database"
	"github.com/Grupo-Astra/apmd-go-api/models"
	"github.com/gin-gonic/gin"
)

func CreateSensor(c *gin.Context) {
	var sensor models.Sensor

	if err := c.ShouldBindJSON(&sensor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&sensor).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar sensor"})
		return
	}

	history := models.SensorHistory{
		Value:     sensor.CurrentValue,
		Status:    sensor.CurrentStatus,
		Timestamp: time.Now(),
		SensorID:  sensor.ID,
	}

	if err := database.DB.Create(&history).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar histórico do sensor"})
		return
	}

	c.JSON(http.StatusCreated, sensor)
}

func GetAllSensors(c *gin.Context) {
	var sensors []models.Sensor
	if err := database.DB.Find(&sensors).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar sensores"})
		return
	}

	c.JSON(http.StatusOK, sensors)
}

func GetSensorByID(c *gin.Context) {
	var sensor models.Sensor
	id := c.Param("id")

	if err := database.DB.Preload("Historic").First(&sensor, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sensor não encontrado"})
		return
	}

	c.JSON(http.StatusOK, sensor)
}

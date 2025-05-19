package handlers

import (
	"net/http"

	"github.com/Grupo-Astra/apmd-go-api/database"
	"github.com/Grupo-Astra/apmd-go-api/models"
	"github.com/gin-gonic/gin"
)

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

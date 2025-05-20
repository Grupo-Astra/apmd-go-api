package utils

import (
	"time"

	"github.com/Grupo-Astra/apmd-go-api/database"
	"github.com/Grupo-Astra/apmd-go-api/models"
)

func StartSensorSimulation(interval time.Duration) {
	for {
		var sensors []models.Sensor
		if err := database.DB.Find(&sensors).Error; err == nil {
			for i := range sensors {
				UpdateSensorData(&sensors[i])
			}
		}
		time.Sleep(interval)
	}
}

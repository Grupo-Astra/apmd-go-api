package utils

import (
	"log"
	"math/rand"
	"time"

	"github.com/Grupo-Astra/apmd-go-api/database"
	"github.com/Grupo-Astra/apmd-go-api/models"
)

func UpdateSensorData(sensor *models.Sensor) error {
	newValue := rand.Float64() * 100

	newStatus := "OK"
	if newValue <= 20 || newValue >= 70 {
		newStatus = "Alerta"
	}

	sensor.CurrentValue = newValue
	sensor.CurrentStatus = newStatus

	if err := database.DB.Save(sensor).Error; err != nil {
		return err
	}

	history := models.SensorHistory{
		Value:     newValue,
		Status:    newStatus,
		Timestamp: time.Now(),
		SensorID:  sensor.ID,
	}

	if err := database.DB.Create(&history).Error; err != nil {
		return err
	}

	log.Printf(
		"Sensor %d atualizado: valor=%.2f, status=%s",
		sensor.ID,
		newValue,
		newStatus,
	)
	return nil
}

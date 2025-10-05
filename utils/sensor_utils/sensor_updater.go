package sensorutils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/Grupo-Astra/apmd-go-api/models"
	"github.com/Grupo-Astra/apmd-go-api/repositories"
	"github.com/Grupo-Astra/apmd-go-api/utils"
)

func UpdateSensorData(repo repositories.SensorRepositoryInterface, sensor *models.Sensor) error {
	newValue, newStatus := generateSensorReading(sensor.Name)

	sensor.CurrentValue = newValue
	sensor.CurrentStatus = newStatus

	history := models.SensorHistory{
		Value:     newValue,
		Status:    newStatus,
		Timestamp: time.Now(),
		SensorID:  sensor.ID,
	}

	if err := repo.Update(sensor, &history); err != nil {
		return err
	}

	utils.LogInfo(
		fmt.Sprintf(
			"Sensor atualizado: [%s] valor=%.2f | status=%s",
			sensor.Name, newValue, newStatus,
		),
	)

	return nil
}

func generateSensorReading(sensorName string) (float64, string) {
	name := strings.ToLower(sensorName)
	var value float64
	var status string

	switch {
	case strings.Contains(name, "pressão"):
		value = 5 + rand.Float64()*3
		if value < 7 {
			status = "OK"
		} else {
			status = "Alerta"
		}

	case strings.Contains(name, "curso"):
		value = rand.Float64() * 100
		status = "OK"

	case strings.Contains(name, "ciclos"):
		value = sensorNameToCycleCount(sensorName) + 1
		status = "OK"

	case strings.Contains(name, "força"):
		value = rand.Float64() * 150
		if value < 100 {
			status = "OK"
		} else {
			status = "Alerta"
		}

	case strings.Contains(name, "vazamento"):
		value = rand.Float64() * 3
		if value < 1 {
			status = "OK"
		} else {
			status = "Alerta"
		}

	default:
		value = rand.Float64() * 100
		status = "OK"
	}

	return value, status
}

var fakeCounter = make(map[string]float64)

func sensorNameToCycleCount(name string) float64 {
	fakeCounter[name] += 1
	return fakeCounter[name]
}

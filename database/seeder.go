package database

import (
	"github.com/Grupo-Astra/apmd-go-api/models"
	"github.com/Grupo-Astra/apmd-go-api/utils"
)

func SeedSensors() {
	utils.LogSection("Seeder de Sensores")

	var count int64
	if err := DB.Model(&models.Sensor{}).Count(&count).Error; err != nil {
		utils.LogError("Erro ao verificar sensorse: " + err.Error())
		return
	}

	if count > 0 {
		utils.LogInfo("Já existem sensores cadastrados no banco. Seeder ignorado.")
		return
	}

	sensors := []models.Sensor{
		{
			Name:          "Sensor de Pressão",
			CurrentValue:  6.2,
			CurrentStatus: "OK",
		},
		{
			Name:          "Sensor de Curso (Posição)",
			CurrentValue:  100,
			CurrentStatus: "OK",
		},
		{
			Name:          "Contador de Ciclos",
			CurrentValue:  0,
			CurrentStatus: "OK",
		},
		{
			Name:          "Sensor de Força",
			CurrentValue:  75.0,
			CurrentStatus: "OK",
		},
		{
			Name:          "Sensor de Vazamento de Ar",
			CurrentValue:  0,
			CurrentStatus: "OK",
		},
	}

	for _, sensor := range sensors {
		if err := DB.Create(&sensor).Error; err != nil {
			utils.LogError("Erro ao inserir sensor: " + err.Error())
			return
		}

		history := models.SensorHistory{
			Value:    sensor.CurrentValue,
			Status:   sensor.CurrentStatus,
			SensorID: sensor.ID,
		}

		if err := DB.Create(&history).Error; err != nil {
			utils.LogError(
				"Erro ao inserir histórico inicial para " +
					sensor.Name + ": " + err.Error(),
			)
			continue
		}

		utils.LogSuccess("Sensor criado: " + sensor.Name)
	}

	utils.LogSuccess("Sensores iniciais populados com sucesso.")
}

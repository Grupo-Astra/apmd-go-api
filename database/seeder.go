package database

import (
	"log"

	"github.com/Grupo-Astra/apmd-go-api/models"
)

func SeedSensors() {
	var count int64
	if err := DB.Model(&models.Sensor{}).Count(&count).Error; err != nil {
		log.Println("Erro ao verificar sensorse:", err)
		return
	}

	if count > 0 {
		log.Println("Já existem sensores cadastrados no banco. Seeder ignorado.")
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
			log.Println("Erro ao inserir sensor:", err)
			return
		}

		history := models.SensorHistory{
			Value:    sensor.CurrentValue,
			Status:   sensor.CurrentStatus,
			SensorID: sensor.ID,
		}

		if err := DB.Create(&history).Error; err != nil {
			log.Println("Erro ao inserir histórico inicial:", err)
		}
	}

	log.Println("Sensores iniciais populados com sucesso.")
}

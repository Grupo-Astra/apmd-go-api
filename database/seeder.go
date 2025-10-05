// Package database gerencia a conexão com o banco
// de dados e as migrações de schema
package database

import (
	"github.com/Grupo-Astra/apmd-go-api/models"
	"github.com/Grupo-Astra/apmd-go-api/repositories"
	"github.com/Grupo-Astra/apmd-go-api/utils"
)

// SeedSensors popula o banco de dados com um conjunto inicial de sensores.
//
// A função verifica se já existem sensores antes de executar para evitar duplicação.
//
// Requer uma instância do repositório de sensores para interagir com o banco de dados.
func SeedSensors(repo repositories.SensorRepositoryInterface) {
	utils.LogSection("Seeder de Sensores")

	count, err := repo.Count()
	if err != nil {
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
		history := models.SensorHistory{
			Value:  sensor.CurrentValue,
			Status: sensor.CurrentStatus,
		}

		if err := repo.Create(&sensor, &history); err != nil {
			utils.LogError("Erro ao inserir sensor '" + sensor.Name + "': " + err.Error())
			return
		}

		utils.LogSuccess("Sensor criado: " + sensor.Name)
	}

	utils.LogSuccess("Sensores iniciais populados com sucesso.")
}

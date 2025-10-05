package sensorutils

import (
	"fmt"
	"time"

	"github.com/Grupo-Astra/apmd-go-api/repositories"
	"github.com/Grupo-Astra/apmd-go-api/utils"
)

func StartSensorSimulation(repo repositories.SensorRepositoryInterface, interval time.Duration) {
	utils.LogInfo("Iniciando simulação de atualização de sensores...")
	for {
		sensors, err := repo.FindAll()
		if err != nil {
			utils.LogError(fmt.Sprintf("Erro ao buscar sensores para simulação: %v", err))
		} else {
			for i := range sensors {
				if err := UpdateSensorData(repo, &sensors[i]); err != nil {
					utils.LogError(
						fmt.Sprintf(
							"Erro ao atualizar sensor [%s]: %v", sensors[i].Name,
							err,
						),
					)
				}
			}
		}
		time.Sleep(interval)
	}
}

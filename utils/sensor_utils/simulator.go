package sensorutils

import (
	"log"
	"time"

	"github.com/Grupo-Astra/apmd-go-api/repositories"
)

func StartSensorSimulation(repo repositories.SensorRepositoryInterface, interval time.Duration) {
	log.Println("Iniciando simulação de atualização de sensores...")
	for {
		sensors, err := repo.FindAll()
		if err != nil {
			log.Printf("Erro ao buscar sensores para simulação: %v", err)
		} else {
			for i := range sensors {
				if err := UpdateSensorData(repo, &sensors[i]); err != nil {
					log.Printf("Erro ao atualizar sensor [%s]: %v", sensors[i].Name, err)
				}
			}
		}
		time.Sleep(interval)
	}
}

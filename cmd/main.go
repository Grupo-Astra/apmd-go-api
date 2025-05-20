package main

import (
	"log"
	"time"

	"github.com/Grupo-Astra/apmd-go-api/database"
	"github.com/Grupo-Astra/apmd-go-api/routes"
	"github.com/Grupo-Astra/apmd-go-api/utils"
)

func main() {
	database.InitDatabase()

	router := routes.SetupRouter()

	go utils.StartSensorSimulation(5 * time.Second)

	log.Println("Servidor inicializado na porta :8080")
	router.Run(":8080")
}

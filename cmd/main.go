package main

import (
	"log"
	"time"

	"github.com/Grupo-Astra/apmd-go-api/database"
	"github.com/Grupo-Astra/apmd-go-api/routes"
	sensorutils "github.com/Grupo-Astra/apmd-go-api/utils/sensor_utils"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	log.Println("Inicializando banco de dados...")
	database.InitDatabase()

	log.Println("Inicializando seeder do banco de dados...")
	database.SeedSensors()

	router := routes.SetupRouter()

	go sensorutils.StartSensorSimulation(5 * time.Second)

	log.Println("Servidor inicializado na porta :8080")
	router.Run(":8080")
}

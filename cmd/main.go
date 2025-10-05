package main

import (
	"time"

	"github.com/Grupo-Astra/apmd-go-api/database"
	"github.com/Grupo-Astra/apmd-go-api/repositories"
	"github.com/Grupo-Astra/apmd-go-api/routes"
	"github.com/Grupo-Astra/apmd-go-api/utils"
	sensorutils "github.com/Grupo-Astra/apmd-go-api/utils/sensor_utils"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	utils.LogInfo("Inicializando banco de dados...")
	database.InitDatabase()

	sensorRepository := repositories.NewSensorRepository(
		database.PostgresDB,
	)

	database.SeedSensors(sensorRepository)

	router := routes.SetupRouter(sensorRepository)

	go sensorutils.StartSensorSimulation(sensorRepository, 5*time.Second)

	utils.LogSuccess("Servidor inicializado na porta :8080")
	router.Run(":8080")
}

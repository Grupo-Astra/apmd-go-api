package routes

import (
	"time"

	"github.com/Grupo-Astra/apmd-go-api/handlers"
	"github.com/Grupo-Astra/apmd-go-api/repositories"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(sensorRepo repositories.SensorRepositoryInterface) *gin.Engine {
	r := gin.Default()

	sensorHandler := handlers.NewSensorHandler(sensorRepo)

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8081"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	apiV1 := r.Group("/api")
	{
		apiV1.GET("api/readings", sensorHandler.GetAllSensors)
		apiV1.GET("api/readings/:id", sensorHandler.GetSensorByID)
		apiV1.POST("api/readings", sensorHandler.CreateSensor)
	}

	return r
}

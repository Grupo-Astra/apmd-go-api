package routes

import (
	"time"

	"github.com/Grupo-Astra/apmd-go-api/handlers"
	"github.com/Grupo-Astra/apmd-go-api/repositories"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(
	sensorRepo repositories.SensorRepositoryInterface,
	userRepo repositories.UserRepositoryInterface,
) *gin.Engine {
	r := gin.Default()

	sensorHandler := handlers.NewSensorHandler(sensorRepo)
	authHandler := handlers.NewAuthHandler(userRepo)

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
		authRoutes := apiV1.Group("/auth")
		{
			authRoutes.POST("/register", authHandler.Register)
			authRoutes.POST("/login", authHandler.Login)
		}

		readings := apiV1.Group("/readings")
		{
			readings.GET("", sensorHandler.GetAllSensors)
			readings.GET("/:id", sensorHandler.GetSensorByID)
			readings.POST("", sensorHandler.CreateSensor)
		}

		dbAdmin := apiV1.Group("/database")
		{
			dbAdmin.POST("/reset", sensorHandler.ResetAndSeedDatabase)
		}
	}

	return r
}

// Package routes é responsável pela configuração e definição
// de todas as rotas da API.
package routes

import (
	"time"

	"github.com/Grupo-Astra/apmd-go-api/handlers"
	"github.com/Grupo-Astra/apmd-go-api/middleware"
	"github.com/Grupo-Astra/apmd-go-api/repositories"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter configura o motor do Gin, atribui permissões de CORS,
// agrupa as rotas e as associa aos seus respectivos handlers.
func SetupRouter(
	sensorRepo repositories.SensorRepositoryInterface,
	userRepo repositories.UserRepositoryInterface,
) *gin.Engine {
	r := gin.Default()

	sensorHandler := handlers.NewSensorHandler(sensorRepo)
	authHandler := handlers.NewAuthHandler(userRepo)
	databaseHandler := handlers.NewDatabaseAdminHandler(sensorRepo, userRepo)

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8081"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := r.Group("/api")
	{
		authRoutes := api.Group("/auth")
		{
			authRoutes.POST("/register", authHandler.Register)
			authRoutes.POST("/login", authHandler.Login)
		}

		// Legado -> rotas públicas
		readingsV1 := api.Group("/readings")
		{
			readingsV1.GET("", sensorHandler.GetAllSensors)
			readingsV1.GET("/:id", sensorHandler.GetSensorByID)
			readingsV1.POST("", sensorHandler.CreateSensor)
		}

		dbAdmin := api.Group("/database")
		{
			dbAdmin.POST("/reset", databaseHandler.ResetAndSeedDatabase)
		}

		v2 := api.Group("/v2")
		v2.Use(middleware.JWTAuthMiddleware())
		{
			readingsV2 := v2.Group("/readings")
			readingsV2.GET("", sensorHandler.GetAllSensors)
			readingsV2.GET("/:id", sensorHandler.GetSensorByID)
			readingsV2.POST("", sensorHandler.CreateSensor)
		}
	}

	return r
}

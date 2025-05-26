package routes

import (
	"time"

	"github.com/Grupo-Astra/apmd-go-api/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8081"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/sensors", handlers.GetAllSensors)
	r.GET("/sensors/:id", handlers.GetSensorByID)
	r.POST("/sensors", handlers.CreateSensor)

	return r
}

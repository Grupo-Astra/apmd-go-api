package routes

import (
	"github.com/Grupo-Astra/apmd-go-api/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/sensors", handlers.GetAllSensors)
	r.GET("/sensors/:id", handlers.GetSensorByID)
	r.POST("/sensors", handlers.CreateSensor)

	return r
}

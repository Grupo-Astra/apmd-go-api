package main

import (
	"github.com/Grupo-Astra/apmd-go-api/database"
	"github.com/Grupo-Astra/apmd-go-api/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	database.InitDatabase()

	router := gin.Default()

	router.GET("/sensors", handlers.GetAllSensors)
	router.GET("/sensors/:id", handlers.GetSensorByID)

	router.Run(":8080")
}

package main

import (
	"github.com/Grupo-Astra/apmd-go-api/database"
	"github.com/gin-gonic/gin"
)

func main() {
	database.InitDatabase()

	router := gin.Default()

	router.Run(":8080")
}

package main

import (
	"log"

	"github.com/Grupo-Astra/apmd-go-api/database"
	"github.com/Grupo-Astra/apmd-go-api/routes"
)

func main() {
	database.InitDatabase()

	router := routes.SetupRouter()

	log.Println("Servidor inicializado na porta :8080")
	router.Run(":8080")
}

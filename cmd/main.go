package main

import (
	"log"

	"github.com/Grupo-Astra/apmd-go-api/config"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar .env")
	}

	config.ConnectToOracle()
	log.Println("Servidor iniciado na porta 8080")
}

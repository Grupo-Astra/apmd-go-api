package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Grupo-Astra/apmd-go-api/config"
	"github.com/Grupo-Astra/apmd-go-api/migrations"
	"github.com/Grupo-Astra/apmd-go-api/routes"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar .env")
	}

	config.ConnectToOracle()
	defer config.DB.Close()

	migrations.RunMigrations()

	router := routes.SetupRouter()

	go func() {
		if err := router.Run(":8080"); err != nil {
			log.Fatalf("Erro ao iniciar o servidor: %v", err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	log.Println("Encerrando aplicação...")
}

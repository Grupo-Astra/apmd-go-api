package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/sijms/go-ora/v2"
)

var DB *sql.DB

func ConnectToOracle() {
	var (
		username = os.Getenv("ORACLE_USER")
		password = os.Getenv("ORACLE_PASSWORD")
		host     = "oracle.fiap.com.br"
		port     = 1521
		sid      = "orcl"
	)

	dsn := fmt.Sprintf("oracle://%s:%s@%s:%d/%s", username, password, host, port, sid)

	var err error
	DB, err = sql.Open("oracle", dsn)
	if err != nil {
		log.Fatalf("Erro ao abrir conexão com o Oracle: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("Erro ao conectar ao Oracle: %v", err)
	}

	log.Println("Conectado ao Oracle com sucesso")
}

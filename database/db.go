// Package database gerencia a conexão com o banco
// de dados e as migrações de schema
package database

import (
	"os"

	"github.com/Grupo-Astra/apmd-go-api/models"
	"github.com/Grupo-Astra/apmd-go-api/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// PostgresDB é a instância global da conexão com o banco de dados.
var PostgresDB *gorm.DB

// InitDatabase inicializa a conexão com o banco de dados PostgreSQL
// e executa o AutoMigrate do GORM para garantir que o esquema esteja atualizado.
//
// A aplicação será encerrada se a conexão ou a migração falharem.
func InitDatabase() {
	var err error

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		utils.LogFatal("Variável de ambiente DATABASE_URL não definida.")
	}

	PostgresDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		utils.LogFatal("Erro ao conectar ao PostgresSQL: " + err.Error())
	}
	utils.LogSuccess("Conexão com PostgresSQL estabelecida com sucesso.")

	utils.LogInfo("Migrando tabelas para o PostgresSQL...")
	err = PostgresDB.AutoMigrate(&models.Sensor{}, &models.SensorHistory{}, &models.User{})
	if err != nil {
		utils.LogFatal("PostgresSQL - Erro ao migrar tabelas: " + err.Error())
	}

	utils.LogSuccess("Tabelas migradas com sucesso.")
}

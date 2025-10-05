package database

import (
	"os"

	"github.com/Grupo-Astra/apmd-go-api/models"
	"github.com/Grupo-Astra/apmd-go-api/utils"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	DB         *gorm.DB
	PostgresDB *gorm.DB
)

func InitDatabase() {
	var err error
	DB, err = gorm.Open(sqlite.Open("data/readings.db"), &gorm.Config{})
	if err != nil {
		utils.LogFatal("Erro ao conectar ao SQLite: " + err.Error())
	}
	utils.LogSuccess("Conexão com SQLite estabelecida com sucesso.")

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		utils.LogFatal("Variável de ambiente DATABASE_URL não definida.")
	}

	PostgresDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		utils.LogFatal("Erro ao conectar ao PostgresSQL: " + err.Error())
	}
	utils.LogSuccess("Conexão com PostgresSQL estabelecida com sucesso.")

	utils.LogInfo("Migrando tabelas para o SQLite...")
	err = DB.AutoMigrate(&models.Sensor{}, &models.SensorHistory{})
	if err != nil {
		utils.LogFatal("SQLite - Erro ao migrar tabelas: " + err.Error())
	}

	utils.LogInfo("Migrando tabelas para o PostgresSQL...")
	err = PostgresDB.AutoMigrate(&models.Sensor{}, &models.SensorHistory{})
	if err != nil {
		utils.LogFatal("PostgresSQL - Erro ao migrar tabelas: " + err.Error())
	}

	utils.LogSuccess("Tabelas migradas com sucesso em ambos os bancos de dados.")
}

package database

import (
	"log"
	"os"

	"github.com/Grupo-Astra/apmd-go-api/models"
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
		log.Fatal("Erro ao conectar ao SQLite: ", err)
	}
	log.Println("Conexão com SQLite estabelecida com sucesso.")

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("Variável de ambiente DATABASE_URL não definida.")
	}

	PostgresDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Erro ao conectar ao PostgresSQL: ", err)
	}
	log.Println("Conexão com PostgresSQL estabelecida com sucesso.")

	log.Println("Migrando tabelas para o SQLite...")
	err = DB.AutoMigrate(&models.Sensor{}, &models.SensorHistory{})
	if err != nil {
		log.Fatal("SQLite - Erro ao migrar tabelas: ", err)
	}

	log.Println("Migrando tabelas para o PostgresSQL...")
	err = PostgresDB.AutoMigrate(&models.Sensor{}, &models.SensorHistory{})
	if err != nil {
		log.Fatal("PostgresSQL - Erro ao migrar tabelas: ", err)
	}

	log.Println("Tabelas migradas com sucesso em ambos os bancos de dados.")
}

package database

import (
	"log"

	"github.com/Grupo-Astra/apmd-go-api/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() {
	var err error
	DB, err = gorm.Open(sqlite.Open("data/readings.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Erro ao conectar ao SQLite: ", err)
	}

	err = DB.AutoMigrate(&models.Sensor{}, &models.SensorHistory{})
	if err != nil {
		log.Fatal("Erro ao migrar tabelas: ", err)
	}
}

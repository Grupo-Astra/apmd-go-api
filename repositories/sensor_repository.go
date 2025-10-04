package repositories

import (
	"github.com/Grupo-Astra/apmd-go-api/models"
	"gorm.io/gorm"
)

type SensorRepositoryInterface interface {
	Create(sensor *models.Sensor, history *models.SensorHistory) error
}

type sensorRepository struct {
	postgresDB *gorm.DB
	sqliteDB   *gorm.DB
}

func NewSensorRepository(postgres *gorm.DB, sqlite *gorm.DB) SensorRepositoryInterface {
	return &sensorRepository{
		postgresDB: postgres,
		sqliteDB:   sqlite,
	}
}

func (r *sensorRepository) Create(sensor *models.Sensor, history *models.SensorHistory) error {
	// TODO
	return nil
}

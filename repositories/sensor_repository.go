package repositories

import (
	"github.com/Grupo-Astra/apmd-go-api/models"
	"gorm.io/gorm"
)

type SensorRepositoryInterface interface {
	Create(sensor *models.Sensor, history *models.SensorHistory) error
	FindAll() ([]models.Sensor, error)
	FindByID(id int) (models.Sensor, error)
	Update(sensor *models.Sensor, history *models.SensorHistory) error
	Count() (int64, error)
	ClearSensorData() error
}

type sensorRepository struct {
	postgresDB *gorm.DB
}

func NewSensorRepository(postgres *gorm.DB) SensorRepositoryInterface {
	return &sensorRepository{
		postgresDB: postgres,
	}
}

func (r *sensorRepository) Create(sensor *models.Sensor, history *models.SensorHistory) error {
	return r.postgresDB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(sensor).Error; err != nil {
			return err
		}
		history.SensorID = sensor.ID
		if err := tx.Create(history).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *sensorRepository) FindAll() ([]models.Sensor, error) {
	var sensors []models.Sensor
	err := r.postgresDB.Find(&sensors).Error
	return sensors, err
}

func (r *sensorRepository) FindByID(id int) (models.Sensor, error) {
	var sensor models.Sensor
	err := r.postgresDB.Preload("Historic").First(&sensor, id).Error
	return sensor, err
}

func (r *sensorRepository) Update(sensor *models.Sensor, history *models.SensorHistory) error {
	return r.postgresDB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(sensor).Error; err != nil {
			return err
		}
		history.SensorID = sensor.ID
		if err := tx.Create(history).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *sensorRepository) Count() (int64, error) {
	var count int64
	err := r.postgresDB.Model(&models.Sensor{}).Count(&count).Error
	return count, err
}

// ClearSensorData remove todos os registros das tabelas de sensor e histórico.
func (r *sensorRepository) ClearSensorData() error {
	return r.postgresDB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("DELETE FROM sensor_histories").Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM sensors").Error; err != nil {
			return err
		}
		return nil
	})
}

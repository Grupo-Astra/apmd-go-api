// Package repositories contém as abstrações para
// acesso ao banco de dados, seguindo o Repository Pattern.
package repositories

import (
	"github.com/Grupo-Astra/apmd-go-api/models"
	"gorm.io/gorm"
)

// SensorRepositoryInterface define o contrato
// para as operações de banco de dados relacionadas a sensores.
type SensorRepositoryInterface interface {
	Create(sensor *models.Sensor, history *models.SensorHistory) error
	FindAll() ([]models.Sensor, error)
	FindByID(id int) (models.Sensor, error)
	Update(sensor *models.Sensor, history *models.SensorHistory) error
	Count() (int64, error)
	ClearSensorData() error
}

// sensorRepository é a implementação concreta da SensorRepositoryInterface.
type sensorRepository struct {
	postgresDB *gorm.DB
}

// NewSensorRepository cria uma nova instância do repositório de sensores.
func NewSensorRepository(postgres *gorm.DB) SensorRepositoryInterface {
	return &sensorRepository{
		postgresDB: postgres,
	}
}

// Create insere um novo sensor e seu primeiro registro
// de histórico no banco de dados.
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

// FindAll retorna todos os sensores cadastrados.
func (r *sensorRepository) FindAll() ([]models.Sensor, error) {
	var sensors []models.Sensor
	err := r.postgresDB.Find(&sensors).Error
	return sensors, err
}

// FindByID busca um sensor específico pelo seu ID, incluindo seu histórico.
func (r *sensorRepository) FindByID(id int) (models.Sensor, error) {
	var sensor models.Sensor
	err := r.postgresDB.Preload("Historic").First(&sensor, id).Error
	return sensor, err
}

// Update atualiza os dados de um sensor existente e cria um novo
// registro de histórico.
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

// Count retorna a contagem total de sensores no banco de dados.
func (r *sensorRepository) Count() (int64, error) {
	var count int64
	err := r.postgresDB.Model(&models.Sensor{}).Count(&count).Error
	return count, err
}

// ClearSensorData remove todos os registros das tabelas de sensores e históricos.
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

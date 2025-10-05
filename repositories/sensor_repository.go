package repositories

import (
	"fmt"

	"github.com/Grupo-Astra/apmd-go-api/models"
	"github.com/Grupo-Astra/apmd-go-api/utils"
	"gorm.io/gorm"
)

type SensorRepositoryInterface interface {
	Create(sensor *models.Sensor, history *models.SensorHistory) error
	FindAll() ([]models.Sensor, error)
	FindByID(id int) (models.Sensor, error)
	Update(sensor *models.Sensor, history *models.SensorHistory) error
	Count() (int64, error)
	ClearAllData() error
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
	txPostgres := r.postgresDB.Begin()
	txSqlite := r.sqliteDB.Begin()

	defer func() {
		if r := recover(); r != nil {
			txPostgres.Rollback()
			txSqlite.Rollback()
			utils.LogWarn(fmt.Sprintf("Pane recuperada, transações revertidas: ", r))
		}
	}()

	if err := txPostgres.Create(sensor).Error; err != nil {
		txPostgres.Rollback()
		txSqlite.Rollback()
		return err
	}

	history.SensorID = sensor.ID
	if err := txPostgres.Create(history).Error; err != nil {
		txPostgres.Rollback()
		txSqlite.Rollback()
		return err
	}

	if err := txSqlite.Create(sensor).Error; err != nil {
		txPostgres.Rollback()
		txSqlite.Rollback()
		return err
	}
	if err := txSqlite.Create(history).Error; err != nil {
		txPostgres.Rollback()
		txSqlite.Rollback()
		return err
	}

	if err := txPostgres.Commit().Error; err != nil {
		return err
	}
	if err := txSqlite.Commit().Error; err != nil {
		utils.LogError(
			fmt.Sprintf(
				"CRÍTICO: Falha ao commitar transação no SQLite após sucesso no PostgreSQL: %v",
				err,
			),
		)
		return err
	}

	return nil
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
	txPostgres := r.postgresDB.Begin()
	txSqlite := r.sqliteDB.Begin()

	defer func() {
		if r := recover(); r != nil {
			txPostgres.Rollback()
			txSqlite.Rollback()
		}
	}()

	if err := txPostgres.Save(sensor).Error; err != nil {
		txPostgres.Rollback()
		txSqlite.Rollback()
		return err
	}
	if err := txSqlite.Save(sensor).Error; err != nil {
		txPostgres.Rollback()
		txSqlite.Rollback()
		return err
	}

	history.SensorID = sensor.ID
	if err := txPostgres.Create(history).Error; err != nil {
		txPostgres.Rollback()
		txSqlite.Rollback()
		return err
	}
	if err := txSqlite.Create(history).Error; err != nil {
		txPostgres.Rollback()
		txSqlite.Rollback()
		return err
	}

	if err := txPostgres.Commit().Error; err != nil {
		return err
	}
	if err := txSqlite.Commit().Error; err != nil {
		utils.LogError(
			fmt.Sprintf(
				"CRITICAL: Falha ao commitar transação no SQLite após sucesso no PostgreSQL: %v",
				err,
			),
		)
		return err
	}

	return nil
}

func (r *sensorRepository) Count() (int64, error) {
	var count int64
	err := r.postgresDB.Model(&models.Sensor{}).Count(&count).Error
	return count, err
}

func (r *sensorRepository) ClearAllData() error {
	txPostgres := r.postgresDB.Begin()
	txSqlite := r.sqliteDB.Begin()

	if err := txPostgres.Exec("DELETE FROM sensor_histories").Error; err != nil {
		txPostgres.Rollback()
		txSqlite.Rollback()
		return err
	}
	if err := txPostgres.Exec("DELETE FROM sensors").Error; err != nil {
		txPostgres.Rollback()
		txSqlite.Rollback()
		return err
	}

	if err := txSqlite.Exec("DELETE FROM sensor_histories").Error; err != nil {
		txPostgres.Rollback()
		txSqlite.Rollback()
		return err
	}
	if err := txSqlite.Exec("DELETE FROM sensors").Error; err != nil {
		txPostgres.Rollback()
		txSqlite.Rollback()
		return err
	}

	if err := txPostgres.Commit().Error; err != nil {
		return err
	}
	if err := txSqlite.Commit().Error; err != nil {
		utils.LogError(
			fmt.Sprintf(
				"CRÍTICO: Falha ao commitar limpeza no SQLite após sucesso no PostgreSQL: %v",
				err,
			),
		)
		return err
	}

	return nil
}

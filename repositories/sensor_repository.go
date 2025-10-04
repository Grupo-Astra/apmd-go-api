package repositories

import (
	"log"

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
	txPostgres := r.postgresDB.Begin()
	txSqlite := r.sqliteDB.Begin()

	defer func() {
		if r := recover(); r != nil {
			txPostgres.Rollback()
			txSqlite.Rollback()
			log.Println("Pane recuperada, transações revertidas: ", r)
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
		log.Printf("CRÍTICO: Falha ao commitar transação no SQLite após sucesso no PostgreSQL: %v", err)
		return err
	}

	return nil
}

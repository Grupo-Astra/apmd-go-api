package migrations

import (
	"log"

	"github.com/Grupo-Astra/apmd-go-api/config"
)

func RunMigrations() {
	createSensorsTable := `
	CREATE TABLE SENSORS (
		SENSOR_ID RAW(16) DEFAULT SYS_GUID() PRIMARY KEY,
		NAME VARCHAR2(100) NOT NULL,
		CURRENT_VALUE NUMBER(10,2) NOT NULL,
		CURRENT_STATUS VARCHAR2(6) NOT NULL
	);`

	_, err := config.DB.Exec(createSensorsTable)
	if err != nil {
		log.Println("Tabela SENSORS já existe ou erro: ", err)
	}

	createSensorHistoryTable := `
	CREATE TABLE SENSOR_HISTORY (
		SENSOR_HISTORY_ID RAW(16) DEFAULT SYS_GUID() PRIMARY KEY,
		VALUE NUMBER(10,2) NOT NULL,
		STATUS VARCHAR2(6) NOT NULL,
		TIMESTAMP TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
		SENSOR_ID RAW(16) NOT NULL,
		CONSTRAINT SENSOR_HISTORY_SENSORS_FK
			FOREIGN KEY (SENSOR_ID)
			REFERENCES SENSORS (SENSOR_ID)
	);`

	_, err = config.DB.Exec(createSensorHistoryTable)
	if err != nil {
		log.Println("Tabela SENSOR_HISTORY pode já existir ou erro:", err)
	}

	log.Println("Migrações executadas com sucesso.")
}

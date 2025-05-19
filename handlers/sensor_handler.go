package handlers

import (
	"database/sql"
	"encoding/hex"
	"log"
	"net/http"
	"strings"

	"github.com/Grupo-Astra/apmd-go-api/config"
	"github.com/Grupo-Astra/apmd-go-api/models"
	"github.com/gin-gonic/gin"
)

func CreateSensor(c *gin.Context) {
	var input models.Sensor

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}

	status := strings.ToUpper(input.CurrentStatus)
	if status != "OK" && status != "ALERTA" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status deve ser 'OK' ou 'Alerta'"})
		return
	}

	var sensorIDRaw []byte
	insertSensorSQL := `
		INSERT INTO SENSORS (SENSOR_ID, NAME, CURRENT_VALUE, CURRENT_STATUS)
		VALUES (SYS_GUID(), :1, :2, :3)
		RETURNING SENSOR_ID INTO :4
	`

	sensorIDRaw = make([]byte, 16)
	_, err := config.DB.Exec(insertSensorSQL, input.Name, input.CurrentValue, status, &sensorIDRaw)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao inserir sensor"})
		return
	}

	insertHistSQL := `
		INSERT INTO SENSOR_HISTORY (SENSOR_HISTORY_ID, VALUE, STATUS, TIMESTAMP, SENSOR_ID)
		VALUES (SYS_GUID(), :1, :2, CURRENT_TIMESTAMP, :3)
	`

	_, err = config.DB.Exec(insertHistSQL, input.CurrentValue, status, sensorIDRaw)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao inserir histórico"})
		return
	}

	input.SensorID = hex.EncodeToString(sensorIDRaw)
	input.CurrentStatus = status

	c.JSON(http.StatusCreated, input)
}

func GetAllSensors(c *gin.Context) {
	rows, err := config.DB.Query(
		"SELECT SENSOR_ID, NAME, CURRENT_VALUE, CURRENT_STATUS FROM SENSORS",
	)
	if err != nil {
		log.Printf("Erro: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao consultar sensores"})
		return
	}
	defer rows.Close()

	var sensors []models.Sensor

	for rows.Next() {
		var s models.Sensor
		var rawID []byte

		if err := rows.Scan(&rawID, &s.Name, &s.CurrentValue, &s.CurrentStatus); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao mapear sensor"})
			return
		}
		s.SensorID = hex.EncodeToString(rawID)
		sensors = append(sensors, s)
	}

	c.JSON(http.StatusOK, sensors)
}

func GetSensorByID(c *gin.Context) {
	id := c.Param("id")

	rawID, err := hex.DecodeString(id)
	if err != nil || len(rawID) != 16 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var s models.Sensor
	row := config.DB.QueryRow(
		"SELECT SENSOR_ID, NAME, CURRENT_VALUE, CURRENT_STATUS FROM SENSORS WHERE SENSOR_ID = :1",
		rawID,
	)

	var sensorID []byte
	if err := row.Scan(&sensorID, &s.Name, &s.CurrentValue, &s.CurrentStatus); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Sensor não encontrado"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao consultar sensor"})
		}
		return
	}
	s.SensorID = hex.EncodeToString(sensorID)

	rows, err := config.DB.Query(`
		SELECT SENSOR_HISTORY_ID, VALUE, STATUS, TIMESTAMP, SENSOR_ID
		FROM SENSOR_HISTORY
		WHERE SENSOR_ID = :1
		ORDER BY TIMESTAMP DESC`,
		rawID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao consultar histórico"})
		return
	}
	defer rows.Close()

	var history []models.SensorHistory
	for rows.Next() {
		var h models.SensorHistory
		var rawHistID, rawSensorID []byte

		if err := rows.Scan(&rawHistID, &h.Value, &h.Status, &h.Timestamp, &rawSensorID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao mapear histórico"})
			return
		}

		h.SensorHistoryID = hex.EncodeToString(rawHistID)
		h.SensorID = hex.EncodeToString(rawSensorID)

		history = append(history, h)
	}

	c.JSON(http.StatusOK, gin.H{
		"sensor":  s,
		"history": history,
	})
}

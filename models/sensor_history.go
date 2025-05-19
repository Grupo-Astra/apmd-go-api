package models

import "time"

type SensorHistory struct {
	SensorHistoryID string    `json:"sensor_history_id"`
	Value           float64   `json:"value"`
	Status          string    `json:"status"`
	Timestamp       time.Time `json:"timestamp"`
	SensorID        string    `json:"sensor_id"`
}

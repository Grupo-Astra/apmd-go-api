package models

type Sensor struct {
	SensorID      string  `json:"sensor_id"`
	Name          string  `json:"name"`
	CurrentValue  float64 `json:"current_value"`
	CurrentStatus string  `json:"current_status"`
}

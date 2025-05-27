package models

import "time"

type SensorHistory struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Value     float64   `json:"value" gorm:"not null"`
	Status    string    `json:"status" gorm:"not null"`
	Timestamp time.Time `json:"timestamp" gorm:"autoCreateTime"`
	SensorID  int       `json:"sensorId" gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

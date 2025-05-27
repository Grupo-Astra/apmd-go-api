package models

type Sensor struct {
	ID            int             `json:"id" gorm:"primaryKey;autoIncrement"`
	Name          string          `json:"name" gorm:"not null"`
	CurrentValue  float64         `json:"currentValue" gorm:"not null"`
	CurrentStatus string          `json:"currentStatus" gorm:"not null"`
	Historic      []SensorHistory `json:"historic" gorm:"foreignKey:SensorID"`
}

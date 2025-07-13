package models

import (
	"time"
)

type Machine struct {
	ID          int       `gorm:"primaryKey;autoIncrement" json:"id"`
	MachineIP   string    `gorm:"column:machine_ip;type:varchar(255);not null" json:"machine_ip"`
	MachinePort string    `gorm:"column:machine_port;type:varchar(255);not null" json:"machine_port"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

// TableName overrides the default table name
func (Machine) TableName() string {
	return "machines"
}

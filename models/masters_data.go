package models

import "time"

type MasterData struct {
	ID              uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	MachineMasterID int       `gorm:"column:machine_master_id;not null" json:"machine_master_id"`
	MachineYear     int       `gorm:"column:machine_year;not null" json:"machine_year"`
	MachineMonth    int       `gorm:"column:machine_month;not null" json:"machine_month"`
	MachineDay      int       `gorm:"column:machine_day;not null" json:"machine_day"`
	MachineHour     int       `gorm:"column:machine_hour;not null" json:"machine_hour"`
	MachineMinute   int       `gorm:"column:machine_minute;not null" json:"machine_minute"`
	MachineSecond   int       `gorm:"column:machine_second;not null" json:"machine_second"`
	TempSet         *int      `gorm:"column:temp_set" json:"temp_set"` // nullable
	HumSet          *int      `gorm:"column:hum_set" json:"hum_set"`   // nullable
	TempAct         *int      `gorm:"column:temp_act" json:"temp_act"` // nullable
	HumAct          *int      `gorm:"column:hum_act" json:"hum_act"`   // nullable
	CreatedAt       time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (MasterData) TableName() string {
	return "masters_data"
}

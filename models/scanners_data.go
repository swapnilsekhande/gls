package models

import "time"

type ScannerData struct {
	ID            uint `gorm:"primaryKey;autoIncrement" json:"id"`
	MachineMaster int  `gorm:"column:machine_master;not null" json:"machine_master"`
	MachineYear   int  `gorm:"column:machine_year;not null" json:"machine_year"`
	MachineMonth  int  `gorm:"column:machine_month;not null" json:"machine_month"`
	MachineDay    int  `gorm:"column:machine_day;not null" json:"machine_day"`
	MachineHour   int  `gorm:"column:machine_hour;not null" json:"machine_hour"`
	MachineMinute int  `gorm:"column:machine_minute;not null" json:"machine_minute"`
	MachineSecond int  `gorm:"column:machine_second;not null" json:"machine_second"`
	TempSet       int  `gorm:"column:temp_set;not null" json:"temp_set"`
	HumSet        int  `gorm:"column:hum_set;not null" json:"hum_set"`

	// Temperature channels
	TempCh0  *int `gorm:"column:temp_ch0" json:"temp_ch0"`
	TempCh1  *int `gorm:"column:temp_ch1" json:"temp_ch1"`
	TempCh2  *int `gorm:"column:temp_ch2" json:"temp_ch2"`
	TempCh3  *int `gorm:"column:temp_ch3" json:"temp_ch3"`
	TempCh4  *int `gorm:"column:temp_ch4" json:"temp_ch4"`
	TempCh5  *int `gorm:"column:temp_ch5" json:"temp_ch5"`
	TempCh6  *int `gorm:"column:temp_ch6" json:"temp_ch6"`
	TempCh7  *int `gorm:"column:temp_ch7" json:"temp_ch7"`
	TempCh8  *int `gorm:"column:temp_ch8" json:"temp_ch8"`
	TempCh9  *int `gorm:"column:temp_ch9" json:"temp_ch9"`
	TempCh10 *int `gorm:"column:temp_ch10" json:"temp_ch10"`
	TempCh11 *int `gorm:"column:temp_ch11" json:"temp_ch11"`

	// Humidity channels
	HumCh0  *int `gorm:"column:hum_ch0" json:"hum_ch0"`
	HumCh1  *int `gorm:"column:hum_ch1" json:"hum_ch1"`
	HumCh2  *int `gorm:"column:hum_ch2" json:"hum_ch2"`
	HumCh3  *int `gorm:"column:hum_ch3" json:"hum_ch3"`
	HumCh4  *int `gorm:"column:hum_ch4" json:"hum_ch4"`
	HumCh5  *int `gorm:"column:hum_ch5" json:"hum_ch5"`
	HumCh6  *int `gorm:"column:hum_ch6" json:"hum_ch6"`
	HumCh7  *int `gorm:"column:hum_ch7" json:"hum_ch7"`
	HumCh8  *int `gorm:"column:hum_ch8" json:"hum_ch8"`
	HumCh9  *int `gorm:"column:hum_ch9" json:"hum_ch9"`
	HumCh10 *int `gorm:"column:hum_ch10" json:"hum_ch10"`
	HumCh11 *int `gorm:"column:hum_ch11" json:"hum_ch11"`

	CreatedAt *time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

// TableName overrides the default table name
func (ScannerData) TableName() string {
	return "scanners_data"
}

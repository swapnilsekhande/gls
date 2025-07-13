package databases

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DiantaDB *gorm.DB

func InitDb() error {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbname,
	)
	var err error
	DiantaDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}
	fmt.Println("✅ Connected to MySQL database!")
	sqlDB, err := DiantaDB.DB()
	if err != nil {
		log.Fatal(err)
	}
	// defer sqlDB.Close()
	if err := sqlDB.Ping(); err != nil {
		log.Fatal("❌ Ping failed:", err)
	}
	fmt.Println("✅ Ping successful!")
	return err
}

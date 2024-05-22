package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
	"user-service/config"
)

var db *gorm.DB

func InitDB() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.AppConfig.DBUser, config.AppConfig.DBPassword, config.AppConfig.DBHost, config.AppConfig.DBPort, config.AppConfig.DBName)
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	// Lấy instance db connection cấu hình connection pool
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get db instance: %v", err)
	}

	sqlDB.SetMaxOpenConns(10) // Set số lượng connection tối đa

	sqlDB.SetMaxIdleConns(10) // Set số lượng connection tối đa được sử dụng

	sqlDB.SetConnMaxLifetime(time.Hour) // Set thời gian tối đa mà một connection có thể được sử dụng

}

func GetDB() *gorm.DB {
	return db
}

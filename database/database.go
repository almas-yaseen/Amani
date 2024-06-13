package database

import (
	"fmt"
	"ginapp/config"
	"ginapp/domain"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase(cfg config.Config) (*gorm.DB, error) {

	psqlInfo := fmt.Sprintf("host=localhost user=%s dbname=%s  port=%s password=%s sslmode=disable TimeZone=Asia/Shanghai", cfg.DBUser, cfg.DBName, cfg.DBPort, cfg.DBPassword)

	db, dberr := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})

	if dberr != nil {
		return nil, fmt.Errorf("Failed to connect the database")
	}
	DB = db

	DB.AutoMigrate(&domain.Car{})
	DB.AutoMigrate(&domain.Image{})
	return DB, nil

}

package database

import (
	"fmt"
	"log"
	"valorant-app/config"
	"valorant-app/database/seeders"
	"valorant-app/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(cfg *config.Config) {
	// Build DSN (Data Source Name)
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate the schema
	err = DB.AutoMigrate(&models.User{}, &models.Team{}, &models.Role{}, &models.Permission{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Seed the database with initial data
	seeders.SeedAll(DB)

	log.Println("Database connected and migrated successfully")
}

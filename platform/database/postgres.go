package database

import (
	"fmt"
	"log"
	"os"

	"github.com/bearts/go-fiber/app/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DbInstance struct {
	Db *gorm.DB
}

var Database DbInstance

// PostgreSQLConnection func for connection to PostgreSQL database.
func PostgreSQLConnection() {
	// Define database connection settings.
	DSN := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	// Define database connection.
	db, err := gorm.Open(postgres.Open(DSN), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
		os.Exit(2)
	}
	log.Println("Connected to PostgreSQL database")
	db.Logger = logger.Default.LogMode(logger.Info)
	log.Println("Running Migration")
	// Migrate the schema
	db.AutoMigrate(&models.User{})

	Database = DbInstance{
		Db: db,
	}
}

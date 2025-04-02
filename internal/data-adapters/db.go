package data_adapters

import (
	"campaign-optimization-engine/internal/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
)

// DB instance
var DB *gorm.DB
var err error

// InitDB initializes the PostgreSQL database connection
func InitDB() {
	// Connection string to PostgreSQL database
	dsn := "host=localhost user=your_user dbname=campaign_db sslmode=disable password=your_password"
	DB, err = gorm.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}
	log.Println("Successfully connected to the database")

	// Migrate the Campaign model to the database
	DB.AutoMigrate(&models.Campaign{})
}

// CloseDB closes the database connection
func CloseDB() {
	DB.Close()
}

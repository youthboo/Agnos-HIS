package database

import (
	"fmt"
	"log"
	"HIS-api/config"
	"HIS-api/models"
)

func MigrateDB() {
	if config.DB == nil {
		log.Fatal("Database connection is not initialized")
	}

	err := config.DB.AutoMigrate(&models.Patient{}, &models.Staff{})
	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	fmt.Println("Database migrated successfully.")
}

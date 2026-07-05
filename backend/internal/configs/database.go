package configs

import (
	"fmt"
	"log"

	"github.com/Kar-Su/uas-mobile.git/internal/package/env"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetUpDatabaseConnection() *gorm.DB {
	dbUser := env.GetWithDefault[string]("DB_USERNAME", "admin")
	dbPass := env.GetWithDefault[string]("DB_PASSWORD", "admin123")
	dbHost := env.GetWithDefault[string]("DB_HOST", "db")
	dbName := env.GetWithDefault[string]("MAIN_DB", "inventory")
	dbPort := env.GetWithDefault[int]("DB_PORT", 5432)

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		dbHost, dbUser, dbPass, dbName, dbPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: SetupLogger(),
	})

	if err != nil {
		log.Printf("Error Init Database: %s", dsn)
		panic(err)
	}

	return db
}

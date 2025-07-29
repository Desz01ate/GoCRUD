package main

import (
	"arise_tech_assessment/internal/infrastructure"
	"log"
	"os"
)

func main() {
	dsn := os.Getenv("CONNECTION_STRINGS_DEFAULT")
	if dsn == "" {
		log.Fatal("CONNECTION_STRINGS_DEFAULT environment variable is required")
	}

	initializer := infrastructure.CreateDbInitializer(dsn)

	log.Println("Running database migrations...")
	if err := initializer.Init(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	log.Println("Seeding database...")
	if err := initializer.Seed(); err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}

	log.Println("Database seeding completed successfully!")
}
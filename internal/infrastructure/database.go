package infrastructure

import (
	"arise_tech_assessment/internal/domain"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DatabaseInitializer struct {
	*gorm.DB
}

func CreateDbInitializer(connectionString string) DatabaseInitializer {
	return DatabaseInitializer{
		DB: NewGormDB(connectionString),
	}
}

func NewGormDB(connectionString string) *gorm.DB {
	dsn := os.Getenv("CONNECTION_STRINGS_DEFAULT")
	if dsn == "" {
		log.Fatal("Fatal: CONNECTION_STRINGS_DEFAULT environment variable is not set.")
	}

	parsedDSN, err := url.Parse(dsn)
	if err != nil {
		log.Fatalf("Fatal: Invalid CONNECTION_STRINGS_DEFAULT DSN format: %v", err)
	}

	dbName := parsedDSN.Path
	if len(dbName) > 0 && dbName[0] == '/' {
		dbName = dbName[1:]
	}
	if dbName == "" {
		log.Fatalf("Fatal: Database name not found in CONNECTION_STRINGS_DEFAULT DSN path.")
	}

	user := parsedDSN.User.Username()
	password, _ := parsedDSN.User.Password()
	host := parsedDSN.Hostname()
	port := parsedDSN.Port()
	sslmode := parsedDSN.Query().Get("sslmode")

	if port == "" {
		port = "5432"
	}

	rootDSN := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=postgres sslmode=%s",
		host, port, user, password, sslmode)

	rootDB, err := sql.Open("pgx", rootDSN)
	if err != nil {
		log.Fatalf("Fatal: Failed to open connection to default 'postgres' database: %v", err)
	}
	defer func() {
		if closeErr := rootDB.Close(); closeErr != nil {
			log.Printf("Warning: Failed to close root database connection: %v", closeErr)
		}
	}()

	rootDB.SetMaxIdleConns(10)
	rootDB.SetMaxOpenConns(100)
	rootDB.SetConnMaxLifetime(time.Hour)

	if err = rootDB.Ping(); err != nil {
		log.Fatalf("Fatal: Failed to ping default 'postgres' database. Check credentials, host/port, or server status: %v", err)
	}
	log.Println("Successfully connected to PostgreSQL server (default database).")

	createDBSQL := fmt.Sprintf("CREATE DATABASE %s", dbName)
	_, err = rootDB.Exec(createDBSQL)
	if err != nil {
		if strings.Contains(err.Error(), fmt.Sprintf(`database "%s" already exists`, dbName)) {
			log.Printf("Database '%s' already exists. Skipping creation.", dbName)
		} else {
			log.Fatalf("Fatal: Failed to create database '%s': %v", dbName, err)
		}
	} else {
		log.Printf("Database '%s' created successfully.", dbName)
	}

	log.Printf("Attempting to connect GORM to application database '%s'...", dbName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Fatal: Failed to connect GORM to application database '%s': %v", dbName, err)
	}
	log.Println("Successfully connected GORM to application database.")

	return db
}

func (initializer *DatabaseInitializer) Init() error {
	err := initializer.DB.AutoMigrate(&domain.Account{}, &domain.Transaction{})

	if err != nil {
		return errors.New("Failed to run auto migration.")
	}

	return nil
}

func (initializer *DatabaseInitializer) Seed() error {
	seeder := NewSeeder(initializer.DB)
	return seeder.SeedData()
}

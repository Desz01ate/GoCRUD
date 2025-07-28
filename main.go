// @title Account API
// @version 1.0
// @description A demo account and transaction management system
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https

package main

import (
	"fmt"
	"log"
	"os"

	"arise_tech_assetment/internal/application"
	"arise_tech_assetment/internal/infrastructure"
	"arise_tech_assetment/internal/infrastructure/router"

	_ "arise_tech_assetment/docs"
)

func main() {
	dsn := os.Getenv("CONNECTION_STRINGS_DEFAULT")
	initializer := infrastructure.CreateDbInitializer(dsn)

	if err := initializer.Init(); err != nil {
		panic(fmt.Errorf("failed to initialize database: %w", err))
	}

	application.RegisterHandlers(initializer.DB)

	r := router.New()

	router.SetupRoutes(r)

	log.Println("Starting server on :8080")

	if err := r.Start(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

package main

import (
	"golang-auth-api-mysql/config"
	"golang-auth-api-mysql/database"
	_ "golang-auth-api-mysql/docs"
	"golang-auth-api-mysql/internal/routes"
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Golang Auth API
// @version 1.0
// @description Complete REST API with JWT authentication, role-based access control, audit logging, and product management
// @contact.name API Support
// @contact.email support@example.com
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	// Load config
	cfg := config.LoadConfig()

	// Connect to database
	db := database.ConnectMySQL(cfg)

	// Run migrations
	database.Migrate(db)

	// Seed data (optional, untuk development)
	if cfg.Environment == "development" {
		database.Seed(db)
	}

	// Setup Gin
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Setup routes
	routes.SetupRoutes(router, db, cfg)

	// Swagger endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Printf("Server running on port %s", cfg.ServerPort)
	log.Printf("Swagger UI available at http://localhost:%s/swagger/index.html", cfg.ServerPort)

	// Start server
	log.Printf("Server running on port %s", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

package database

import (
	"golang-auth-api-mysql/internal/domain"
	"golang-auth-api-mysql/pkg/hash"
	"log"

	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {
	// Check if admin already exists
	var count int64
	db.Model(&domain.User{}).Where("role = ?", "admin").Count(&count)

	if count > 0 {
		log.Println("Seed data already exists, skipping...")
		return
	}

	// Create admin user
	hashedPassword, _ := hash.HashPassword("admin123")
	admin := domain.User{
		Name:     "Admin User",
		Email:    "admin@example.com",
		Password: hashedPassword,
		Role:     "admin",
	}

	if err := db.Create(&admin).Error; err != nil {
		log.Println("Failed to seed admin:", err)
		return
	}

	// Create regular user
	hashedPassword, _ = hash.HashPassword("user123")
	user := domain.User{
		Name:     "Regular User",
		Email:    "user@example.com",
		Password: hashedPassword,
		Role:     "user",
	}

	if err := db.Create(&user).Error; err != nil {
		log.Println("Failed to seed user:", err)
		return
	}

	// Create sample products
	products := []domain.Product{
		{
			Name:        "Product 1",
			Description: "Description for product 1",
			Price:       10000,
			Stock:       100,
		},
		{
			Name:        "Product 2",
			Description: "Description for product 2",
			Price:       25000,
			Stock:       50,
		},
	}

	for _, p := range products {
		if err := db.Create(&p).Error; err != nil {
			log.Println("Failed to seed product:", err)
		}
	}

	log.Println("Seed data created successfully")
}

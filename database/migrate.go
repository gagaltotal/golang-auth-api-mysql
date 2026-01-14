package database

import (
	"golang-auth-api-mysql/internal/domain"
	"log"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&domain.User{},
		&domain.Product{},
		&domain.RefreshToken{},
		&domain.AuditLog{},
	)

	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	log.Println("Database migration completed")
}

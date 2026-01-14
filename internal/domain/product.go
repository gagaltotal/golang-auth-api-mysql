package domain

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:200;not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Price       float64        `gorm:"not null" json:"price"`
	Stock       int            `gorm:"default:0" json:"stock"`
	ImageURL    string         `gorm:"size:500" json:"image_url"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required,min=3"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Stock       int     `json:"stock" binding:"required,gte=0"`
}

type UpdateProductRequest struct {
	Name        string  `json:"name" binding:"omitempty,min=3"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"omitempty,gt=0"`
	Stock       int     `json:"stock" binding:"omitempty,gte=0"`
}

type ProductQuery struct {
	Page     int     `form:"page,default=1"`
	Limit    int     `form:"limit,default=10"`
	Search   string  `form:"search"`
	MinPrice float64 `form:"min_price"`
	MaxPrice float64 `form:"max_price"`
}

type ProductListResponse struct {
	Data       []Product `json:"data"`
	Page       int       `json:"page"`
	Limit      int       `json:"limit"`
	TotalItems int64     `json:"total_items"`
	TotalPages int       `json:"total_pages"`
}

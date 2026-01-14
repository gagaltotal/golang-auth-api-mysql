package repository

import "golang-auth-api-mysql/internal/domain"

type ProductRepository interface {
	Create(product *domain.Product) error
	FindByID(id uint) (*domain.Product, error)
	Update(product *domain.Product) error
	Delete(id uint) error
	FindAll(query domain.ProductQuery) ([]domain.Product, int64, error)
}

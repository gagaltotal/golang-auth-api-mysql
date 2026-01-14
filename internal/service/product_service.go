package service

import "golang-auth-api-mysql/internal/domain"

type ProductService interface {
	Create(req domain.CreateProductRequest) (*domain.Product, error)
	GetByID(id uint) (*domain.Product, error)
	Update(id uint, req domain.UpdateProductRequest) (*domain.Product, error)
	Delete(id uint) error
	GetAll(query domain.ProductQuery) (*domain.ProductListResponse, error)
	UpdateImage(id uint, imageURL string) (*domain.Product, error)
}

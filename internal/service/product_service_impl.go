package service

import (
	"errors"
	"golang-auth-api-mysql/internal/domain"
	"golang-auth-api-mysql/internal/repository"
	"math"

	"gorm.io/gorm"
)

type productServiceImpl struct {
	productRepo repository.ProductRepository
}

func NewProductService(productRepo repository.ProductRepository) ProductService {
	return &productServiceImpl{productRepo: productRepo}
}

func (s *productServiceImpl) Create(req domain.CreateProductRequest) (*domain.Product, error) {
	product := &domain.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
	}

	if err := s.productRepo.Create(product); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *productServiceImpl) GetByID(id uint) (*domain.Product, error) {
	return s.productRepo.FindByID(id)
}

func (s *productServiceImpl) Update(id uint, req domain.UpdateProductRequest) (*domain.Product, error) {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Price > 0 {
		product.Price = req.Price
	}
	if req.Stock >= 0 {
		product.Stock = req.Stock
	}

	if err := s.productRepo.Update(product); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *productServiceImpl) Delete(id uint) error {
	_, err := s.productRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("product not found")
		}
		return err
	}
	return s.productRepo.Delete(id)
}

func (s *productServiceImpl) GetAll(query domain.ProductQuery) (*domain.ProductListResponse, error) {
	if query.Page < 1 {
		query.Page = 1
	}
	if query.Limit < 1 {
		query.Limit = 10
	}

	products, total, err := s.productRepo.FindAll(query)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.Limit)))

	return &domain.ProductListResponse{
		Data:       products,
		Page:       query.Page,
		Limit:      query.Limit,
		TotalItems: total,
		TotalPages: totalPages,
	}, nil
}

func (s *productServiceImpl) UpdateImage(id uint, imageURL string) (*domain.Product, error) {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	product.ImageURL = imageURL
	if err := s.productRepo.Update(product); err != nil {
		return nil, err
	}

	return product, nil
}

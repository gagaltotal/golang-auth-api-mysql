package tests

import (
	"golang-auth-api-mysql/internal/domain"
	"golang-auth-api-mysql/internal/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock ProductRepository
type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) Create(product *domain.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductRepository) FindByID(id uint) (*domain.Product, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}

func (m *MockProductRepository) Update(product *domain.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockProductRepository) FindAll(query domain.ProductQuery) ([]domain.Product, int64, error) {
	args := m.Called(query)
	return args.Get(0).([]domain.Product), args.Get(1).(int64), args.Error(2)
}

func TestCreateProduct_Success(t *testing.T) {
	mockProductRepo := new(MockProductRepository)
	productService := service.NewProductService(mockProductRepo)

	req := domain.CreateProductRequest{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       10000,
		Stock:       100,
	}

	mockProductRepo.On("Create", mock.AnythingOfType("*domain.Product")).Return(nil)

	result, err := productService.Create(req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, req.Name, result.Name)
	assert.Equal(t, req.Price, result.Price)
	assert.Equal(t, req.Stock, result.Stock)

	mockProductRepo.AssertExpectations(t)
}

func TestGetProductByID_Success(t *testing.T) {
	mockProductRepo := new(MockProductRepository)
	productService := service.NewProductService(mockProductRepo)

	product := &domain.Product{
		ID:    1,
		Name:  "Test Product",
		Price: 10000,
		Stock: 100,
	}

	mockProductRepo.On("FindByID", product.ID).Return(product, nil)

	result, err := productService.GetByID(product.ID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, product.ID, result.ID)
	assert.Equal(t, product.Name, result.Name)

	mockProductRepo.AssertExpectations(t)
}

func TestUpdateProduct_Success(t *testing.T) {
	mockProductRepo := new(MockProductRepository)
	productService := service.NewProductService(mockProductRepo)

	product := &domain.Product{
		ID:    1,
		Name:  "Old Product",
		Price: 10000,
		Stock: 100,
	}

	req := domain.UpdateProductRequest{
		Name:  "Updated Product",
		Price: 15000,
		Stock: 150,
	}

	mockProductRepo.On("FindByID", product.ID).Return(product, nil)
	mockProductRepo.On("Update", mock.AnythingOfType("*domain.Product")).Return(nil)

	result, err := productService.Update(product.ID, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, req.Name, result.Name)
	assert.Equal(t, req.Price, result.Price)
	assert.Equal(t, req.Stock, result.Stock)

	mockProductRepo.AssertExpectations(t)
}

func TestDeleteProduct_Success(t *testing.T) {
	mockProductRepo := new(MockProductRepository)
	productService := service.NewProductService(mockProductRepo)

	productID := uint(1)
	product := &domain.Product{
		ID:   productID,
		Name: "Test Product",
	}

	mockProductRepo.On("FindByID", productID).Return(product, nil)
	mockProductRepo.On("Delete", productID).Return(nil)

	err := productService.Delete(productID)

	assert.NoError(t, err)
	mockProductRepo.AssertExpectations(t)
}

func TestGetAllProducts_Success(t *testing.T) {
	mockProductRepo := new(MockProductRepository)
	productService := service.NewProductService(mockProductRepo)

	query := domain.ProductQuery{
		Page:  1,
		Limit: 10,
	}

	products := []domain.Product{
		{ID: 1, Name: "Product 1", Price: 10000, Stock: 100},
		{ID: 2, Name: "Product 2", Price: 20000, Stock: 50},
	}

	mockProductRepo.On("FindAll", query).Return(products, int64(2), nil)

	result, err := productService.GetAll(query)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, len(result.Data))
	assert.Equal(t, int64(2), result.TotalItems)
	assert.Equal(t, 1, result.TotalPages)

	mockProductRepo.AssertExpectations(t)
}

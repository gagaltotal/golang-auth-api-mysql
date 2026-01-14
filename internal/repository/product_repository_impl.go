package repository

import (
	"golang-auth-api-mysql/internal/domain"

	"gorm.io/gorm"
)

type productRepositoryImpl struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepositoryImpl{db: db}
}

func (r *productRepositoryImpl) Create(product *domain.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepositoryImpl) FindByID(id uint) (*domain.Product, error) {
	var product domain.Product
	err := r.db.First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepositoryImpl) Update(product *domain.Product) error {
	return r.db.Save(product).Error
}

func (r *productRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&domain.Product{}, id).Error
}

func (r *productRepositoryImpl) FindAll(query domain.ProductQuery) ([]domain.Product, int64, error) {
	var products []domain.Product
	var total int64

	offset := (query.Page - 1) * query.Limit
	db := r.db.Model(&domain.Product{})

	if query.Search != "" {
		db = db.Where("name LIKE ? OR description LIKE ?", "%"+query.Search+"%", "%"+query.Search+"%")
	}
	if query.MinPrice > 0 {
		db = db.Where("price >= ?", query.MinPrice)
	}
	if query.MaxPrice > 0 {
		db = db.Where("price <= ?", query.MaxPrice)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Offset(offset).Limit(query.Limit).Find(&products).Error
	return products, total, err
}

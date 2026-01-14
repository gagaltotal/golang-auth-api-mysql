package repository

import (
	"golang-auth-api-mysql/internal/domain"

	"gorm.io/gorm"
)

type userRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepositoryImpl{db: db}
}

func (r *userRepositoryImpl) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *userRepositoryImpl) FindByID(id uint) (*domain.User, error) {
	var user domain.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepositoryImpl) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepositoryImpl) Update(user *domain.User) error {
	return r.db.Save(user).Error
}

func (r *userRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&domain.User{}, id).Error
}

func (r *userRepositoryImpl) FindAll(page, limit int) ([]domain.User, int64, error) {
	var users []domain.User
	var total int64

	offset := (page - 1) * limit

	if err := r.db.Model(&domain.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Offset(offset).Limit(limit).Find(&users).Error
	return users, total, err
}

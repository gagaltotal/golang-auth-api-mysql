package repository

import "golang-auth-api-mysql/internal/domain"

type UserRepository interface {
	Create(user *domain.User) error
	FindByID(id uint) (*domain.User, error)
	FindByEmail(email string) (*domain.User, error)
	FindByVerificationToken(token string) (*domain.User, error)
	FindByPasswordResetToken(token string) (*domain.User, error)
	Update(user *domain.User) error
	Delete(id uint) error
	FindAll(page, limit int) ([]domain.User, int64, error)
}

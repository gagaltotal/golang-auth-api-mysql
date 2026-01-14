package service

import "golang-auth-api-mysql/internal/domain"

type UserService interface {
	GetProfile(userID uint) (*domain.User, error)
	UpdateProfile(userID uint, req domain.UpdateUserRequest) (*domain.User, error)
	ChangePassword(userID uint, req domain.ChangePasswordRequest) error
	GetAllUsers(page, limit int) ([]domain.User, int64, error)
	DeleteUser(userID uint) error
}

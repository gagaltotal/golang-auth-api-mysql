package repository

import "golang-auth-api-mysql/internal/domain"

type RefreshTokenRepository interface {
	Create(token *domain.RefreshToken) error
	FindByToken(token string) (*domain.RefreshToken, error)
	RevokeByUserID(userID uint) error
	RevokeByToken(token string) error
	DeleteExpired() error
}

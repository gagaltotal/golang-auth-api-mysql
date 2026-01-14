package repository

import (
	"golang-auth-api-mysql/internal/domain"
	"time"

	"gorm.io/gorm"
)

type refreshTokenRepositoryImpl struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) RefreshTokenRepository {
	return &refreshTokenRepositoryImpl{db: db}
}

func (r *refreshTokenRepositoryImpl) Create(token *domain.RefreshToken) error {
	return r.db.Create(token).Error
}

func (r *refreshTokenRepositoryImpl) FindByToken(token string) (*domain.RefreshToken, error) {
	var refreshToken domain.RefreshToken
	err := r.db.Where("token = ? AND revoked = ? AND expires_at > ?",
		token, false, time.Now()).
		Preload("User").
		First(&refreshToken).Error

	if err != nil {
		return nil, err
	}
	return &refreshToken, nil
}

func (r *refreshTokenRepositoryImpl) RevokeByUserID(userID uint) error {
	return r.db.Model(&domain.RefreshToken{}).
		Where("user_id = ? AND revoked = ?", userID, false).
		Update("revoked", true).Error
}

func (r *refreshTokenRepositoryImpl) RevokeByToken(token string) error {
	return r.db.Model(&domain.RefreshToken{}).
		Where("token = ?", token).
		Update("revoked", true).Error
}

func (r *refreshTokenRepositoryImpl) DeleteExpired() error {
	return r.db.Where("expires_at < ?", time.Now()).
		Delete(&domain.RefreshToken{}).Error
}

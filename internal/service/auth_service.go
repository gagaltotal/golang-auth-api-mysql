package service

import "golang-auth-api-mysql/internal/domain"

type AuthService interface {
	Register(req domain.RegisterRequest) (*domain.AuthResponse, error)
	Login(req domain.LoginRequest) (*domain.AuthResponse, error)
	RefreshToken(req domain.RefreshRequest) (*domain.AuthResponse, error)
	Logout(refreshToken string) error
}

package service

import (
	"errors"
	"golang-auth-api-mysql/config"
	"golang-auth-api-mysql/internal/domain"
	"golang-auth-api-mysql/internal/repository"
	"golang-auth-api-mysql/pkg/hash"
	"golang-auth-api-mysql/pkg/jwt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type authServiceImpl struct {
	userRepo         repository.UserRepository
	refreshTokenRepo repository.RefreshTokenRepository
	cfg              *config.Config
}

func NewAuthService(
	userRepo repository.UserRepository,
	refreshTokenRepo repository.RefreshTokenRepository,
	cfg *config.Config,
) AuthService {
	return &authServiceImpl{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		cfg:              cfg,
	}
}

func (s *authServiceImpl) Register(req domain.RegisterRequest) (*domain.AuthResponse, error) {
	// Check if user exists
	_, err := s.userRepo.FindByEmail(req.Email)
	if err == nil {
		return nil, errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := hash.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     "user",
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Generate tokens
	return s.generateAuthResponse(user)
}

func (s *authServiceImpl) Login(req domain.LoginRequest) (*domain.AuthResponse, error) {
	// Find user
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	// Check password
	if !hash.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	// Revoke old refresh tokens
	s.refreshTokenRepo.RevokeByUserID(user.ID)

	// Generate tokens
	return s.generateAuthResponse(user)
}

func (s *authServiceImpl) RefreshToken(req domain.RefreshRequest) (*domain.AuthResponse, error) {
	refreshToken, err := s.refreshTokenRepo.FindByToken(req.RefreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	if refreshToken.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("refresh token expired")
	}

	// Revoke old token (rotation)
	s.refreshTokenRepo.RevokeByToken(req.RefreshToken)

	// Generate new tokens
	return s.generateAuthResponse(&refreshToken.User)
}

func (s *authServiceImpl) Logout(refreshToken string) error {
	return s.refreshTokenRepo.RevokeByToken(refreshToken)
}

func (s *authServiceImpl) generateAuthResponse(user *domain.User) (*domain.AuthResponse, error) {
	// Parse durations
	jwtDuration, _ := jwt.ParseDuration(s.cfg.JWTExpiration)
	refreshDuration, _ := jwt.ParseDuration(s.cfg.RefreshExpiration)

	// Generate access token
	accessToken, err := jwt.GenerateToken(user.ID, user.Email, user.Role, s.cfg.JWTSecret, jwtDuration)
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshTokenString := uuid.New().String()
	refreshToken := &domain.RefreshToken{
		UserID:    user.ID,
		Token:     refreshTokenString,
		ExpiresAt: time.Now().Add(refreshDuration),
	}

	if err := s.refreshTokenRepo.Create(refreshToken); err != nil {
		return nil, err
	}

	return &domain.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenString,
		User:         user,
	}, nil
}

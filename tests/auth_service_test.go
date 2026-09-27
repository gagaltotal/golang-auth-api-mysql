package tests

import (
	"golang-auth-api-mysql/config"
	"golang-auth-api-mysql/internal/domain"
	"golang-auth-api-mysql/internal/service"
	"golang-auth-api-mysql/pkg/hash"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// Mock UserRepository
type MockUserRepository struct {
	mock.Mock
}

// FindByPasswordResetToken implements [repository.UserRepository].
func (m *MockUserRepository) FindByPasswordResetToken(token string) (*domain.User, error) {
	panic("unimplemented")
}

func (m *MockUserRepository) Create(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByID(id uint) (*domain.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(email string) (*domain.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) Update(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) FindAll(page, limit int) ([]domain.User, int64, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]domain.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserRepository) FindByVerificationToken(token string) (*domain.User, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

// Mock RefreshTokenRepository
type MockRefreshTokenRepository struct {
	mock.Mock
}

func (m *MockRefreshTokenRepository) Create(token *domain.RefreshToken) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) FindByToken(token string) (*domain.RefreshToken, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.RefreshToken), args.Error(1)
}

func (m *MockRefreshTokenRepository) RevokeByUserID(userID uint) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) RevokeByToken(token string) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) DeleteExpired() error {
	args := m.Called()
	return args.Error(0)
}

func TestRegister_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRefreshTokenRepo := new(MockRefreshTokenRepository)

	cfg := &config.Config{
		JWTSecret:         "test-secret",
		JWTExpiration:     "15m",
		RefreshExpiration: "7d",
	}

	authService := service.NewAuthService(mockUserRepo, mockRefreshTokenRepo, cfg)

	req := domain.RegisterRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	// Mock expectations
	mockUserRepo.On("FindByEmail", req.Email).Return(nil, gorm.ErrRecordNotFound)
	mockUserRepo.On("Create", mock.AnythingOfType("*domain.User")).Return(nil)
	mockRefreshTokenRepo.On("Create", mock.AnythingOfType("*domain.RefreshToken")).Return(nil)

	// Execute
	response, err := authService.Register(req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.AccessToken)
	assert.NotEmpty(t, response.RefreshToken)
	assert.Equal(t, req.Name, response.User.Name)
	assert.Equal(t, req.Email, response.User.Email)
	assert.Equal(t, "user", response.User.Role)

	mockUserRepo.AssertExpectations(t)
	mockRefreshTokenRepo.AssertExpectations(t)
}

func TestRegister_EmailAlreadyExists(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRefreshTokenRepo := new(MockRefreshTokenRepository)

	cfg := &config.Config{
		JWTSecret:         "test-secret",
		JWTExpiration:     "15m",
		RefreshExpiration: "7d",
	}

	authService := service.NewAuthService(mockUserRepo, mockRefreshTokenRepo, cfg)

	req := domain.RegisterRequest{
		Name:     "Test User",
		Email:    "existing@example.com",
		Password: "password123",
	}

	existingUser := &domain.User{
		ID:    1,
		Email: req.Email,
	}

	// Mock expectations
	mockUserRepo.On("FindByEmail", req.Email).Return(existingUser, nil)

	// Execute
	response, err := authService.Register(req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, "email already registered", err.Error())

	mockUserRepo.AssertExpectations(t)
}

func TestLogin_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRefreshTokenRepo := new(MockRefreshTokenRepository)

	cfg := &config.Config{
		JWTSecret:         "test-secret",
		JWTExpiration:     "15m",
		RefreshExpiration: "7d",
	}

	authService := service.NewAuthService(mockUserRepo, mockRefreshTokenRepo, cfg)

	password := "password123"
	hashedPassword, _ := hash.HashPassword(password)

	user := &domain.User{
		ID:       1,
		Name:     "Test User",
		Email:    "test@example.com",
		Password: hashedPassword,
		Role:     "user",
	}

	req := domain.LoginRequest{
		Email:    user.Email,
		Password: password,
	}

	// Mock expectations
	mockUserRepo.On("FindByEmail", req.Email).Return(user, nil)
	mockRefreshTokenRepo.On("RevokeByUserID", user.ID).Return(nil)
	mockRefreshTokenRepo.On("Create", mock.AnythingOfType("*domain.RefreshToken")).Return(nil)

	// Execute
	response, err := authService.Login(req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.AccessToken)
	assert.NotEmpty(t, response.RefreshToken)
	assert.Equal(t, user.Email, response.User.Email)

	mockUserRepo.AssertExpectations(t)
	mockRefreshTokenRepo.AssertExpectations(t)
}

func TestLogin_InvalidCredentials(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRefreshTokenRepo := new(MockRefreshTokenRepository)

	cfg := &config.Config{
		JWTSecret:         "test-secret",
		JWTExpiration:     "15m",
		RefreshExpiration: "7d",
	}

	authService := service.NewAuthService(mockUserRepo, mockRefreshTokenRepo, cfg)

	req := domain.LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	// Mock expectations - user not found
	mockUserRepo.On("FindByEmail", req.Email).Return(nil, gorm.ErrRecordNotFound)

	// Execute
	response, err := authService.Login(req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, "invalid credentials", err.Error())

	mockUserRepo.AssertExpectations(t)
}

func TestRefreshToken_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRefreshTokenRepo := new(MockRefreshTokenRepository)

	cfg := &config.Config{
		JWTSecret:         "test-secret",
		JWTExpiration:     "15m",
		RefreshExpiration: "7d",
	}

	authService := service.NewAuthService(mockUserRepo, mockRefreshTokenRepo, cfg)

	user := &domain.User{
		ID:    1,
		Name:  "Test User",
		Email: "test@example.com",
		Role:  "user",
	}

	refreshToken := &domain.RefreshToken{
		ID:        1,
		UserID:    user.ID,
		Token:     "valid-refresh-token",
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		Revoked:   false,
		User:      *user,
	}

	req := domain.RefreshRequest{
		RefreshToken: refreshToken.Token,
	}

	// Mock expectations
	mockRefreshTokenRepo.On("FindByToken", req.RefreshToken).Return(refreshToken, nil)
	mockRefreshTokenRepo.On("RevokeByToken", req.RefreshToken).Return(nil)
	mockRefreshTokenRepo.On("Create", mock.AnythingOfType("*domain.RefreshToken")).Return(nil)

	// Execute
	response, err := authService.RefreshToken(req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.AccessToken)
	assert.NotEmpty(t, response.RefreshToken)

	mockRefreshTokenRepo.AssertExpectations(t)
}

func TestLogout_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockRefreshTokenRepo := new(MockRefreshTokenRepository)

	cfg := &config.Config{
		JWTSecret:         "test-secret",
		JWTExpiration:     "15m",
		RefreshExpiration: "7d",
	}

	authService := service.NewAuthService(mockUserRepo, mockRefreshTokenRepo, cfg)

	refreshToken := "valid-refresh-token"

	// Mock expectations
	mockRefreshTokenRepo.On("RevokeByToken", refreshToken).Return(nil)

	// Execute
	err := authService.Logout(refreshToken)

	// Assert
	assert.NoError(t, err)
	mockRefreshTokenRepo.AssertExpectations(t)
}

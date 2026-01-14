package tests

import (
	"golang-auth-api-mysql/internal/domain"
	"golang-auth-api-mysql/internal/service"
	"golang-auth-api-mysql/pkg/hash"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestGetProfile_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	userService := service.NewUserService(mockUserRepo)

	user := &domain.User{
		ID:    1,
		Name:  "Test User",
		Email: "test@example.com",
		Role:  "user",
	}

	mockUserRepo.On("FindByID", user.ID).Return(user, nil)

	result, err := userService.GetProfile(user.ID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, user.ID, result.ID)
	assert.Equal(t, user.Email, result.Email)

	mockUserRepo.AssertExpectations(t)
}

func TestUpdateProfile_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	userService := service.NewUserService(mockUserRepo)

	user := &domain.User{
		ID:    1,
		Name:  "Old Name",
		Email: "old@example.com",
		Role:  "user",
	}

	req := domain.UpdateUserRequest{
		Name:  "New Name",
		Email: "new@example.com",
	}

	mockUserRepo.On("FindByID", user.ID).Return(user, nil)
	mockUserRepo.On("FindByEmail", req.Email).Return(nil, gorm.ErrRecordNotFound)
	mockUserRepo.On("Update", mock.AnythingOfType("*domain.User")).Return(nil)

	result, err := userService.UpdateProfile(user.ID, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, req.Name, result.Name)
	assert.Equal(t, req.Email, result.Email)

	mockUserRepo.AssertExpectations(t)
}

func TestChangePassword_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	userService := service.NewUserService(mockUserRepo)

	oldPassword := "oldpassword123"
	hashedPassword, _ := hash.HashPassword(oldPassword)

	user := &domain.User{
		ID:       1,
		Password: hashedPassword,
	}

	req := domain.ChangePasswordRequest{
		OldPassword: oldPassword,
		NewPassword: "newpassword123",
	}

	mockUserRepo.On("FindByID", user.ID).Return(user, nil)
	mockUserRepo.On("Update", mock.AnythingOfType("*domain.User")).Return(nil)

	err := userService.ChangePassword(user.ID, req)

	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
}

func TestChangePassword_InvalidOldPassword(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	userService := service.NewUserService(mockUserRepo)

	oldPassword := "oldpassword123"
	hashedPassword, _ := hash.HashPassword(oldPassword)

	user := &domain.User{
		ID:       1,
		Password: hashedPassword,
	}

	req := domain.ChangePasswordRequest{
		OldPassword: "wrongpassword",
		NewPassword: "newpassword123",
	}

	mockUserRepo.On("FindByID", user.ID).Return(user, nil)

	err := userService.ChangePassword(user.ID, req)

	assert.Error(t, err)
	assert.Equal(t, "invalid old password", err.Error())

	mockUserRepo.AssertExpectations(t)
}

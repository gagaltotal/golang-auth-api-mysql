package service

import (
	"errors"
	"golang-auth-api-mysql/internal/domain"
	"golang-auth-api-mysql/internal/repository"
	"golang-auth-api-mysql/pkg/hash"

	"gorm.io/gorm"
)

type userServiceImpl struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userServiceImpl{userRepo: userRepo}
}

func (s *userServiceImpl) GetProfile(userID uint) (*domain.User, error) {
	return s.userRepo.FindByID(userID)
}

func (s *userServiceImpl) UpdateProfile(userID uint, req domain.UpdateUserRequest) (*domain.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		user.Name = req.Name
	}

	if req.Email != "" {
		existingUser, _ := s.userRepo.FindByEmail(req.Email)
		if existingUser != nil && existingUser.ID != userID {
			return nil, errors.New("email already in use")
		}
		user.Email = req.Email
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userServiceImpl) ChangePassword(userID uint, req domain.ChangePasswordRequest) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	if !hash.CheckPassword(req.OldPassword, user.Password) {
		return errors.New("invalid old password")
	}

	hashedPassword, err := hash.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	return s.userRepo.Update(user)
}

func (s *userServiceImpl) GetAllUsers(page, limit int) ([]domain.User, int64, error) {
	return s.userRepo.FindAll(page, limit)
}

func (s *userServiceImpl) DeleteUser(userID uint) error {
	_, err := s.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}
	return s.userRepo.Delete(userID)
}

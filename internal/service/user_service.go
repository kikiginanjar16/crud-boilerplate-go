package service

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/yourname/go-fiber-gorm-todo-auth-swagger/internal/models"
	"github.com/yourname/go-fiber-gorm-todo-auth-swagger/internal/repository"
)

type UserService interface {
	GetByID(id uint) (*models.User, error)
	UpdateProfile(id uint, name string) (*models.User, error)
	ChangePassword(id uint, oldPwd, newPwd string) error
	UpdateAvatarURL(id uint, url string) (*models.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(r repository.UserRepository) UserService {
	return &userService{repo: r}
}

func (s *userService) GetByID(id uint) (*models.User, error) {
	return s.repo.FindByID(id)
}

func (s *userService) UpdateProfile(id uint, name string) (*models.User, error) {
	u, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if name != "" {
		u.Name = name
	}
	if err := s.repo.Update(u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *userService) ChangePassword(id uint, oldPwd, newPwd string) error {
	u, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(oldPwd)) != nil {
		return errors.New("old password mismatch")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newPwd), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.repo.SetPassword(id, string(hash))
}

func (s *userService) UpdateAvatarURL(id uint, url string) (*models.User, error) {
	if err := s.repo.SetAvatar(id, url); err != nil {
		return nil, err
	}
	return s.repo.FindByID(id)
}

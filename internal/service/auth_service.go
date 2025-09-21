package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/yourname/go-fiber-gorm-todo-auth-swagger/internal/config"
	"github.com/yourname/go-fiber-gorm-todo-auth-swagger/internal/models"
	"github.com/yourname/go-fiber-gorm-todo-auth-swagger/internal/repository"
)

type AuthService interface {
	Register(name, email, password string, role models.Role) (*models.User, error)
	Login(email, password string) (string, *models.User, error)
}

type authService struct {
	cfg  *config.Config
	repo repository.UserRepository
}

func NewAuthService(cfg *config.Config, r repository.UserRepository) AuthService {
	return &authService{cfg: cfg, repo: r}
}

func (s *authService) Register(name, email, password string, role models.Role) (*models.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &models.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(hash),
		Role:         role,
	}
	if err := s.repo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *authService) Login(email, password string) (string, *models.User, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return "", nil, errors.New("invalid email or password")
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return "", nil, errors.New("invalid email or password")
	}

	claims := jwt.MapClaims{
		"sub": user.ID,
		"email": user.Email,
		"role": user.Role,
		"exp": time.Now().Add(time.Duration(s.cfg.JWTExpireMinute) * time.Minute).Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return "", nil, err
	}
	return signed, user, nil
}

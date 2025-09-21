package repository

import (
	"gorm.io/gorm"
	"github.com/yourname/go-fiber-gorm-todo-auth-swagger/internal/models"
)

type UserRepository interface {
	FindByEmail(email string) (*models.User, error)
	FindByID(id uint) (*models.User, error)
	Create(u *models.User) error
	Update(u *models.User) error
	SetPassword(id uint, hash string) error
	SetAvatar(id uint, url string) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var u models.User
	if err := r.db.Where("email = ?", email).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) FindByID(id uint) (*models.User, error) {
	var u models.User
	if err := r.db.First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) Create(u *models.User) error {
	return r.db.Create(u).Error
}

func (r *userRepository) Update(u *models.User) error {
	return r.db.Save(u).Error
}

func (r *userRepository) SetPassword(id uint, hash string) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).Update("password_hash", hash).Error
}

func (r *userRepository) SetAvatar(id uint, url string) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).Update("avatar_url", url).Error
}

package service

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/yourname/go-fiber-gorm-todo-auth-swagger/internal/models"
	"github.com/yourname/go-fiber-gorm-todo-auth-swagger/internal/repository"
)

type TodoService interface {
	List(limit, page int, search string, completed *bool, priority *models.Priority, sort string, ownerID *uint) ([]models.Todo, int64, error)
	Get(id uint) (*models.Todo, error)
	Create(input *models.Todo) (*models.Todo, error)
	Update(id uint, input *models.Todo) (*models.Todo, error)
	Delete(id uint) error
	ToggleComplete(id uint, completed bool) (*models.Todo, error)
}

type todoService struct {
	repo      repository.TodoRepository
	validator *validator.Validate
}

func NewTodoService(r repository.TodoRepository) TodoService {
	return &todoService{repo: r, validator: validator.New()}
}

func (s *todoService) List(limit, page int, search string, completed *bool, priority *models.Priority, sort string, ownerID *uint) ([]models.Todo, int64, error) {
	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit
	return s.repo.FindAll(limit, offset, search, completed, priority, sort, ownerID)
}

func (s *todoService) Get(id uint) (*models.Todo, error) {
	return s.repo.FindByID(id)
}

func (s *todoService) Create(input *models.Todo) (*models.Todo, error) {
	if input.Priority == "" {
		input.Priority = models.PriorityMedium
	}
	if input.Description != "" && len(input.Description) > 2000 {
		input.Description = input.Description[:2000]
	}
	if input.DueDate != nil && input.DueDate.Before(time.Now().AddDate(-10, 0, 0)) {
		d := time.Now()
		input.DueDate = &d
	}
	if err := s.validator.Struct(input); err != nil {
		return nil, err
	}
	if err := s.repo.Create(input); err != nil {
		return nil, err
	}
	return input, nil
}

func (s *todoService) Update(id uint, input *models.Todo) (*models.Todo, error) {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if input.Title != "" {
		existing.Title = input.Title
	}
	existing.Description = input.Description
	if input.Priority != "" {
		existing.Priority = input.Priority
	}
	existing.DueDate = input.DueDate
	existing.Completed = input.Completed

	if err := s.validator.Struct(existing); err != nil {
		return nil, err
	}
	if err := s.repo.Update(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *todoService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *todoService) ToggleComplete(id uint, completed bool) (*models.Todo, error) {
	return s.repo.ToggleComplete(id, completed)
}

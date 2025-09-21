package repository

import (
	"gorm.io/gorm"
	"github.com/yourname/go-fiber-gorm-todo-auth-swagger/internal/models"
)

type TodoRepository interface {
	FindAll(limit, offset int, search string, completed *bool, priority *models.Priority, sort string, ownerID *uint) ([]models.Todo, int64, error)
	FindByID(id uint) (*models.Todo, error)
	Create(todo *models.Todo) error
	Update(todo *models.Todo) error
	Delete(id uint) error
	ToggleComplete(id uint, completed bool) (*models.Todo, error)
}

type todoRepository struct {
	db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) TodoRepository {
	return &todoRepository{db: db}
}

func (r *todoRepository) FindAll(limit, offset int, search string, completed *bool, priority *models.Priority, sort string, ownerID *uint) ([]models.Todo, int64, error) {
	var todos []models.Todo
	q := r.db.Model(&models.Todo{})

	if search != "" {
		q = q.Where("title ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if completed != nil {
		q = q.Where("completed = ?", *completed)
	}
	if priority != nil {
		q = q.Where("priority = ?", *priority)
	}
	if ownerID != nil {
		q = q.Where("owner_id = ?", *ownerID)
	}

	var count int64
	if err := q.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	switch sort {
	case "due_asc":
		q = q.Order("due_date ASC NULLS LAST")
	case "due_desc":
		q = q.Order("due_date DESC NULLS LAST")
	case "created_desc":
		q = q.Order("created_at DESC")
	default:
		q = q.Order("created_at ASC")
	}

	if err := q.Limit(limit).Offset(offset).Find(&todos).Error; err != nil {
		return nil, 0, err
	}
	return todos, count, nil
}

func (r *todoRepository) FindByID(id uint) (*models.Todo, error) {
	var todo models.Todo
	if err := r.db.First(&todo, id).Error; err != nil {
		return nil, err
	}
	return &todo, nil
}

func (r *todoRepository) Create(todo *models.Todo) error {
	return r.db.Create(todo).Error
}

func (r *todoRepository) Update(todo *models.Todo) error {
	return r.db.Save(todo).Error
}

func (r *todoRepository) Delete(id uint) error {
	return r.db.Delete(&models.Todo{}, id).Error
}

func (r *todoRepository) ToggleComplete(id uint, completed bool) (*models.Todo, error) {
	todo, err := r.FindByID(id)
	if err != nil {
		return nil, err
	}
	todo.Completed = completed
	if err := r.db.Save(todo).Error; err != nil {
		return nil, err
	}
	return todo, nil
}

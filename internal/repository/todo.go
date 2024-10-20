// Package repository provides the database operations for the todo endpoint.
package repository

import (
	"github.com/fardinabir/todo-manager-app/internal/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Todo is the repository for the todo endpoint.
type Todo interface {
	Create(t *model.Todo) error
	Delete(id int) error
	Update(t *model.Todo) error
	Find(id int) (*model.Todo, error)
	FindAll(qry map[string]interface{}) ([]*model.Todo, error)
}

type todo struct {
	db *gorm.DB
}

// NewTodo returns a new instance of the todo repository.
func NewTodo(db *gorm.DB) Todo {
	return &todo{
		db: db,
	}
}

func (td *todo) Create(t *model.Todo) error {
	if err := td.db.Create(t).Error; err != nil {
		return err
	}
	return nil
}

func (td *todo) Update(t *model.Todo) error {
	if err := td.db.Save(t).Error; err != nil {
		return err
	}
	return nil
}

func (td *todo) Delete(id int) error {
	result := td.db.Where("id = ?", id).Delete(&model.Todo{})
	if result.RowsAffected == 0 {
		return model.ErrNotFound
	}
	if result.Error != nil {
		return result.Error
	}
	log.Info("Deleted todo with id: ", id)
	return nil
}

func (td *todo) Find(id int) (*model.Todo, error) {
	var todo *model.Todo
	err := td.db.Where("id = ?", id).Take(&todo).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, model.ErrNotFound
		}
		return nil, err
	}
	return todo, nil
}

func (td *todo) FindAll(qry map[string]interface{}) ([]*model.Todo, error) {
	var todos []*model.Todo
	tx := td.db

	if val, ok := qry["task"].(string); ok {
		tx = tx.Where("task LIKE ?", "%"+val+"%")
		delete(qry, "task")
	}
	if len(qry) > 0 {
		tx = tx.Where(qry)
	}
	err := tx.Order("priority desc").Order("created_at desc").Find(&todos).Error
	if err != nil {
		return nil, err
	}
	return todos, nil
}

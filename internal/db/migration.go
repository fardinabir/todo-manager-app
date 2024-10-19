package db

import (
	"github.com/fardinabir/todo-manager-app/internal/model"
	"gorm.io/gorm"
)

// Migrate runs the auto-migration for the database
func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&model.Todo{}); err != nil {
		return err
	}
	return nil
}

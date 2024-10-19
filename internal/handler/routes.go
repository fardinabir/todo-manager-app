package handler

import (
	"github.com/fardinabir/todo-manager-app/internal/repository"
	"github.com/fardinabir/todo-manager-app/internal/service"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// Register registers the routes for the application.
func Register(e *echo.Echo, db *gorm.DB) {
	e.Validator = NewCustomValidator()

	api := e.Group("/api/v1")

	// Health check
	healthHandler := NewHealth()
	api.GET("/healthz", healthHandler.Healthz)

	// Todo
	repository := repository.NewTodo(db)
	service := service.NewTodo(repository)
	todoHandler := NewTodo(service)
	todo := api.Group("/todos")
	{
		todo.POST("", todoHandler.Create)
		todo.GET("", todoHandler.FindAll)
		todo.GET("/:id", todoHandler.Find)
		todo.PUT("/:id", todoHandler.Update)
		todo.DELETE("/:id", todoHandler.Delete)
	}
}

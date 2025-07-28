package usecases

import (
	domain "task_management/Domain"

	"github.com/gin-gonic/gin"
)

// user related interfaces
type IUserRepository interface {
	CreateUser(user *domain.User) error
	FindByUsername(username string) (*domain.User, error)
	CountByUsername(username string) (int64, error)
	CountAll() (int64, error)
	PromoteUser(userID string) error
}

type IPasswordService interface {
	HashPassword(password string) (string, error)
	ComparePassword(hashedPassword, inputPassword string) bool
}

type IJWTService interface {
	GenerateToken(userID string, role domain.Role) (string, error)
}

// task related interfaces

type ITaskRepo interface {
	CreateTask(task *domain.Task) error
	GetAllTasks() ([]domain.Task, error)
	GetTaskByID(taskID string) (*domain.Task, error)
	UpdateTaskByID(taskID string, updatedTask *domain.Task) error
	DeleteTaskByID(taskID string) error
}
type IAuthService interface {
	AuthWithRole(roles ...string) gin.HandlerFunc
}

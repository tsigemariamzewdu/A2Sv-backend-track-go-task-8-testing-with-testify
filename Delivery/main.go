package main

import (
	"task_management/Delivery/controllers"
	"task_management/Delivery/router"
	infrastructure "task_management/infrastructure"
	repositories "task_management/Repositories"
	usecases "task_management/usecases"
	
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Gin router
	r := gin.Default()
	
	// Initialize dependencies
	userRepo := repositories.NewUserRepository()
	taskRepo := repositories.NewTaskRepository()
	passwordService := infrastructure.NewPasswordService()
	jwtService := infrastructure.NewJWTService("key")
	authService:=infrastructure.NewAuthService("key")
	
	// Create use cases
	userUseCase := usecases.NewUserUseCase(userRepo, passwordService, jwtService)
	taskUseCase := usecases.NewTaskUseCase(taskRepo)
	
	// Create controllers
	userController := controllers.NewUserController(userUseCase)
	taskController := controllers.NewTaskController(taskUseCase)
	
	// Setup routes
	if err := router.SetUpRoutes(r, userController, taskController,authService); err != nil {
		panic(err) 
	}
	
	// Start server
	r.Run(":8080")
}
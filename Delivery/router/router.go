package router

import (
	"task_management/Delivery/controllers"

	usecases "task_management/usecases"

	"github.com/gin-gonic/gin"
)

func SetUpRoutes(
	router *gin.Engine,
	userController *controllers.UserController,
	taskController *controllers.TaskController,
	authService usecases.IAuthService,
) error {
	// Public routes
	router.POST("/register", userController.Register)
	router.POST("/login", userController.Login)
	router.POST("/logout", userController.Logout)
	router.POST("/promote", userController.PromoteUser)

	userRoutes := router.Group("/tasks")
	userRoutes.Use(authService.AuthWithRole("Admin","User"))
	{
		userRoutes.GET("/", taskController.GetTasks)
		userRoutes.GET("/:id", taskController.GetTaskByID)
	}
	
	adminTaskRoutes := router.Group("/tasks")
	adminTaskRoutes.Use(authService.AuthWithRole("Admin"))
	{
		adminTaskRoutes.PUT("/:id", taskController.UpdateTaskByID)
		adminTaskRoutes.DELETE("/:id", taskController.DeleteTaskByID)
		adminTaskRoutes.POST("/", taskController.AddTask)
	}
	
	adminUserRoutes := router.Group("/admin")
	adminUserRoutes.Use(authService.AuthWithRole("Admin"))
	{
		adminUserRoutes.POST("/promote/", userController.PromoteUser)
	}
	
	return nil
}
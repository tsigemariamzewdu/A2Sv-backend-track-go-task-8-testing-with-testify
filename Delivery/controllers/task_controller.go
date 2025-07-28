package controllers

import (
	"net/http"
	// "time"

	domain "task_management/Domain"
	usecases "task_management/usecases"

	"github.com/gin-gonic/gin"
	// "go.mongodb.org/mongo-driver/bson/primitive"
)

//holds a reference to the user usecase
type UserController struct {
	UserUseCase *usecases.UserUseCase
}
type TaskController struct{
	TaskUseCase *usecases.TaskUseCase
}
type RegisterUserInputDTO struct{
	Username string  `json:"username"`
	Password string   `json:"password"`
}

//constructor

func NewUserController (uc *usecases.UserUseCase) *UserController{
	return &UserController{
		UserUseCase: uc,
	}
}

func NewTaskController ( tc *usecases.TaskUseCase)*TaskController{
	return &TaskController{
		TaskUseCase: tc,
	}
}

// Register controller
func (userctrl *UserController) Register(c *gin.Context) {
	var newUser RegisterUserInputDTO

	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input:" + err.Error()})
		return
	}

	user, err := userctrl.UserUseCase.Register(userctrl.ChangeToDomain(&newUser))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.Password = ""
	c.IndentedJSON(http.StatusOK, user)
}

// Login controller
func (userctrl *UserController) Login(c *gin.Context) {
	var input RegisterUserInputDTO

	if err := c.ShouldBindJSON(&input); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid input format"})
		return
	}

	token, user, err := userctrl.UserUseCase.Login( *userctrl.ChangeToDomain(&input))
	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("auth_token", token, 3600*24, "/", "", false, true)
	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "login successful",
		"user": gin.H{
			"id":       user.ID.Hex(),
			"username": user.Username,
			"role":     user.Role,
		},
	})
}
//Logout controller
func (userctrl *UserController) Logout(c *gin.Context) {
	c.SetCookie("auth_token", "", -1, "/", "", false, true)
	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "logged out successfully",
	})
}

// PromoteUser controller
func (userctrl *UserController) PromoteUser(c *gin.Context) {
	var req struct {
		UserID string `json:"userId"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	err := userctrl.UserUseCase.PromoteUser( req.UserID)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "user promoted to admin"})
}
//function to change the dto to domain
func (userctrl *UserController)ChangeToDomain(input *RegisterUserInputDTO)*domain.RegisterUserInput{
	var user domain.RegisterUserInput
	user.Username = input.Username
	user.Password = input.Password
	return &user
}
//get all tasks controller

func (taskctrl *TaskController)GetTasks(c *gin.Context) {
	tasks,err := taskctrl.TaskUseCase.GetAllTasks()
	if err!=nil{
		c.IndentedJSON(http.StatusInternalServerError,  gin.H{"error":err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, tasks)
}

func (taskctrl *TaskController)GetTaskByID(c *gin.Context) {
	id := c.Param("id")

	
	task, err := taskctrl.TaskUseCase.GetTaskByID(id)
	if err !=nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	c.IndentedJSON(http.StatusOK, task)
}

func (taskctrl *TaskController) AddTask(c *gin.Context) {
	var newTask domain.InputTask

	if err := c.ShouldBindJSON(&newTask); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error here ": err.Error()})
		return
	}
	tasknew,err := taskctrl.TaskUseCase.AddTask(&newTask)
	if err !=nil  {
		c.IndentedJSON(http.StatusConflict, gin.H{"error": "task already exits"})
		return
	}
	c.IndentedJSON(http.StatusCreated, tasknew)
}

//controller to delete a task
func (taskctrl *TaskController) DeleteTaskByID(c *gin.Context) {
	id := c.Param("id")
	

	
	err :=taskctrl.TaskUseCase.DeleteTaskByID(id)
	if err !=nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "task deleted"})

}

//controller to update task by id 
func (taskctrl *TaskController) UpdateTaskByID(c *gin.Context) {
	id := c.Param("id")
	
	

	var updatedTask domain.Task
	if err := c.ShouldBindJSON(&updatedTask); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := taskctrl.TaskUseCase.UpdateTaskByID(id, &updatedTask)
	if err !=nil{
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	c.IndentedJSON(http.StatusOK, updatedTask)

}




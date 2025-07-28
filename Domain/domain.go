package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// define taks and user struct

type TaskStatus string

const (
	StatusNotStarted TaskStatus = "not-started"
	StatusInProgress TaskStatus = "in-progress"
	StatusCompleted  TaskStatus = "completed"
)

type Task struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	DueDate     time.Time          `bson:"dueDate" json:"dueDate"`
	Status      TaskStatus         `bson:"status" json:"status"`
	
}
type InputTask struct{
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	DueDate     time.Time          `bson:"dueDate" json:"dueDate"`
	Status      TaskStatus         `bson:"status" json:"status"`
}
type Role string

const (
	RoleAdmin Role = "Admin"
	RoleUser  Role = "User"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username string             `bson:"username" json:"username"`
	Password string             `bson:"password,omitempty" json:"-"` 
	Role     Role               `bson:"role" json:"role"`
}
type RegisterUserInput struct{
	Username string
	Password string
}




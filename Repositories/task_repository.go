package repositories

import (
	"context"
	"errors"
	"fmt"
	domain "task_management/Domain"
	"task_management/db"
	"task_management/usecases"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	// "go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
)
type ITaskMongoCollection interface {
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) // Returns ICursor
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult // Returns ISingleResult
	CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error)
	UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
}
type TaskRepository struct {
	Collection ITaskMongoCollection
	 Context context.Context
	
}
func NewTaskRepository() usecases.ITaskRepo {
	col:=db.GetTasksCollection()
	ctx:= context.Background()
	return &TaskRepository{
		Collection: col,
		Context: ctx,
	}
}
// function to create a new task in the database
func ( r *TaskRepository) CreateTask(task *domain.Task) error{
	_,err:= r.Collection.InsertOne(r.Context,task)
	return err
}

//function to get all tasks 
func (r *TaskRepository) GetAllTasks() ([]domain.Task, error) {
    // Initialize empty slice to return empty slice in case of no tasks
    tasks := make([]domain.Task, 0)

    cur, err := r.Collection.Find(r.Context, bson.M{})
    if err != nil {
        return nil, fmt.Errorf("failed to fetch tasks: %v", err)  
    }
    defer cur.Close(r.Context)

    // Decode all 
    if err := cur.All(r.Context, &tasks); err != nil {
        return nil, fmt.Errorf("failed to decode tasks: %v", err)
    }

    return tasks, nil
}
//function to get task by id
func (r *TaskRepository) GetTaskByID( taskID string) (*domain.Task, error) {
	//check id 
	objID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return nil, errors.New("invalid task ID")
	}

	//find task mapped with that id 

	var task domain.Task
	err = r.Collection.FindOne(r.Context, bson.M{"_id": objID}).Decode(&task)
	if err != nil {
		return nil, err
	}
	return &task, nil
}
//function to update task by id 

func (r *TaskRepository) UpdateTaskByID(taskID string, updatedTask *domain.Task) error {
	objID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return errors.New("invalid task ID")
	}

	update := bson.M{
		"$set": bson.M{
			"title":       updatedTask.Title,
			"description": updatedTask.Description,
			"status":   updatedTask.Status,
		},
	}

	result, err := r.Collection.UpdateOne(r.Context, bson.M{"_id": objID}, update)
	if err != nil {
		return err
	}
	//count and if 0 it means no task
	if result.MatchedCount == 0 {
		return errors.New("task not found")
	}
	return nil
}
//function to delete task by id
func (r *TaskRepository) DeleteTaskByID( taskID string) error {
	objID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return errors.New("invalid task ID")
	}

	result, err := r.Collection.DeleteOne(r.Context, bson.M{"_id": objID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("task not found")
	}
	return nil
}
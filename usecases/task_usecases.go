package usecases

import (
	"errors"
	domain "task_management/Domain"

	// repositories "task_management/Repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// define TaskUseCase struct
type TaskUseCase struct {
	TaskRepo ITaskRepo
}

func NewTaskUseCase(repo ITaskRepo) *TaskUseCase {
	return &TaskUseCase{
		TaskRepo: repo,
	}
}

// add new task usecase
func (uc *TaskUseCase) AddTask(input *domain.InputTask) (*domain.Task, error) {

	task := &domain.Task{
		ID:          primitive.NewObjectID(),
		Title:       input.Title,
		Description: input.Description,
		Status:      input.Status,
	}

	err := uc.TaskRepo.CreateTask(task)
	if err != nil {
		return nil, errors.New("failed to create task")
	}
	return task, nil
}

// getalltasksusecase
func (uc *TaskUseCase) GetAllTasks() ([]domain.Task, error) {

	tasks, err := uc.TaskRepo.GetAllTasks()
	if err != nil {
		return nil, errors.New("failed to retrieve")
	}
	return tasks, nil

}

// get task byID use case
func (uc *TaskUseCase) GetTaskByID(id string) (*domain.Task, error) {

	task, err := uc.TaskRepo.GetTaskByID(id)
	if err != nil {
		return nil, errors.New("task not found")
	}
	return task, nil
}

// update task by id
func (uc *TaskUseCase) UpdateTaskByID(id string, input *domain.Task) error {

	return uc.TaskRepo.UpdateTaskByID(id, input)
}

// delete task by id
func (uc *TaskUseCase) DeleteTaskByID(id string) error {

	return uc.TaskRepo.DeleteTaskByID(id)
}

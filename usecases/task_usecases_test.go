package usecases_test

import (
	domain "task_management/Domain"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository)CreateTask(task *domain.Task) error{
	args:=m.Called(task)
	return args.Error(0)
}

func (m *MockTaskRepository) GetAllTasks() ([]domain.Task, error) {
    args:=m.Called()
	return args.Get(0).([]domain.Task),args.Error(1)
}

func (m *MockTaskRepository) GetTaskByID( taskID string) (*domain.Task, error) {
	args:=m.Called(taskID)
	return args.Get(0).(*domain.Task),args.Error(1)
}

func (m *MockTaskRepository) UpdateTaskByID(taskID string, updatedTask *domain.Task) error {
	args:=m.Called(taskID)
	return args.Error(0)
	
}

func (m *MockTaskRepository) DeleteTaskByID( taskID string) error {
	args:=m.Called(taskID)
	return args.Error(0)
	
}

//test suite
type TaskUsecaseTestSuite struct{
	suite.Suite
	taskRepo *MockTaskRepository
}

//setting up the test
func (suite *TaskUsecaseTestSuite) SetupTest(){
	suite.taskRepo=new(MockTaskRepository)
}

func TestTaskUseCaseSuite( t *testing.T){
	suite.Run(t,new(TaskUsecaseTestSuite))
}

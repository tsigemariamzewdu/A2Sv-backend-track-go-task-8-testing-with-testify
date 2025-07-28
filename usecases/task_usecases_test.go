package usecases_test

import (
	"errors"
	domain "task_management/Domain"
	"task_management/usecases"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	useCase *usecases.TaskUseCase
}

//setting up the test
func (suite *TaskUsecaseTestSuite) SetupTest(){
	suite.taskRepo=new(MockTaskRepository)
	suite.useCase=usecases.NewTaskUseCase(
		suite.taskRepo,
	)
}

func TestTaskUseCaseSuite( t *testing.T){
	suite.Run(t,new(TaskUsecaseTestSuite))
}

func (suite *TaskUsecaseTestSuite) TestAddTask() {
    // Setup test input
    input := &domain.InputTask{
        Title:       "Test Task",
        Description: "Test Description",
        Status:      "pending",
    }

    // Test 1 Successful task creation
    suite.Run("successful task creation", func() {
        suite.SetupTest()
        

        suite.taskRepo.On("CreateTask", mock.AnythingOfType("*domain.Task")).Return(nil).Once()

        task, err := suite.useCase.AddTask(input)
        
        suite.NoError(err)
        suite.NotNil(task)
        suite.Equal(input.Title, task.Title)
        suite.Equal(input.Description, task.Description)
        suite.Equal(input.Status, task.Status)
        suite.NotEmpty(task.ID)
        suite.taskRepo.AssertExpectations(suite.T())
    })

    // Test 2  Repository returns error
    suite.Run("repository returns error", func() {
        suite.SetupTest()
        
        expectedErr := errors.New("database error")
        suite.taskRepo.On("CreateTask", mock.AnythingOfType("*domain.Task")).Return(expectedErr).Once()

        task, err := suite.useCase.AddTask(input)
        
        suite.Error(err)
        suite.Nil(task)
        suite.EqualError(err, "failed to create task")
        suite.taskRepo.AssertExpectations(suite.T())
    })
}

func (suite *TaskUsecaseTestSuite) TestGetAllTasks() {
    // Setup test data
    mockTasks := []domain.Task{
        {
            ID:          primitive.NewObjectID(),
            Title:       "Task 1",
            Description: "Description 1",
            Status:      "pending",
        },
        {
            ID:          primitive.NewObjectID(),
            Title:       "Task 2",
            Description: "Description 2",
            Status:      "completed",
        },
    }

    // Test 1  Successful retrieval
    suite.Run("successful retrieval", func() {
        suite.SetupTest()
        
        suite.taskRepo.On("GetAllTasks").Return(mockTasks, nil).Once()

        tasks, err := suite.useCase.GetAllTasks()
        
        suite.NoError(err)
        suite.Equal(mockTasks, tasks)
        suite.taskRepo.AssertExpectations(suite.T())
    })

    // Test 2  Repository returns error
    suite.Run("repository returns error", func() {
        suite.SetupTest()
        
        expectedErr := errors.New("database error")
        suite.taskRepo.On("GetAllTasks").Return(nil, expectedErr).Once()

        tasks, err := suite.useCase.GetAllTasks()
        
        suite.Error(err)
        suite.Nil(tasks)
        suite.EqualError(err, "failed to retrieve")
        suite.taskRepo.AssertExpectations(suite.T())
    })
}

func (suite *TaskUsecaseTestSuite) TestGetTaskByID() {
    // Setup test data
    taskID := primitive.NewObjectID().Hex()
    mockTask := &domain.Task{
        ID:          primitive.NewObjectID(),
        Title:       "Test Task",
        Description: "Test Description",
        Status:      "pending",
    }

    // Test 1: Successful retrieval
    suite.Run("successful retrieval", func() {
        suite.SetupTest()
        
        suite.taskRepo.On("GetTaskByID", taskID).Return(mockTask, nil).Once()

        task, err := suite.useCase.GetTaskByID(taskID)
        
        suite.NoError(err)
        suite.Equal(mockTask, task)
        suite.taskRepo.AssertExpectations(suite.T())
    })

    // Test 2 Task not found
    suite.Run("task not found", func() {
        suite.SetupTest()
        
        expectedErr := errors.New("not found")
        suite.taskRepo.On("GetTaskByID", taskID).Return(nil, expectedErr).Once()

        task, err := suite.useCase.GetTaskByID(taskID)
        
        suite.Error(err)
        suite.Nil(task)
        suite.EqualError(err, "task not found")
        suite.taskRepo.AssertExpectations(suite.T())
    })

    // Test 3  Invalid ID format
    suite.Run("invalid ID format", func() {
        suite.SetupTest()

        
        task, err := suite.useCase.GetTaskByID("invalid-id")
        
        suite.Error(err)
        suite.Nil(task)
        suite.taskRepo.AssertNotCalled(suite.T(), "GetTaskByID")
    })
}

func (suite *TaskUsecaseTestSuite) TestUpdateTaskByID() {
    // Setup test data
    taskID := primitive.NewObjectID().Hex()
    updatedTask := &domain.Task{
        Title:       "Updated Task",
        Description: "Updated Description",
        Status:      "completed",
    }

    // Test 1  Successful update
    suite.Run("successful update", func() {
        suite.SetupTest()
        
        suite.taskRepo.On("UpdateTaskByID", taskID, updatedTask).Return(nil).Once()

        err := suite.useCase.UpdateTaskByID(taskID, updatedTask)
        
        suite.NoError(err)
        suite.taskRepo.AssertExpectations(suite.T())
    })

    // Test 2  Repository returns error
    suite.Run("repository returns error", func() {
        suite.SetupTest()
        
        expectedErr := errors.New("database error")
        suite.taskRepo.On("UpdateTaskByID", taskID, updatedTask).Return(expectedErr).Once()

        err := suite.useCase.UpdateTaskByID(taskID, updatedTask)
        
        suite.Error(err)
        suite.Equal(expectedErr, err) 
        suite.taskRepo.AssertExpectations(suite.T())
    })

    // Test 3: Invalid ID format
    suite.Run("invalid ID format", func() {
        suite.SetupTest()

        err := suite.useCase.UpdateTaskByID("invalid-id", updatedTask)
        
        suite.Error(err)
        suite.taskRepo.AssertNotCalled(suite.T(), "UpdateTaskByID")
    })
}

func (suite *TaskUsecaseTestSuite) TestDeleteTaskByID() {
    // Setup test data
    taskID := primitive.NewObjectID().Hex()

    // Test 1  Successful deletion
    suite.Run("successful deletion", func() {
        suite.SetupTest()
        
        suite.taskRepo.On("DeleteTaskByID", taskID).Return(nil).Once()

        err := suite.useCase.DeleteTaskByID(taskID)
        
        suite.NoError(err)
        suite.taskRepo.AssertExpectations(suite.T())
    })

    // Test 2  Repository returns error
    suite.Run("repository returns error", func() {
        suite.SetupTest()
        
        expectedErr := errors.New("database error")
        suite.taskRepo.On("DeleteTaskByID", taskID).Return(expectedErr).Once()

        err := suite.useCase.DeleteTaskByID(taskID)
        
        suite.Error(err)
        suite.Equal(expectedErr, err) 
        suite.taskRepo.AssertExpectations(suite.T())
    })

    // Test 3  Invalid ID format
    suite.Run("invalid ID format", func() {
        suite.SetupTest()

        err := suite.useCase.DeleteTaskByID("invalid-id")
        
        suite.Error(err)
        suite.taskRepo.AssertNotCalled(suite.T(), "DeleteTaskByID")
    })
}
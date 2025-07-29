package repositories_test

import (
	"context"
	"errors"
	"testing"
	

	domain "task_management/Domain"
	repositories "task_management/Repositories"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MockTaskCollection mocks the MongoDB  task Collection
type MockTaskCollection struct{
	mock.Mock
}

func (m *MockTaskCollection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	args := m.Called(ctx, document)
	var result *mongo.InsertOneResult
	
	if res, ok := args.Get(0).(*mongo.InsertOneResult); ok {
		result = res
	} else if args.Get(0) != nil {
		panic(errors.New("mocked InsertOne result was not of type *mongo.InsertOneResult"))
	}
	return result, args.Error(1)
}

func (m *MockTaskCollection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*mongo.Cursor), args.Error(1)
}
func (m *MockTaskCollection) CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}


func (m *MockTaskCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	args := m.Called(ctx, filter)
	return args.Get(0).(*mongo.SingleResult)
}

func (m *MockTaskCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	args := m.Called(ctx, filter, update)
	return args.Get(0).(*mongo.UpdateResult), args.Error(1)
}

func (m *MockTaskCollection) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*mongo.DeleteResult), args.Error(1)
}

// MockCursor mocks the MongoDB Cursor
type MockCursor struct {
	mock.Mock
	data  []interface{}
	index int
}

func (m *MockCursor) Next(ctx context.Context) bool {
	m.index++
	return m.index <= len(m.data)
}

func (m *MockCursor) Decode(val interface{}) error {
	if m.index-1 >= len(m.data) {
		return mongo.ErrNoDocuments
	}
	bsonBytes, _ := bson.Marshal(m.data[m.index-1])
	return bson.Unmarshal(bsonBytes, val)
}

func (m *MockCursor) Close(ctx context.Context) error {
	return nil
}

func (m *MockCursor) Err() error {
	return nil
}

func (m *MockCursor) All(ctx context.Context, results interface{}) error {
	
	return nil
}

// MockSingleResult mocks the MongoDB SingleResult
type MockSingleResult struct {
	mock.Mock
	data interface{}
}

func (m *MockSingleResult) Decode(v interface{}) error {
	bsonBytes, _ := bson.Marshal(m.data)
	return bson.Unmarshal(bsonBytes, v)
}

func (m *MockSingleResult) Err() error {
	return nil
}

type TaskRepositoryTestSuite struct {
	suite.Suite
	repo        *repositories.TaskRepository
	mockCol     *MockTaskCollection
	mockContext context.Context
}

func (suite *TaskRepositoryTestSuite) SetupTest() {
	suite.mockCol = new(MockTaskCollection)
	suite.mockContext = context.Background()
	suite.repo = &repositories.TaskRepository{
		Collection: suite.mockCol,
		Context:    suite.mockContext,
	}
}

func TestTaskRepositorySuite(t *testing.T) {
	suite.Run(t, new(TaskRepositoryTestSuite))
}

func (suite *TaskRepositoryTestSuite) TestCreateTask() {
	task := &domain.Task{
		Title:       "Test Task",
		Description: "Test Description",
		Status:      "pending",
	}

	// Test Case 1  Successful task creation
	suite.Run("Success", func() {
		suite.SetupTest()

		suite.mockCol.On("InsertOne",
			suite.mockContext,
			mock.AnythingOfType("*domain.Task"),
		).Return(&mongo.InsertOneResult{InsertedID: primitive.NewObjectID()}, nil).Run(func(args mock.Arguments) {
			insertedTask := args.Get(1).(*domain.Task)
			suite.NotNil(insertedTask.ID)
			suite.Equal(task.Title, insertedTask.Title)
			suite.Equal(task.Description, insertedTask.Description)
			suite.Equal(task.Status, insertedTask.Status)
		}).Once()

		err := suite.repo.CreateTask(task)
		suite.NoError(err)
		suite.mockCol.AssertExpectations(suite.T())
	})

	// Test Case 2 Error during task creation
	suite.Run("Error", func() {
		suite.SetupTest()

		suite.mockCol.On("InsertOne",
			suite.mockContext,
			mock.AnythingOfType("*domain.Task"),
		).Return(nil, errors.New("db insert error")).Once()

		err := suite.repo.CreateTask(task)
		suite.Error(err)
		suite.EqualError(err, "db insert error")
		suite.mockCol.AssertExpectations(suite.T())
	})
}

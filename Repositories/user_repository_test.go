package repositories_test

import (
	"context"
	"errors"
	domain "task_management/Domain"
	repositories "task_management/Repositories"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	// "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//MockCollection mocks the MongoDB Collection
type MockCollection struct {
	mock.Mock
}

func (m *MockCollection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	args := m.Called(ctx, document)
	var result *mongo.InsertOneResult
	
	if res, ok := args.Get(0).(*mongo.InsertOneResult); ok {
		result = res
	} else if args.Get(0) != nil {
		
		panic(errors.New("mocked InsertOne result was not of type *mongo.InsertOneResult"))
	}
	return result, args.Error(1)
}	

func (m *MockCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	args := m.Called(ctx, filter)
	return args.Get(0).(*mongo.SingleResult)
	// 	return res
	// }
	// if args.Get(0) == nil {
	// 	return nil
	// }
	// panic(errors.New("mocked FindOne result was not of type *mongo.SingleResult"))
}

func (m *MockCollection) CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	args := m.Called(ctx, filter, update)
	return args.Get(0).(*mongo.UpdateResult), args.Error(1)
}

//MockSingleResult mocks the mongo.SingleResult
// type MockableSingleResult struct{
// 	mongo.SingleResult
// 	mock.Mock
// }

// func (m *MockableSingleResult) Decode(v interface{})error{
// 	args:=m.Called(v)
// 	if args.Get(0)!=nil{
// 		if user,ok:=args.Get(0).(*domain.User);ok{
// 			targetUser,ok:=v.(*domain.User)
// 			if ok{
// 				*targetUser=*user
// 			}
// 		}else if rawBson,ok:=args.Get(0).(bson.Raw);ok{
// 			return bson.Unmarshal(rawBson,v)
// 		}
// 	}
// 	return args.Error(1)
// }


// Test suite
type UserRepositoryTestSuite struct {
	suite.Suite
	repo        *repositories.UserRepository
	mockCol     *MockCollection
	mockContext context.Context
}

func (suite *UserRepositoryTestSuite) SetupTest() {
	suite.mockCol = new(MockCollection)
	suite.mockContext = context.Background()
	suite.repo = &repositories.UserRepository{
		Collection: suite.mockCol,
		Context:    suite.mockContext,
	}
}

func TestUserRepositorySuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}

func (suite *UserRepositoryTestSuite) TestCreateUser() {
	user := &domain.User{
		Username: "tsige",
		Password: "hashedpassword",
		Role:     domain.RoleUser,
	}

	// test 1 Successful user creation
	suite.Run("Success", func() {
		suite.SetupTest() 

		
		suite.mockCol.On("InsertOne",
			suite.mockContext,
			mock.AnythingOfType("*domain.User"), 
		).Return(&mongo.InsertOneResult{InsertedID: primitive.NewObjectID()}, nil).Run(func(args mock.Arguments) {
			
			insertedUser := args.Get(1).(*domain.User)
			suite.NotNil(insertedUser.ID) 
			suite.Equal(user.Username, insertedUser.Username)
			suite.Equal(user.Password, insertedUser.Password)
			suite.Equal(user.Role, insertedUser.Role)
		}).Once()

		err := suite.repo.CreateUser(user)
		suite.NoError(err)
		suite.mockCol.AssertExpectations(suite.T())
	})

	//test 2 error during user creation

	suite.Run("Error", func() {
		suite.SetupTest() 

		
		suite.mockCol.On("InsertOne",
			suite.mockContext,
			mock.AnythingOfType("*domain.User"),
		).Return(nil, errors.New("db insert error")).Once()

		err := suite.repo.CreateUser(user)
		suite.Error(err)
		suite.EqualError(err, "db insert error")
		suite.mockCol.AssertExpectations(suite.T())
	})
}

// func (suite *UserRepositoryTestSuite) TestFindByUsername() {
// 	username := "existinguser"
// 	expectedUser := &domain.User{
// 		ID:       primitive.NewObjectID(),
// 		Username: username,
// 		Password: "hashedpassword",
// 		Role:     domain.RoleUser,
// 	}

// 	// Test Case 1: User found
// 	suite.Run("User_Found", func() {
// 		suite.SetupTest()

// 		mockSingleResult := new(MockableSingleResult)
		
// 		mockSingleResult.On("Decode", mock.AnythingOfType("*domain.User")).Return(expectedUser, nil).Once()

		
// 		suite.mockCol.On("FindOne",
// 			suite.mockContext,
// 			bson.M{"username": username},
// 		).Return(mockSingleResult).Once()

// 		user, err := suite.repo.FindByUsername(username)
// 		suite.NoError(err)
// 		suite.NotNil(user)
// 		suite.Equal(expectedUser.Username, user.Username)
// 		suite.Equal(expectedUser.ID, user.ID)
// 		mockSingleResult.AssertExpectations(suite.T())
// 		suite.mockCol.AssertExpectations(suite.T())
// 	})}
package usecases_test

import (
	"errors"
	domain "task_management/Domain"
	"task_management/usecases"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo"
)

// mock user repository
type MockUserRepostitoy struct {
	mock.Mock
}

// mocks the repo method createuser
func (m *MockUserRepostitoy) CreateUser(user *domain.User) error {
	args := m.Called(user) //record calls
	return args.Error(0)   //return configured error
}

// mocks findbyusername method
func (m *MockUserRepostitoy) FindByUsername(username string) (*domain.User, error) {
    args := m.Called(username)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*domain.User), args.Error(1)
}
//mocks countbyusername method

func (m *MockUserRepostitoy) CountByUsername(username string) (int64, error) {
	args := m.Called(username)
	return args.Get(0).(int64), args.Error(1)
}

//mocks countall method

func (m *MockUserRepostitoy) CountAll() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

//mocks promoteuser method

func (m *MockUserRepostitoy) PromoteUser(userID string) error {
	args := m.Called(userID)
	return args.Error(0)
}

//mock password service

type MockPasswordService struct {
	mock.Mock
}

// mocks hashpassword method
func (m *MockPasswordService) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockPasswordService) ComparePassword(hashedPassword, inputPassword string) bool {
	args := m.Called(hashedPassword, inputPassword)
	return args.Bool(0)
}

//mock jwt service

type MockJWTService struct {
	mock.Mock
}

func (m *MockJWTService) GenerateToken(userID string, role domain.Role) (string, error) {
	args := m.Called(userID, role)
	return args.String(0), args.Error(1)
}

//Test suite

type UserUseCaseTestSuite struct {
	suite.Suite
	userRepo        *MockUserRepostitoy
	passwordService *MockPasswordService
	jwtService      *MockJWTService
	useCase         *usecases.UserUseCase
}

// setting up the test
func (suite *UserUseCaseTestSuite) SetupTest() {
	suite.userRepo = new(MockUserRepostitoy)
	suite.passwordService = new(MockPasswordService)
	suite.jwtService = new(MockJWTService)
	suite.useCase = usecases.NewUserUseCase(
		suite.userRepo,
		suite.passwordService,
		suite.jwtService,
	)
}

func TestUserUseCaseSuite(t *testing.T) {
	suite.Run(t, new(UserUseCaseTestSuite))
}

// TestRegister tests the Register method of UserUseCase

func (suite *UserUseCaseTestSuite) TestRegister() {
	//define sample registeruserinput
	input := &domain.RegisterUserInput{
		Username: "tsige",
		Password: "123123123",
	}
	hashedPassword := "hashed123123123"

	//successful registration -first user becomes admin
	suite.Run("succesfull registration first user admin", func() {
		//first setup the test which resets the mocks
		suite.SetupTest()
		suite.userRepo.On("CountByUsername", input.Username).Return(int64(0), nil).Once()
		suite.userRepo.On("CountAll").Return(int64(0), nil).Once()
		suite.passwordService.On("HashPassword", input.Password).Return(hashedPassword, nil).Once()

		
		suite.userRepo.On("CreateUser", mock.AnythingOfType("*domain.User")).Return(nil).Run(func(args mock.Arguments) {
			userArg := args.Get(0).(*domain.User)
			suite.Equal(input.Username, userArg.Username)
			suite.Equal(hashedPassword, userArg.Password)
			suite.Equal(domain.RoleAdmin, userArg.Role)
		}).Once()
		//call the register method
		user, err := suite.useCase.Register(input)

		//assertions
		suite.NoError(err)
		suite.NotNil(user)
		suite.Equal(input.Username, user.Username)
		suite.Equal(hashedPassword, user.Password)
		suite.Equal(domain.RoleAdmin, user.Role)
		suite.userRepo.AssertExpectations(suite.T())
		suite.passwordService.AssertExpectations(suite.T())
	})

	//test 2 successful registration-subsequent user becomes normal user
	suite.Run("successful registration of susequent user", func() {
    suite.SetupTest()

    suite.userRepo.On("CountByUsername", input.Username).Return(int64(0), nil).Once()
    suite.userRepo.On("CountAll").Return(int64(5), nil).Once()
    suite.passwordService.On("HashPassword", input.Password).Return(hashedPassword, nil).Once()
    suite.userRepo.On("CreateUser", mock.AnythingOfType("*domain.User")).Return(nil).Run(func(args mock.Arguments) {
        userArg := args.Get(0).(*domain.User)
        suite.Equal(input.Username, userArg.Username)
        suite.Equal(hashedPassword, userArg.Password)
        suite.Equal(domain.RoleUser, userArg.Role)
    }).Once()

    user, err := suite.useCase.Register(input)

    suite.NoError(err)
    suite.NotNil(user)
    suite.Equal(input.Username, user.Username)
    suite.Equal(hashedPassword, user.Password)
    suite.Equal(domain.RoleUser, user.Role)
    suite.userRepo.AssertExpectations(suite.T())
    suite.passwordService.AssertExpectations(suite.T())
})

	//test 3 registration with existing username
	suite.Run("registration with exitsing username ", func() {

		//setup the test
		suite.SetupTest()

		suite.userRepo.On("CountByUsername", input.Username).Return(int64(1), nil).Once()

		//call the register method
		user, err := suite.useCase.Register(input)

		//assertions
		suite.Error(err)
		suite.Nil(user)
		suite.EqualError(err, "username already exists")
		suite.userRepo.AssertExpectations(suite.T())
		suite.passwordService.AssertNotCalled(suite.T(), "HashPassword")
		suite.userRepo.AssertNotCalled(suite.T(), "CreateUser")

	})
	//test 4 error during countbyusername
	suite.Run("error countbyusername", func() {
		suite.SetupTest()

		suite.userRepo.On("CountByUsername", input.Username).Return(int64(0), errors.New("error counting"))

		//call the register method
		user, err := suite.useCase.Register(input)

		//assertions
		suite.Error(err)
		suite.Nil(user)
		suite.EqualError(err, "error while checking existing user")
		suite.userRepo.AssertExpectations(suite.T())
		suite.passwordService.AssertNotCalled(suite.T(), "HashPassword")
		suite.userRepo.AssertNotCalled(suite.T(), "CreateUser")

	})

	//test 5 error during count all
	suite.Run("error during count all", func() {
		suite.SetupTest()

		suite.userRepo.On("CountByUsername", input.Username).Return(int64(0), nil).Once()
		suite.userRepo.On("CountAll").Return(int64(0), errors.New("error counting")).Once()

		//calll the register method
		user, err := suite.useCase.Register(input)

		//assertions
		suite.Error(err)
		suite.Nil(user)
		suite.EqualError(err, "error checking total users")
		suite.userRepo.AssertExpectations(suite.T())
		suite.passwordService.AssertNotCalled(suite.T(), "HashPassword")
		suite.userRepo.AssertNotCalled(suite.T(), "CreateUser")

	})
	// test 6 error duing hashing
	suite.Run("error hashing", func() {
    suite.SetupTest()

    suite.userRepo.On("CountByUsername", input.Username).Return(int64(0), nil).Once()
    suite.userRepo.On("CountAll").Return(int64(0), nil).Once()
    suite.passwordService.On("HashPassword", input.Password).Return("", errors.New("error during hashing")).Once()

    user, err := suite.useCase.Register(input)

    suite.Error(err)
    suite.Nil(user)
    suite.EqualError(err, "failed to hash password")
    suite.userRepo.AssertExpectations(suite.T())
    suite.passwordService.AssertExpectations(suite.T())
    suite.userRepo.AssertNotCalled(suite.T(), "CreateUser")
})
	//test 7 error during createuser
	suite.Run("error while creating user", func() {
		suite.SetupTest()

		suite.userRepo.On("CountByUsername", input.Username).Return(int64(0), nil).Once()
		suite.userRepo.On("CountAll").Return(int64(0), nil).Once()
		suite.passwordService.On("HashPassword", input.Password).Return(hashedPassword, nil).Once()
		suite.userRepo.On("CreateUser", mock.AnythingOfType("*domain.User")).Return(errors.New("insertion error")).Once()

		// Call the Register method

		user, err := suite.useCase.Register(input)

		//assertions
		suite.Error(err)
		suite.Nil(user)
		suite.EqualError(err, "failed to add user")
		suite.userRepo.AssertExpectations(suite.T())
		suite.passwordService.AssertExpectations(suite.T())

	})
}

// TestLogin tests the login methid of userusercase
func (suite *UserUseCaseTestSuite) TestLogin() {
	//define a sample loginUserInput for testing
	input := &domain.RegisterUserInput{
		Username: "tsige",
		Password: "123123123",
	}
	hashedPassword := "hashed123123123"
	// userID:=primitive.NewObjectID().Hex()
	expectedToken := "mockedjwttoken"

	existingUser := &domain.User{
		ID:       primitive.NewObjectID(),
		Username: input.Username,
		Password: hashedPassword,
		Role:     domain.RoleUser,
	}

	// test 1 succesfull login

	suite.Run("succesfull login", func() {
		suite.SetupTest()

		suite.userRepo.On("FindByUsername", input.Username).Return(existingUser, nil).Once()

		suite.passwordService.On("ComparePassword", hashedPassword, input.Password).Return(true).Once()

		suite.jwtService.On("GenerateToken", existingUser.ID.Hex(), existingUser.Role).Return(expectedToken, nil).Once()

		token, user, err := suite.useCase.Login(*input)

		suite.NoError(err)
		suite.NotNil(user)
		suite.Equal(expectedToken, token)
		suite.NotNil(user)
		suite.Equal(existingUser.ID, user.ID)
		suite.Equal(existingUser.Username, user.Username)
		suite.userRepo.AssertExpectations(suite.T())
		suite.passwordService.AssertExpectations(suite.T())
		suite.jwtService.AssertExpectations(suite.T())
	})

	//test 2 user not found
	//test 2 user not found
	suite.Run("user not found", func() {
		suite.SetupTest()

		// This setup is correct - returns nil user and mongo.ErrNoDocuments
		suite.userRepo.On("FindByUsername", input.Username).Return(nil, mongo.ErrNoDocuments).Once()

		token, loggedInUser, err := suite.useCase.Login(*input)

		suite.Error(err)
		suite.Empty(token)
		suite.Nil(loggedInUser) 
		suite.EqualError(err, "invalid username or password")
		suite.userRepo.AssertExpectations(suite.T())
		suite.passwordService.AssertNotCalled(suite.T(), "ComparePassword")
		suite.jwtService.AssertNotCalled(suite.T(), "GenerateToken")
	})

	//test 3 incorrect password

	suite.Run("incorrect password",func() {
		suite.SetupTest()

		
		suite.userRepo.On("FindByUsername", input.Username).Return(existingUser, nil).Once()
		
		suite.passwordService.On("ComparePassword", hashedPassword, input.Password).Return(false).Once()

		
		token, loggedInUser, err := suite.useCase.Login(*input)

		suite.Error(err)
		suite.Empty(token)
		suite.Nil(loggedInUser) 
		suite.EqualError(err, "invalid username or password")
		suite.userRepo.AssertExpectations(suite.T())
		suite.passwordService.AssertExpectations(suite.T())
		suite.jwtService.AssertNotCalled(suite.T(), "GenerateToken")

	})

}
func (suite *UserUseCaseTestSuite) TestPromoteUser() {
    validUserID := primitive.NewObjectID().Hex()

    // Test 1  Successful promotion
    suite.Run("successful promotion", func() {
        suite.SetupTest()
        

        suite.userRepo.On("PromoteUser", validUserID).Return(nil).Once()

        err := suite.useCase.PromoteUser(validUserID)
        
        suite.NoError(err)
        suite.userRepo.AssertExpectations(suite.T())
    })

    // Test 2 Repository returns error
    suite.Run("repository returns error", func() {
        suite.SetupTest()
        
        expectedErr := errors.New("database error")
        suite.userRepo.On("PromoteUser", validUserID).Return(expectedErr).Once()

        err := suite.useCase.PromoteUser(validUserID)
        
        suite.Error(err)
        suite.Equal(expectedErr, err) 
        suite.userRepo.AssertExpectations(suite.T())
    })

    // Test 3 Empty user ID
    suite.Run("empty user ID", func() {
        suite.SetupTest()

       
        suite.userRepo.On("PromoteUser", "").Return(nil).Once()

        err := suite.useCase.PromoteUser("")
        
        suite.NoError(err) 
        suite.userRepo.AssertExpectations(suite.T())
    })

    // Test 4  Invalid user ID format
    suite.Run("invalid user ID format", func() {
        suite.SetupTest()
        
        invalidID := "invalid-id"
        
        suite.userRepo.On("PromoteUser", invalidID).Return(nil).Once()

        err := suite.useCase.PromoteUser(invalidID)
        
        suite.NoError(err) 
        suite.userRepo.AssertExpectations(suite.T())
    })
}
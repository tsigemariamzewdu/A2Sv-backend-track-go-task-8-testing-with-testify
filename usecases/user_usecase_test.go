package usecases_test

import (
	
	domain "task_management/Domain"
	"task_management/usecases"
	"testing"


	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// mock user repository
type MockUserRepostitoy struct {
	mock.Mock
}
//mocks the repo method createuser
func (m *MockUserRepostitoy) CreateUser(user *domain.User) error {
	args := m.Called(user) //record calls
	return args.Error(0) //return configured error
}

//mocks findbyusername method
func (m *MockUserRepostitoy) FindByUsername(username string)(*domain.User,error){
	args:=m.Called(username)
	return args.Get(0).(*domain.User),args.Error(1)
}

//mocks countbyusername method

func (m *MockUserRepostitoy) CountByUsername(username string)(int64,error){
	args:=m.Called(username)
	return args.Get(0).(int64),args.Error(1)
}

//mocks countall method

func (m *MockUserRepostitoy) CountAll()(int64,error){
	args:=m.Called()
	return args.Get(0).(int64),args.Error(1)
}

//mocks promoteuser method

func (m *MockUserRepostitoy) PromoteUser(userID string)error{
	args:=m.Called(userID)
	return args.Error(0)
}



//mock password service

type MockPasswordService struct{
	mock.Mock
}

//mocks hashpassword method
func (m *MockPasswordService) HashPassword(password string)(string,error){
	args:=m.Called(password)
	return args.String(0),args.Error(1)
}

func (m *MockPasswordService)ComparePassword(hashedPassword,inputPassword string)bool{
	args:=m.Called(hashedPassword,inputPassword)
	return args.Bool(0)
}


//mock jwt service

type MockJWTService struct{
	mock.Mock
}

func (m *MockJWTService)GenerateToken(userID ,role string) (string ,error){
	args:=m.Called(userID,role)
	return args.String(0),args.Error(1)
}

//Test suite

type UserUseCaseTestSuite struct{
	suite.Suite
	userRepo *MockUserRepostitoy
	passwordService *MockPasswordService
	jwtService *MockJWTService
	useCase *usecases.UserUseCase
}
//setting up the test
func (suite *UserUseCaseTestSuite) SetupTest(){
	suite.userRepo=new(MockUserRepostitoy)
	suite.passwordService=new(MockPasswordService)
	suite.jwtService=new(MockJWTService)
	suite.useCase=usecases.NewUserUseCase(
		suite.userRepo,
		suite.passwordService,
		suite.jwtService,
	)
}

func TestUserUseCaseSuite(t *testing.T){
	suite.Run(t,new(UserUseCaseTestSuite))
}

// TestRegister tests the Register method of UserUseCase

func (suite *UserUseCaseTestSuite) TestRegister(){
	//define sample registeruserinput
	input := &domain.RegisterUserInput{
		Username: "tsige",
		Password:"123123123",
	}
	hashedPassword:="hashed123123123"

	//first user becomes admin
	suite.Run("succesfull registration first user admin",func() {
		//first setup the test which resets the mocks
		suite.SetupTest()
		suite.userRepo.On("CountByUsername",input.Username).Return (int64(0),nil).Once()
		suite.userRepo.On("CountAll").Return(int64(0),nil).Once()
		suite.passwordService.On("HashPassword",input.Password).Return(hashedPassword,nil).Once()

		//expect Createuser to be called with the new user(admin role) and return no error
		suite.userRepo.On("CreateUser",mock.AnythingOfType("*domain.User")).Return(nil).Run(func(args mock.Arguments){
			userArg:=args.Get(0).(*domain.User)
			suite.Equal(input.Username,userArg.Username)
			suite.Equal(hashedPassword,userArg.Password)
			suite.Equal(domain.RoleAdmin,userArg.Role)
		}).Once()
		//call the register method
		user,err:=suite.useCase.Register(input)
		
		//assertions
		suite.NoError(err)
		suite.NotNil(user)
		suite.Equal(input.Username,user.Username)
		suite.Equal(hashedPassword,user.Password)
		suite.Equal(domain.RoleAdmin,user.Role)
		suite.userRepo.AssertExpectations(suite.T())
		suite.passwordService.AssertExpectations(suite.T())
	})
}
package usecases

import (
	"errors"
	domain "task_management/Domain"

	
)

type UserUseCase struct {
	UserRepo        IUserRepository
	PasswordService IPasswordService
	JWTService      IJWTService
	
}

func NewUserUseCase(repo IUserRepository, ps IPasswordService, jw IJWTService) *UserUseCase {
	return &UserUseCase{
		UserRepo:        repo,
		PasswordService: ps,
		JWTService:      jw,
	}
}

// register use case
func (uc *UserUseCase) Register(input *domain.RegisterUserInput) (*domain.User, error) {

	count, err := uc.UserRepo.CountByUsername(input.Username)
	if err != nil {
		return nil, errors.New("error while checking existing user")
	}
	if count > 0 {
		return nil, errors.New("username already exists")
	}
	//check number of total users and if 0 make the first user and admin
	totalUsers, err := uc.UserRepo.CountAll()
	if err != nil {
		return nil, errors.New("error checking total users")

	}
	//hash the password
	hashedPassword, err := uc.PasswordService.HashPassword(input.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")

	}
	//set the role as user first then check if it the first user and if so make it admin
	role := domain.RoleUser
	if totalUsers == 0 {
		role = domain.RoleAdmin
	}
	newUser := &domain.User{
		
		Username: input.Username,
		Password: hashedPassword,
		Role:     role,
	}
	err = uc.UserRepo.CreateUser(newUser)
	if err != nil {
		return nil, errors.New("failed to add user")
	}
	return newUser, nil
}

//login use case

func (uc *UserUseCase) Login(input domain.RegisterUserInput) (string, *domain.User, error) {

	//find username
	user, err := uc.UserRepo.FindByUsername(input.Username)
	if err != nil {
		return "", nil, errors.New("invalid username or password")
	}

	//compare password
	ok := uc.PasswordService.ComparePassword(user.Password, input.Password)
	if !ok {
		return "", nil, errors.New("invalid username or password")
	}
	//generate token
	token, err := uc.JWTService.GenerateToken(user.ID.Hex(), string(user.Role))
	if err != nil {
		return "", nil, errors.New("failed to generate token")
	}

	return token, user, nil
}

// promoteuser usecase
func (uc *UserUseCase) PromoteUser(userID string) error {
	return uc.UserRepo.PromoteUser(userID)
}

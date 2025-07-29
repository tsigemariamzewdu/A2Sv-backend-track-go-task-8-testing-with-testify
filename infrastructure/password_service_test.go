package infrastruture_test

import (
    "testing"
    
    infrastruture "task_management/infrastructure"
    "github.com/stretchr/testify/suite"
)

type PasswordServiceTestSuite struct {
    suite.Suite
    service *infrastruture.PasswordService
}

func (s *PasswordServiceTestSuite) SetupTest() {
    s.service = infrastruture.NewPasswordService().(*infrastruture.PasswordService)
}

func (s *PasswordServiceTestSuite) TestHashPassword() {
    hashed, err := s.service.HashPassword("123123123")
    s.NoError(err)
    s.NotEmpty(hashed)
    s.NotEqual("123123123", hashed)
}

func (s *PasswordServiceTestSuite) TestPasswordComparison() {
   
    password := "securePassword123"
    hashed, _ := s.service.HashPassword(password)
    
    
    s.True(s.service.ComparePassword(hashed, password))
    s.False(s.service.ComparePassword(hashed, "wrongPassword"))
    s.False(s.service.ComparePassword("invalid-hash", password))
}

func (s *PasswordServiceTestSuite) TestEdgeCases() {
    s.Run("Empty password", func() {
        hashed, err := s.service.HashPassword("")
        s.NoError(err)
        s.True(s.service.ComparePassword(hashed, ""))
    })
    
   
}

func TestPasswordServiceSuite(t *testing.T) {
    suite.Run(t, new(PasswordServiceTestSuite))
}

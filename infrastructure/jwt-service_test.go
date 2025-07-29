package infrastruture_test

import (
	domain "task_management/Domain"
	infrastruture "task_management/infrastructure"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/suite"
)

type JWTServiceTestSuite struct {
	suite.Suite
	service *infrastruture.JWTService

}

func (s *JWTServiceTestSuite) TestGenerateToken(){
	token,err:=s.service.GenerateToken("1",domain.RoleUser)
	s.NoError(err)
	s.NotEmpty(token)

	//verify token contents
	parsed,err:=jwt.Parse(token,func(t *jwt.Token)(interface {},error){
		return []byte("wellwellwell"),nil
	})
	s.NoError(err)
	s.True(parsed.Valid)

	claims:=parsed.Claims.(jwt.MapClaims)
	s.Equal("1",claims["sub"])
	s.Equal(domain.RoleUser,claims["role"])
	
}

func (s *JWTServiceTestSuite) TestInvalidSecret() {
	
	token, _ := s.service.GenerateToken("1", domain.RoleUser)
	
	// Try to parse with wrong secret
	_, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte("notwell"), nil
	})
	s.Error(err)
}

func TestJWTServiceSuite(t *testing.T) {
	suite.Run(t, new(JWTServiceTestSuite))
}
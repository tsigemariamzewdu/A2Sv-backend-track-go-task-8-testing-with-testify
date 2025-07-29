package infrastruture_test

import (
	"net/http"
	"net/http/httptest"
	domain "task_management/Domain"
	infrastruture "task_management/infrastructure"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

//define a test suite struct for authmiddleware
type AuthMiddlewareTestSuite struct {
	suite.Suite
	authService *infrastruture.AuthService
	secret      string
}

func (suite *AuthMiddlewareTestSuite) setupTest(){
	suite.secret="wellwellwell"
	suite.authService = infrastruture.NewAuthService(suite.secret).(*infrastruture.AuthService)

}

//function to create token for test 
func (suite *AuthMiddlewareTestSuite) createToken(userID int, role domain.Role) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"role": role,
		"sub":  userID, 
		"exp":  time.Now().Add(time.Hour).Unix(),
	})

	tokenStr, _ := token.SignedString([]byte(suite.secret))
	return tokenStr
}
//test with valid token and role
func (suite *AuthMiddlewareTestSuite) TestAuthWithValidTokenAndRole(){
	router:=gin.New()
	router.Use(suite.authService.AuthWithRole("Admin"))

	router.GET("/protected",func(c *gin.Context){
		userID:=c.GetString("userID")
		role:=c.GetString("Role")
		c.IndentedJSON(http.StatusOK,gin.H{"userID":userID,"role":role})

	})
	token:=suite.createToken(1,"Admin")
	req:=httptest.NewRequest(http.MethodGet,"/protected",nil)
	req.AddCookie(&http.Cookie{Name:"auth_token",Value:token})
	w:=httptest.NewRecorder()

	router.ServeHTTP(w,req)
	assert.Equal(suite.T(),http.StatusOK,w.Code)
	assert.Contains(suite.T(), w.Body.String(), 1)
	assert.Contains(suite.T(), w.Body.String(), "Admin")

	
}

// Test with valid token but wrong role
func (suite *AuthMiddlewareTestSuite) TestAuthWithInvalidRole() {
	router := gin.New()
	router.Use(suite.authService.AuthWithRole("Admin")) // only admin allowed

	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	token := suite.createToken(1, "User") 
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.AddCookie(&http.Cookie{Name: "auth_token", Value: token})
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusForbidden, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "role not authorized")
}
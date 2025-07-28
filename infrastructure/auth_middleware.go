package infrastruture

import (
	"fmt"
	"net/http"
	usecases "task_management/usecases"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)
type AuthService struct{
	jwtSecret []byte
}

func NewAuthService( secret string)usecases.IAuthService{
	return &AuthService{jwtSecret: []byte(secret)}

}

var jwtSecret = []byte("wellwellwell")

func (a *AuthService)AuthWithRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		//read the token from cookie
		cookie, err := c.Request.Cookie("auth_token")

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized: no auth cookie"})
			return
		}
		tokenstr := cookie.Value

		//parse the token using jwt.Parse

		token, err := jwt.Parse(tokenstr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("error in signing method")
			}

			return jwtSecret, nil
		})

		//check error
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized: invalid token" + err.Error()})
			return
		}

		//extract claims -role and id from the token

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized: invalid token claims"})
			return
		}
		userID, ok1 := claims["sub"].(string)
		role, ok2 := claims["role"].(string)

		if !ok1 || !ok2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized: missing user info in token"})
			return
		}
		c.Set("userID", userID)
		c.Set("userRole", role)

		//check if role is allowed
		authorized := false
		for _, r := range allowedRoles {
			if role == r {
				authorized = true
				break
			}
		}
		if !authorized {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "role not authorized"})
			return
		}

		c.Next()

	}
}

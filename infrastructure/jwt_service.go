package infrastruture

import (
	usecases "task_management/usecases"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// JWTServiceImpl implements domain.JWTService
type JWTService struct {
	secretKey []byte
}

// NewJWTService returns a new instance of JWTServiceImpl with a secret key
func NewJWTService(secret string)usecases.IJWTService {
	return &JWTService{
		secretKey: []byte(secret),
	}
}


// GenerateToken creates a signed JWT token for the given user ID and role
func (j *JWTService) GenerateToken(userID, role string) (string, error) {
	claims := jwt.MapClaims{
		"sub":  userID,
		"role": role,
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}



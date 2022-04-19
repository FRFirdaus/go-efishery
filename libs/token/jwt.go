package token

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

// TokenService is interface for token utilities
type TokenService interface {
	Validate(token string) (*jwt.Token, error)
	Generate(claim jwt.Claims) (string, error)
}
type jwtServices struct {
	secretKey string
}

func (service *jwtServices) Validate(encodedToken string) (*jwt.Token, error) {
	return jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		if _, isvalid := token.Method.(*jwt.SigningMethodHMAC); !isvalid {
			return nil, fmt.Errorf("Invalid token %v", token.Header["alg"])

		}
		return []byte(service.secretKey), nil
	})

}

func (service *jwtServices) Generate(claim jwt.Claims) (string, error) {

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claim,
	)

	signedToken, err := token.SignedString([]byte(service.secretKey))

	return signedToken, err
}

// NewService return token utilities
func NewService(secretKey string) TokenService {
	return &jwtServices{
		secretKey: secretKey,
	}
}

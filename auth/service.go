package auth

import (
	"errors"

	"github.com/golang-jwt/jwt/v4"
)

type Service interface {
	GenerateToken(userID uint) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}

type jwtService struct {
}

var SECRET_KEY = []byte("kapantamat00225588")

func NewService() *jwtService {
	return &jwtService{}
}
func (s *jwtService) GenerateToken(userID uint) (string, error) {
	claim := jwt.MapClaims{}
	claim["user_id"] = userID

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	singedToken, err := token.SignedString(SECRET_KEY)
	if err != nil {
		return singedToken, err
	}
	return singedToken, nil
}

func (s *jwtService) ValidateToken(token string) (*jwt.Token, error) {
	Token, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("Invalid token")
		}
		return []byte(SECRET_KEY), nil
	})
	if err != nil {
		return Token, err
	}
	return Token, nil
}

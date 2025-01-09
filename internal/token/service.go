package token

import (
	"catalog-digital-product/internal/custom"
	"time"

	"github.com/golang-jwt/jwt"
)

type TokenService interface {
	GenerateToken(userId int) (string, error)
	ValidateToken(encodedToken string) (*jwt.Token, error)
}

type TokenServiceImpl struct {
	Key []byte
}

func NewTokenService(key []byte) TokenService {
	return &TokenServiceImpl{Key: key}
}

func (t *TokenServiceImpl) GenerateToken(userId int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(time.Hour * 24 * 3).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(t.Key)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (t *TokenServiceImpl) ValidateToken(encodedToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(encodedToken, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, custom.ErrUnauthorized
		}
		return t.Key, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				return nil, custom.ErrUnauthorized
			}
		}
		return token, nil
	}

	return nil, custom.ErrUnauthorized
}

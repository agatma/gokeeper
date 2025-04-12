package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/google/uuid"
)

type Authenticator struct {
	secretKey string
	tokenExp  time.Duration
}

type Claims struct {
	jwt.RegisteredClaims
	UserID uuid.UUID
	Login  string
}

type Token string

func NewAuthJWT(secretKey string, tokenExp time.Duration) *Authenticator {
	return &Authenticator{
		secretKey: secretKey,
		tokenExp:  tokenExp,
	}
}

func (a *Authenticator) MakeJWT(ID uuid.UUID, login string) (Token, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(a.tokenExp)),
		},
		UserID: ID,
		Login:  login,
	})

	tokenString, err := token.SignedString([]byte(a.secretKey))
	if err != nil {
		return "", err
	}

	return Token(tokenString), nil
}

func (a *Authenticator) GetUserID(tokenStr string) (uuid.UUID, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("parse JWT error")
		}
		return []byte(a.secretKey), nil
	})
	if err != nil {
		return uuid.UUID{}, errors.New("auth error")
	}

	if !token.Valid {
		return uuid.UUID{}, errors.New("token is invalid")
	}

	return claims.UserID, nil
}

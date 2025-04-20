package fileworkers

import (
	"bytes"
	"gokeeper/pkg/auth"
	"gokeeper/pkg/domain"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JwtFileWorker struct {
	filePath string
}

func NewJwtFileWorker(filePath string) *JwtFileWorker {
	return &JwtFileWorker{
		filePath: filePath,
	}
}

func (jfw *JwtFileWorker) Set(jwt string) error {
	if !jfw.validateDate(jwt) {
		return domain.ErrJWTTokenError
	}

	file, err := os.OpenFile(jfw.filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write([]byte(jwt))
	return err
}

func (jfw *JwtFileWorker) Get() (string, error) {
	file, err := os.OpenFile(jfw.filePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var buf bytes.Buffer
	_, err = buf.ReadFrom(file)
	if err != nil {
		return "", err
	}

	jwt := buf.String()
	if jfw.validateDate(jwt) {
		return jwt, nil
	}

	return "", domain.ErrJWTTokenError
}

func (jfw *JwtFileWorker) validateDate(tokenStr string) bool {
	claims := &auth.Claims{}
	_, _, err := jwt.NewParser().ParseUnverified(tokenStr, claims)

	if err != nil {
		return false
	}

	return claims.ExpiresAt.After(time.Now())
}

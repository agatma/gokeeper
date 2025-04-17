package domain

import (
	"bytes"
	"database/sql/driver"
	"errors"
	"fmt"
	"time"
)

type Type int32

const (
	LOGIN_PASSWORD Type = iota
	TEXT
	BYTES
	CARD
)

func (t Type) MarshalJSON() ([]byte, error) {
	switch t {
	case LOGIN_PASSWORD:
		return []byte("\"LOGIN_PASSWORD\""), nil
	case TEXT:
		return []byte("\"TEXT\""), nil
	case BYTES:
		return []byte("\"BYTES\""), nil
	case CARD:
		return []byte("\"CARD\""), nil
	default:
		return nil, ErrPrivateDataBadFormat
	}
}

func (t Type) Value() (driver.Value, error) {
	switch t {
	case LOGIN_PASSWORD:
		return "LOGIN_PASSWORD", nil
	case TEXT:
		return "TEXT", nil
	case BYTES:
		return "BYTES", nil
	case CARD:
		return "CARD", nil
	default:
		return nil, errors.New("invalid data type")
	}
}

func (t *Type) Scan(value interface{}) error {
	sv, err := driver.String.ConvertValue(value)
	if err != nil {
		return errors.New("failed to scan Type")
	}
	v, ok := sv.(string)
	if !ok {
		return errors.New("failed to scan Type")
	}

	switch v {
	case "LOGIN_PASSWORD":
		*t = LOGIN_PASSWORD
	case "TEXT":
		*t = TEXT
	case "BYTES":
		*t = BYTES
	case "CARD":
		*t = CARD
	default:
		return errors.New("invalid type")
	}
	return nil
}

func (t *Type) UnmarshalJSON(data []byte) error {
	fmt.Println("data", string(data))
	switch {
	case bytes.Equal(data, []byte("\"LOGIN_PASSWORD\"")):
		*t = LOGIN_PASSWORD
	case bytes.Equal(data, []byte("\"TEXT\"")):
		*t = TEXT
	case bytes.Equal(data, []byte("\"BYTES\"")):
		*t = BYTES
	case bytes.Equal(data, []byte("\"CARD\"")):
		*t = CARD
	default:
		*t = LOGIN_PASSWORD
		//return ErrPrivateDataBadFormat
	}
	return nil
}

type Data struct {
	ID       string    `json:"id"`
	DataType Type      `json:"type"`
	MetaData []byte    `json:"meta"`
	Data     []byte    `json:"data"`
	SavedAt  time.Time `json:"saved_at"`
}

type DeleteRequest struct {
	ID        string    `json:"id"`
	DeletedAt time.Time `json:"deleted_at"`
}

type GetAllRequest struct {
	Limit  uint64 `json:"limit"`
	Offset uint64 `json:"offset"`
}

type LoginPasswordData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type CardData struct {
	Number string `json:"number"`
	Secure uint16 `json:"secure"`
	Name   string `json:"name"`
	Date   string `json:"date"`
}

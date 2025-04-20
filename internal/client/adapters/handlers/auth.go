package clients

import (
	"context"
	"encoding/json"
	"gokeeper/pkg/domain"
	"net/http"

	"github.com/go-resty/resty/v2"
)

type AuthClient struct {
	client *resty.Client
}

func NewAuthClient(client *resty.Client) *AuthClient {
	return &AuthClient{
		client: client,
	}
}

func (ac *AuthClient) Login(ctx context.Context, user domain.InUserRequest) (string, error) {
	body, err := json.Marshal(user)
	if err != nil {
		return "", err
	}
	resp, err := ac.client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post("/api/user/login")
	if err != nil {
		return "", err
	}

	switch resp.StatusCode() {
	case http.StatusUnauthorized:
		return "", domain.ErrUserAuthentication
	case http.StatusOK:
		jwt := resp.Header().Get("authorization")
		return jwt, nil
	default:
		return "", domain.ErrInternalServerError
	}
}

func (ac *AuthClient) Register(ctx context.Context, user domain.InUserRequest) (string, error) {
	body, err := json.Marshal(user)
	if err != nil {
		return "", err
	}
	resp, err := ac.client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post("/api/user/register")
	if err != nil {
		return "", err
	}

	switch resp.StatusCode() {
	case http.StatusConflict:
		return "", domain.ErrUserConflict
	case http.StatusOK:
		jwt := resp.Header().Get("authorization")
		return jwt, nil
	default:
		return "", domain.ErrInternalServerError
	}
}

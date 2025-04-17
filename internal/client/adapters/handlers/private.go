package clients

import (
	"context"
	"encoding/json"
	"gokeeper/pkg/domain"
	"net/http"
	"strconv"

	"github.com/go-resty/resty/v2"
)

type PrivateClient struct {
	client *resty.Client
}

func NewPrivateClient(client *resty.Client) *PrivateClient {
	return &PrivateClient{
		client: client,
	}
}

func (pc *PrivateClient) Save(ctx context.Context, pd domain.Data, jwt string) error {
	body, err := json.Marshal(pd)
	if err != nil {
		return err
	}
	resp, err := pc.client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		SetHeader("Authorization", jwt).
		Post("/api/private")
	if err != nil {
		return err
	}

	switch resp.StatusCode() {
	case http.StatusUnauthorized:
		return domain.ErrUserAuthentication
	case http.StatusConflict:
		return domain.ErrPrivateDataConflict
	case http.StatusBadRequest:
		return domain.ErrPrivateDataBadFormat
	case http.StatusOK:
		return nil
	default:
		return domain.ErrInternalServerError
	}
}

func (pc *PrivateClient) Delete(ctx context.Context, pd domain.DeleteRequest, jwt string) error {
	body, err := json.Marshal(pd)
	if err != nil {
		return err
	}
	resp, err := pc.client.R().
		SetContext(ctx).
		SetHeader("Authorization", jwt).
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Delete("/api/private")
	if err != nil {
		return err
	}

	switch resp.StatusCode() {
	case http.StatusUnauthorized:
		return domain.ErrUserAuthentication
	case http.StatusConflict:
		return domain.ErrPrivateDataConflict
	case http.StatusBadRequest:
		return domain.ErrPrivateDataBadFormat
	case http.StatusOK:
		return nil
	default:
		return domain.ErrInternalServerError
	}
}

func (pc *PrivateClient) Get(ctx context.Context, id string, jwt string) (*domain.Data, error) {
	resp, err := pc.client.R().
		SetContext(ctx).
		SetHeader("Authorization", jwt).
		Get("/api/private/" + id)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode() {
	case http.StatusUnauthorized:
		return nil, domain.ErrUserAuthentication
	case http.StatusNotFound:
		return nil, domain.ErrPrivateDataNotFound
	case http.StatusBadRequest:
		return nil, domain.ErrPrivateDataBadFormat
	case http.StatusOK:
		respBody := resp.Body()
		var pd domain.Data
		err = json.Unmarshal(respBody, &pd)
		if err != nil {
			return nil, err
		}
		return &pd, nil
	default:
		return nil, domain.ErrInternalServerError
	}
}

func (pc *PrivateClient) GetAll(ctx context.Context, pd domain.GetAllRequest, jwt string) ([]domain.Data, error) {
	resp, err := pc.client.R().
		SetContext(ctx).
		SetHeader("Authorization", jwt).
		Get("/api/private?limit=" + strconv.FormatUint(pd.Limit, 10) + "&offset=" + strconv.FormatUint(pd.Offset, 10))
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode() {
	case http.StatusUnauthorized:
		return nil, domain.ErrUserAuthentication
	case http.StatusBadRequest:
		return nil, domain.ErrPrivateDataBadFormat
	case http.StatusOK:
		respBody := resp.Body()
		var pd []domain.Data
		err = json.Unmarshal(respBody, &pd)
		if err != nil {
			return nil, err
		}
		return pd, nil
	default:
		return nil, domain.ErrInternalServerError
	}
}

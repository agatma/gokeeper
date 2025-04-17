package private

import (
	"context"
	"errors"
	"gokeeper/pkg/domain"
	"log"
)

type AuthService interface {
	Register(ctx context.Context, user domain.InUserRequest, saveJWT bool) error
	Login(ctx context.Context, user domain.InUserRequest, saveJWT bool) (string, error)
	GetJwt(ctx context.Context) (string, error)
}

type Client interface {
	Save(ctx context.Context, pd domain.Data, jwt string) error
	Delete(ctx context.Context, pd domain.DeleteRequest, jwt string) error
	Get(ctx context.Context, id string, jwt string) (*domain.Data, error)
	GetAll(ctx context.Context, pd domain.GetAllRequest, jwt string) ([]domain.Data, error)
}

type FileWorker interface {
	SaveMany(pd []domain.Data) error
	GetAll() ([]domain.Data, error)
	DeleteAll() error
}

type Encrypter interface {
	EncryptMessage(msg []byte, secrets ...string) ([]byte, error)
	DecryptMessage(msg []byte, secrets ...string) ([]byte, error)
}

type BulkSender interface {
	Send(ctx context.Context, pds []domain.Data, jwt string) error
}

type Service struct {
	authService       AuthService
	privateClient     Client
	encrypter         Encrypter
	privateFileWorker FileWorker
	privateBulkSender BulkSender
}

func NewPrivateService(
	authService AuthService,
	privateClient Client,
	encrypter Encrypter,
	privateFileWorker FileWorker,
	privateBulkSender BulkSender,
) *Service {
	return &Service{
		authService:       authService,
		privateClient:     privateClient,
		encrypter:         encrypter,
		privateFileWorker: privateFileWorker,
		privateBulkSender: privateBulkSender,
	}
}

func (ps *Service) authorizeUser(ctx context.Context, inputUser *domain.InUserRequest) (string, error) {
	jwt, err := ps.authService.GetJwt(ctx)
	if err != nil {
		if inputUser != nil {
			jwt, err = ps.authService.Login(ctx, *inputUser, true)
			if err != nil {
				return "", err
			}
			return jwt, nil
		}
		return "", err
	}
	return jwt, nil
}

func (ps *Service) Save(ctx context.Context, pd domain.Data, inputUser domain.InUserRequest, saveLocalOnError bool) error {
	jwt, err := ps.authorizeUser(ctx, &inputUser)
	if err != nil {
		return err
	}

	pd.Data, err = ps.encrypter.EncryptMessage(pd.Data, inputUser.Login, inputUser.Password)
	if err != nil {
		return err
	}

	clientErr := ps.privateClient.Save(ctx, pd, jwt)
	if clientErr != nil {
		if errors.Is(clientErr, domain.ErrPrivateDataConflict) || errors.Is(clientErr, domain.ErrPrivateDataBadFormat) {
			return clientErr
		}
		if saveLocalOnError {
			fileWorkerErr := ps.privateFileWorker.SaveMany([]domain.Data{pd})
			if fileWorkerErr != nil {
				return fileWorkerErr
			}
			return domain.WarnServerUnavailable
		}
		return clientErr
	}
	return nil
}

func (ps *Service) GetAll(ctx context.Context, gpr domain.GetAllRequest, inputUser domain.InUserRequest) ([]domain.Data, error) {
	jwt, err := ps.authorizeUser(ctx, &inputUser)
	if err != nil {
		return nil, err
	}

	pds, err := ps.privateClient.GetAll(ctx, gpr, jwt)
	if err != nil {
		return nil, err
	}

	for idx := range pds {
		pds[idx].Data, err = ps.encrypter.DecryptMessage(pds[idx].Data, inputUser.Login, inputUser.Password)
		if err != nil {
			return nil, err
		}
	}
	return pds, nil
}

func (ps *Service) Get(ctx context.Context, id string, inputUser domain.InUserRequest) (*domain.Data, error) {
	jwt, err := ps.authorizeUser(ctx, &inputUser)
	if err != nil {
		return nil, err
	}

	pd, err := ps.privateClient.Get(ctx, id, jwt)
	if err != nil {
		return nil, err
	}

	pd.Data, err = ps.encrypter.DecryptMessage(pd.Data, inputUser.Login, inputUser.Password)
	if err != nil {
		return nil, err
	}
	return pd, nil
}

func (ps *Service) Delete(ctx context.Context, pd domain.DeleteRequest) error {
	jwt, err := ps.authorizeUser(ctx, nil)
	if err != nil {
		return err
	}

	return ps.privateClient.Delete(ctx, pd, jwt)
}

func (ps *Service) Upload(ctx context.Context) error {
	jwt, err := ps.authorizeUser(ctx, nil)
	if err != nil {
		return err
	}

	pds, err := ps.privateFileWorker.GetAll()
	if err != nil {
		return err
	}

	err = ps.privateBulkSender.Send(ctx, pds, jwt)
	if err != nil {
		return err
	}

	err = ps.privateFileWorker.DeleteAll()
	if err != nil {
		log.Printf("Warn: failed to delete data form dist: %v", err)
	}
	return nil
}

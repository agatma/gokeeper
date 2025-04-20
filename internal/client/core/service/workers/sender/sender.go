package sender

import (
	"context"
	"errors"
	"gokeeper/pkg/domain"
	"time"

	"golang.org/x/sync/errgroup"
)

type PrivateClient interface {
	Save(ctx context.Context, pd domain.Data, jwt string) error
}

type Sender struct {
	workersNum     int64
	privateClient PrivateClient
}

func NewSender(workersNum int64, privateClient PrivateClient) *Sender {
	return &Sender{
		workersNum:     workersNum,
		privateClient: privateClient,
	}
}

func (ps *Sender) doWork(ctx context.Context, pdChannel <-chan domain.Data, jwt string) error {
	for pd := range pdChannel {
		select {
		case <-ctx.Done():
			return errors.New("graceful shutdown")
		default:
			err := ps.privateClient.Save(ctx, pd, jwt)
			if err != nil {
				return err
			}
		}
		time.Sleep(time.Second)
	}
	return nil
}

func (ps *Sender) Send(ctx context.Context, pds []domain.Data, jwt string) error {
	pdChannel := make(chan domain.Data, len(pds))
	for _, pd := range pds {
		pdChannel <- pd
	}
	close(pdChannel)

	wg := new(errgroup.Group)

	for w := 0; w < int(ps.workersNum); w++ {
		wg.Go(func() error {
			return ps.doWork(ctx, pdChannel, jwt)
		})
	}

	if err := wg.Wait(); err != nil {
		return err
	}
	return nil
}

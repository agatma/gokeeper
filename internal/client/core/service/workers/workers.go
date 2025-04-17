package workers

import (
	"gokeeper/internal/client/core/config"
	fileworkers "gokeeper/internal/client/core/service/workers/files"
	"gokeeper/internal/client/core/service/workers/sender"
)

type Workers struct {
	FileWorker *fileworkers.FileWorkers
	Sender     *sender.Sender
}

func NewWorkers(cfg *config.Config, privateClient sender.PrivateClient) *Workers {
	return &Workers{
		FileWorker:     fileworkers.NewFileWorkers(cfg),
		Sender: sender.NewSender(int64(cfg.SenderWorkersNum), privateClient),
	}
}

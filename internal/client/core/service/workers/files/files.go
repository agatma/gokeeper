package fileworkers

import "gokeeper/internal/client/core/config"

type FileWorkers struct {
	JWTWorker         *JwtFileWorker
	PrivateFileWorker *PrivateFileWorker
}

func NewFileWorkers(cfg *config.Config) *FileWorkers {
	return &FileWorkers{
		JWTWorker:         NewJwtFileWorker(cfg.JWTPath),
		PrivateFileWorker: NewPrivateFileWorker(cfg.PrivateDataPath),
	}
}

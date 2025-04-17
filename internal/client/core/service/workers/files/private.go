package fileworkers

import (
	"bytes"
	"encoding/json"
	"gokeeper/pkg/domain"
	"os"
)

type PrivateFileWorker struct {
	filePath string
}

func NewPrivateFileWorker(filePath string) *PrivateFileWorker {
	return &PrivateFileWorker{
		filePath: filePath,
	}
}

func (pfw *PrivateFileWorker) SaveMany(pd []domain.Data) error {
	savedPrivateData, err := pfw.GetAll()
	if err != nil {
		return err
	}

	savedPrivateData = append(savedPrivateData, pd...)
	newPrivateData, err := json.Marshal(savedPrivateData)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(pfw.filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(newPrivateData)
	return err
}

func (pfw *PrivateFileWorker) GetAll() ([]domain.Data, error) {
	file, err := os.OpenFile(pfw.filePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var buf bytes.Buffer
	_, err = buf.ReadFrom(file)
	if err != nil {
		return nil, err
	}

	if buf.Len() == 0 {
		return []domain.Data{}, nil
	}

	var savedPrivateData []domain.Data
	err = json.Unmarshal(buf.Bytes(), &savedPrivateData)
	if err != nil {
		return nil, err
	}

	return savedPrivateData, nil
}

func (pfw *PrivateFileWorker) DeleteAll() error {
	return os.Remove(pfw.filePath)
}

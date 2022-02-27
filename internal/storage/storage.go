package storage

import (
	"YP-metrics-and-alerting/internal/repository"
	"encoding/json"
	"io"
	"log"
	"os"
)

type storage struct {
	fileName string
}

type Storage interface {
	Save(data []byte) error
	Retrieve() (*repository.MapStorageRepo, error)
}

func NewFileStorage(fileName string) *storage {
	return &storage{fileName: fileName}
}

func (s *storage) Save(data []byte) error {
	log.Println("save to file...")

	f, err := os.Create(s.fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func (s *storage) Retrieve() (*repository.MapStorageRepo, error) {
	log.Println("restore from file...")

	f, err := os.Open(s.fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	mapStorage := &repository.MapStorageRepo{}
	err = json.Unmarshal(data, mapStorage)
	if err != nil {
		return nil, err
	}
	return mapStorage, nil
}

package storage

import (
	"YP-metrics-and-alerting/internal/models"
	"encoding/json"
	"io"
	"log"
	"os"
)

type JSONFileStorage struct {
	fileName string
}

func NewJSONFileStorage(fileName string) *JSONFileStorage {
	return &JSONFileStorage{fileName: fileName}
}

func (s *JSONFileStorage) Save(data []byte) error {
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

func (s *JSONFileStorage) Retrieve() ([]*models.Metrics, error) {
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

	var slice []*models.Metrics
	err = json.Unmarshal(data, &slice)
	if err != nil {
		return nil, err
	}

	return slice, nil
}

package repository

import (
	"YP-metrics-and-alerting/internal/models"
	"testing"
)

func TestNewMapStorageRepo(t *testing.T) {
	repo := NewMapStorageRepo()

	var delta int64 = 1
	counterMetric := &models.Metrics{
		ID:    "PollCount",
		MType: "counter",
		Delta: &delta,
	}
	err := repo.SetMetric(counterMetric)
	if err != nil {
		t.Errorf("Error while setting counterMetric")
	}

	m, err := repo.GetMetric(counterMetric.ID)
	if err != nil {
		t.Errorf("Error while getting counterMetric: %s", counterMetric.ID)
	}

	expected := int64(1)
	if *m.Delta != expected {
		t.Errorf("Expected %d, received %d", expected, *m.Delta)
	}
	t.Log(*m.Delta)

	err = repo.SetMetric(counterMetric)
	if err != nil {
		t.Errorf("Error while setting counterMetric")
	}
	m, err = repo.GetMetric(counterMetric.ID)
	if err != nil {
		t.Errorf("Error while getting counterMetric: %s", counterMetric.ID)
	}
	expected = int64(2)
	if *m.Delta != expected {
		t.Errorf("Expected %d, received %d", expected, *m.Delta)
	}
	t.Log(*m.Delta)
}

func TestNewMapStorageRepo_GaugeMetric(t *testing.T) {
	repo := NewMapStorageRepo()

	var expected = 65637.019
	gaugeMetric := &models.Metrics{
		ID:    "testSetGet134",
		MType: "gauge",
		Value: &expected,
	}
	err := repo.SetMetric(gaugeMetric)
	if err != nil {
		t.Errorf("Error while setting gaugeMetric")
	}
	m, err := repo.GetMetric(gaugeMetric.ID)
	if err != nil {
		t.Errorf("Error while getting counterMetric: %s", m.ID)
	}

	if *m.Value != expected {
		t.Errorf("Expected %f, received %f", expected, *m.Value)
	}
	t.Log(*m.Value)

	expected = 156519.255
	gaugeMetric2 := &models.Metrics{
		ID:    "testSetGet134",
		MType: "gauge",
		Value: &expected,
	}
	err = repo.SetMetric(gaugeMetric2)
	if err != nil {
		t.Errorf("Error while setting gaugeMetric")
	}
	m, err = repo.GetMetric(gaugeMetric2.ID)
	if err != nil {
		t.Errorf("Error while getting counterMetric: %s", gaugeMetric2.ID)
	}
	if *m.Value != expected {
		t.Errorf("Expected %f, received %f", expected, *m.Value)
	}
	t.Log(*m.Value)
}

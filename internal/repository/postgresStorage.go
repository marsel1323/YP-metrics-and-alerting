package repository

import (
	"YP-metrics-and-alerting/internal/models"
	"context"
	"database/sql"
	"time"
)

type PostgresStorage struct {
	DB *sql.DB
}

func NewPostgresStorage(db *sql.DB) *PostgresStorage {
	return &PostgresStorage{
		DB: db,
	}
}

func (postgres PostgresStorage) GetMetric(id string) (*models.Metrics, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row := postgres.DB.QueryRowContext(
		ctx,
		`SELECT id, type, delta, value from metrics WHERE id = $1`,
		id,
	)

	if err := row.Err(); err != nil {
		return nil, err
	}

	var metric models.Metrics

	err := row.Scan(
		&metric.ID,
		&metric.MType,
		&metric.Delta,
		&metric.Value,
	)
	if err != nil {
		return nil, err
	}

	return &metric, nil
}

func (postgres *PostgresStorage) SetMetric(metric *models.Metrics) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := postgres.DB.ExecContext(
		ctx,
		`
				INSERT INTO metrics (id, type, delta, value) 
				VALUES ($1, $2, $3, $4)
				ON CONFLICT (id)
				DO UPDATE SET delta = metrics.delta + $3, value = $4;
			`,
		metric.ID,
		metric.MType,
		metric.Delta,
		metric.Value,
	)
	if err != nil {
		return err
	}
	return nil
}

func (postgres PostgresStorage) GetMetricsList() ([]*models.Metrics, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := postgres.DB.QueryContext(
		ctx,
		`SELECT id, type, delta, value from metrics`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []*models.Metrics

	for rows.Next() {
		var metric models.Metrics
		err := rows.Scan(
			&metric.ID,
			&metric.MType,
			&metric.Delta,
			&metric.Value,
		)
		if err != nil {
			return nil, err
		}

		metrics = append(metrics, &metric)
	}

	return metrics, nil
}

func (postgres PostgresStorage) SetMetricsList(metricsList []*models.Metrics) error {
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, metric := range metricsList {
		err := postgres.SetMetric(metric)
		if err != nil {
			return err
		}
	}
	return nil
}

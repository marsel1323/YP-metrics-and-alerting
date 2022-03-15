package repository

type DBRepo interface {
	GetAllGaugeMetricValues() (map[string]float64, error)
	GetGaugeMetricValue(string) (float64, error)
	SetGaugeMetricValue(string, float64) error

	GetAllCounterMetricValues() (map[string]int64, error)
	GetCounterMetricValue(string) (int64, error)
	SetCounterMetricValue(string, int64) error

	BunchSetMetrics(mapStorage *MapStorageRepo) error
}

package gql

type MetricsData struct {
	Metrics []Metric `json:"metrics"`
}

type Metric struct {
	Measurement string        `json:"measurement"`
	Values      []MetricValue `json:"values"`
}

type MetricValue struct {
	Timestamp int     `json:"ts"`
	Value     float64 `json:"value"`
}

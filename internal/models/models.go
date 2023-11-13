package models

type MetricType int

const (
	Gauge MetricType = iota
	Counter
)

type MetricReq struct {
	Name  string
	Type  MetricType
	Value float64
}

type MetricResp struct {
	Name  string
	Type  MetricType
	Value float64
}

type GaugeMetric struct {
	Name  string
	Type  MetricType
	Value float64
}

type CounterMetric struct {
	Name  string
	Type  MetricType
	Value int64
}

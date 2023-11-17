package models

type MetricType int

const (
	Gauge MetricType = iota
	Counter
)

type Metric struct {
	Value any
	Name  string
	Type  MetricType
}

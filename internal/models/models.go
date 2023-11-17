package models

type MetricType int

const (
	Gauge MetricType = iota
	Counter
)

type Metric struct {
	Name  string
	Type  MetricType
	Value any
}

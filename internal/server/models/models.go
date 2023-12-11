package models

const (
	Gauge   string = "gauge"
	Counter string = "counter"
)

type Metric struct {
	Value any
	Name  string
	Type  string
}

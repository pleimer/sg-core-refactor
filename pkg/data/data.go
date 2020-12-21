package data

import "time"

// package data defines the data descriptions for objects used in the internal buses

//Type describes internal message types
type Type int

const (
	//EVENT ...
	EVENT Type = iota
	//METRIC ...
	METRIC
)

// MetricType follows standard metric conventions followed by prometheus and
// collectd
type MetricType int

const (
	//COUNTER ...
	COUNTER MetricType = iota
	//GAUGE ...
	GAUGE
	//UNTYPED ...
	UNTYPED
)

func (t Type) String() string {
	return [...]string{"EVENT", "METRIC"}[t]
}

// Event internal event type
type Event struct {
	Handler string
	Message string
}

// Metric internal metric type
type Metric struct {
	Name   string
	Labels map[string]string
	Time   *time.Time
	Type   MetricType
	Value  float64
}

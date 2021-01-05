package data

import (
	"fmt"
	"time"
)

// package data defines the data descriptions for objects used in the internal buses

//Type describes internal message types
type Type int

const (
	//EVENT ...
	EVENT Type = iota
	//METRIC ...
	METRIC
)

func (t Type) String() string {
	return [...]string{"EVENT", "METRIC"}[t]
}

// MetricType follows standard metric conventions from prometheus and
// collectd
type MetricType int

const (
	//COUNTER ...
	COUNTER MetricType = iota
	//GAUGE ...
	GAUGE
	//DERIVE ...
	DERIVE
	//UNTYPED ...
	UNTYPED
)

//FromString set metric type from string
func (mt MetricType) FromString(key string) (MetricType, error) {
	if mt, ok := map[string]MetricType{
		"counter": COUNTER,
		"gauge":   GAUGE,
		"derive":  DERIVE,
		"untyped": UNTYPED,
	}[key]; ok {
		return mt, nil
	}
	return UNTYPED, fmt.Errorf("undefined metric type '%s'", key)
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
	Time   time.Time
	Type   MetricType
	Value  float64
}

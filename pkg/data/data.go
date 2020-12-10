package data

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

// Event internal event type
type Event struct {
	Message string
}

// Metric internal metric type
type Metric struct {
	Message string
}

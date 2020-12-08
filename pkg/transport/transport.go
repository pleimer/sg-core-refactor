package transport

import (
	"strings"
	"sync"
)

// package transport defines the interfaces for interacting with transport
// plugins

//Mode indicates if transport is setup to receive or write
type Mode int

const (
	//WRITE ...
	WRITE = iota
	//READ ...
	READ
)

//String get string representation of mode
func (m Mode) String() string {
	return [...]string{"WRITE", "READ"}[m]
}

//FromString get mode from string
func (m Mode) FromString(s string) {
	m = map[string]Mode{
		"write": WRITE,
		"read":  READ,
	}[strings.ToLower(s)]
}

//Transport type listens on one interface and delivers data to core
type Transport interface {
	Config(interface{}) error
	Run(*sync.WaitGroup, chan []byte)
}

//NewFn transport New function must be of this type
type NewFn func() Transport

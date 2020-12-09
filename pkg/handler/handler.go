package handler

import "github.com/infrawatch/sg-core-refactor/pkg/data"

// package handler contains the interface description for handler plugins

// Handler
type Handler interface {
	Handle([]byte) (interface{}, error)
	Type() data.Type
}

// NewFn New func must be of this type
type NewFn func() Handler

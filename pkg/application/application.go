package application

import "sync"

//package application defines the interface for interacting with application plugins

//Application describes application plugin interfaces
type Application interface {
	Config(interface{}) error
	Run(*sync.WaitGroup, interface{})
}

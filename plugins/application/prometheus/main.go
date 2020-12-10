package main

import (
	"fmt"
	"sync"

	"github.com/infrawatch/sg-core-refactor/pkg/application"
	"github.com/infrawatch/sg-core-refactor/pkg/data"
)

//Application plugins should also provide a callback function with the data
// and the Run() function is separate for

//Prometheus plugin for interfacing with Prometheus
type Prometheus struct {
}

//Run run scrape endpoint
func (p *Prometheus) Run(wg *sync.WaitGroup, eChan chan data.Event, mChan chan data.Metric) {
	defer wg.Done()

	for {
		select {
		case ev := <-eChan:
			fmt.Printf("Prometheus received event with message: %s\n", ev.Message)
		case m := <-mChan:
			fmt.Printf("Prometheus received metric with message: %s\n", m.Message)
		}
	}
}

//Config implements application.Application
func (p *Prometheus) Config(c interface{}) error {
	return nil
}

//New constructor
func New() application.Application {
	return &Prometheus{}
}

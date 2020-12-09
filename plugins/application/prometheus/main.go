package main

import (
	"fmt"
	"sync"

	"github.com/infrawatch/sg-core-refactor/pkg/application"
	"github.com/infrawatch/sg-core-refactor/pkg/bus"
	"github.com/infrawatch/sg-core-refactor/pkg/data"
)

//Prometheus plugin for interfacing with Prometheus
type Prometheus struct {
}

//Run run scrape endpoint
func (p *Prometheus) Run(wg *sync.WaitGroup, t interface{}) {
	defer wg.Done()

	c := make(chan data.Event)
	b := t.(*bus.EventBus)
	b.Subscribe(c)

	fmt.Printf("prometheus app received event: %s\n", (<-c).Message)
}

//Config implements application.Application
func (p *Prometheus) Config(c interface{}) error {
	return nil
}

//New constructor
func New() application.Application {
	return &Prometheus{}
}

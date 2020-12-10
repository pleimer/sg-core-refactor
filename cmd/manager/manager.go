package manager

import (
	"fmt"
	"path/filepath"
	"plugin"
	"strings"
	"sync"

	"github.com/infrawatch/apputils/logging"
	"github.com/infrawatch/sg-core-refactor/pkg/application"
	"github.com/infrawatch/sg-core-refactor/pkg/bus"
	"github.com/infrawatch/sg-core-refactor/pkg/data"
	"github.com/infrawatch/sg-core-refactor/pkg/handler"
	"github.com/infrawatch/sg-core-refactor/pkg/transport"
	"github.com/pkg/errors"
)

var (
	transports     map[string]transport.Transport
	metricHandlers map[string][]handler.MetricHandler
	eventHandlers  map[string][]handler.EventHandler
	applications   map[string]application.Application
	eventBus       bus.EventBus
	metricBus      bus.MetricBus
	pluginPath     string
	logger         *logging.Logger
)

func init() {
	transports = map[string]transport.Transport{}
	metricHandlers = map[string][]handler.MetricHandler{}
	eventHandlers = map[string][]handler.EventHandler{}
	applications = map[string]application.Application{}
	pluginPath = "/usr/lib64/sg-core"
}

func eventHandleDecorator(data []byte, call func([]byte) (data.Event, error)) {
	e, err := call(data)
	if err != nil {
		logger.Metadata(logging.Metadata{"error": err})
		logger.Error("cannot publish event to event bus")
		return
	}
	eventBus.Publish(e)
}

func eventMetricDecorator(data []byte, call func([]byte) (data.Metric, error)) {
	e, err := call(data)
	if err != nil {
		logger.Metadata(logging.Metadata{"error": err})
		logger.Error("cannot publish event to event bus")
		return
	}
	metricBus.Publish(e)
}

//SetPluginDir set directory path containing plugin binaries
func SetPluginDir(path string) {
	pluginPath = path
}

//SetLogger set logger
func SetLogger(l *logging.Logger) {
	logger = l
}

//InitTransport load tranpsort binary and initialize with config
func InitTransport(name string, mode string, config interface{}) error {
	n, err := initPlugin(name)
	if err != nil {
		return errors.Wrap(err, "failed initializing transport")
	}

	new, ok := n.(func() transport.Transport)
	if !ok {
		return fmt.Errorf("plugin %s constructor 'New' did not return type 'transport.Transport'", name)
	}

	transports[name] = new()

	if config == nil {
		return nil
	}

	err = transports[name].Config(config)
	if err != nil {
		return err
	}
	return nil
}

//InitApplication initialize application plugin with configuration
func InitApplication(name string, config interface{}) error {
	n, err := initPlugin(name)
	if err != nil {
		return errors.Wrap(err, "failed initializing application plugin")
	}

	new, ok := n.(func() application.Application)
	if !ok {
		return fmt.Errorf("plugin %s constructor 'New' did not return type 'application.Application'", name)
	}

	applications[name] = new()

	if config == nil {
		return nil
	}

	err = applications[name].Config(config)
	if err != nil {
		return err
	}
	return nil
}

//SetTransportHandlers load handlers binaries for transport
func SetTransportHandlers(name string, handlerNames []string) error {
	for _, hName := range handlerNames {
		n, err := initPlugin(hName)
		if err != nil {
			return errors.Wrap(err, "failed initializing handler")
		}

		var hType string
		switch new := n.(type) {
		case (func() handler.MetricHandler):
			metricHandlers[name] = append(metricHandlers[name], new())
			hType = "MetricHandler"
		case (func() handler.EventHandler):
			eventHandlers[name] = append(eventHandlers[name], new())
			hType = "EventHandler"
		default:
			return fmt.Errorf("handler %s constructor did not return type handler.EventHandler or handler.MetricsHandler", hName)
		}
		logger.Metadata(logging.Metadata{"handler": hName, "type": hType})
		logger.Info("initialized handler")
	}
	return nil
}

//RunTransports spins off tranpsort + handler processes
func RunTransports(wg *sync.WaitGroup) {
	for name, t := range transports {
		exchange := make(chan []byte)

		wg.Add(2)
		go t.Run(wg, exchange)
		go func(wg *sync.WaitGroup, name string) {
			defer wg.Done()
			for _, handler := range metricHandlers[name] {
				res, err := handler.Handle(<-exchange)
				if err != nil {
					logger.Metadata(logging.Metadata{"error": err})
					logger.Error("failed handling message")
					continue
				}
				metricBus.Publish(res)
			}
			for _, handler := range eventHandlers[name] {
				res, err := handler.Handle(<-exchange)
				if err != nil {
					logger.Metadata(logging.Metadata{"error": err})
					logger.Error("failed handling message")
					continue
				}
				eventBus.Publish(res)
			}
		}(wg, name)
	}
}

//RunApplications spins off application processes
func RunApplications(wg *sync.WaitGroup) {
	for _, a := range applications {
		eChan := make(chan data.Event)
		mChan := make(chan data.Metric)

		eventBus.Subscribe(eChan)
		metricBus.Subscribe(mChan)
		wg.Add(1)
		go a.Run(wg, eChan, mChan)
	}
}

// helper functions

func initPlugin(name string) (plugin.Symbol, error) {
	bin := strings.Join([]string{name, "so"}, ".")
	path := filepath.Join(pluginPath, bin)
	p, err := plugin.Open(path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open binary %s", path)
	}

	n, err := p.Lookup("New")
	return n, err
}

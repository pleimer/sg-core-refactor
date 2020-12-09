package manager

import (
	"fmt"
	"path/filepath"
	"plugin"
	"strings"
	"sync"

	"github.com/infrawatch/apputils/logging"
	"github.com/infrawatch/sg-core-refactor/pkg/application"
	"github.com/infrawatch/sg-core-refactor/pkg/handler"
	"github.com/infrawatch/sg-core-refactor/pkg/transport"
	"github.com/pkg/errors"
)

// TODO: give transport and associated handlers the same channel to talk to eachother with

var (
	transports   map[string]transport.Transport
	handlers     map[string][]handler.Handler
	applications map[string]application.Application
	pluginPath   string
	logger       *logging.Logger
)

func init() {
	transports = map[string]transport.Transport{}
	handlers = map[string][]handler.Handler{}
	applications = map[string]application.Application{}
	pluginPath = "/usr/lib64/sg-core"
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
		return fmt.Errorf("plugin %s constructor 'New' is not of type 'func() transport.Transport'", name)
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

//SetTransportHandlers load handlers binaries for transport
func SetTransportHandlers(name string, handlerNames []string) error {
	for _, hName := range handlerNames {
		n, err := initPlugin(hName)
		if err != nil {
			return errors.Wrap(err, "failed initializing handler")
		}
		new, ok := n.(func() handler.Handler)
		if !ok {
			return fmt.Errorf("plugin %s constructor 'New' is not of type 'func() handler.Handler'", name)
		}

		handlers[name] = append(handlers[name], new())
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
			for _, handler := range handlers[name] {
				res, err := handler.Handle(<-exchange)
				if err != nil {
					logger.Metadata(logging.Metadata{"error": err})
					logger.Error("failed handling message")
				}
				logger.Metadata(logging.Metadata{"result": string(res)})
				logger.Info("handled message")
			}
		}(wg, name)
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

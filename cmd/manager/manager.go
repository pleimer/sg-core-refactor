package manager

import (
	"fmt"
	"path/filepath"
	"plugin"
	"strings"
	"sync"

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

//SetReceiver set up transport plugin with configuration.
func SetReceiver(transportName string, mode string, handlers []string, config interface{}) error {

	return nil
}

//InitTransport load tranpsort binary and initialize with config
func InitTransport(name string, mode string, config interface{}) error {
	n, path, err := initPlugin(name)
	if err != nil {
		return errors.Wrapf(err, "failed to open transport constructor 'New' in binary %s", path)
	}

	new, ok := n.(func() transport.Transport)
	if !ok {
		return fmt.Errorf("plugin %s constructor 'New' is not of type 'transport.NewFn'", name)
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
		n, path, err := initPlugin(hName)
		if err != nil {
			return errors.Wrapf(err, "failed to open handler constructor 'New' in binary %s", path)
		}
		new, ok := n.(func() handler.Handler)
		if !ok {
			return fmt.Errorf("plugin %s constructor 'New' is not of type 'handler.NewFn'", name)
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
				handler.Handle(<-exchange)
			}
		}(wg, name)
	}
}

// helper functions

func initPlugin(name string) (plugin.Symbol, string, error) {
	bin := strings.Join([]string{name, "so"}, ".")
	path := filepath.Join(pluginPath, bin)
	p, err := plugin.Open(path)
	if err != nil {
		return nil, "", errors.Wrapf(err, "failed to open plugin binary %s", path)
	}

	n, err := p.Lookup("New")
	return n, path, err
}

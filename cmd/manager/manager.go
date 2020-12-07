package manager

import (
	"fmt"
	"path/filepath"
	"plugin"
	"strings"

	"github.com/infrawatch/sg-core-refactor/pkg/application"
	"github.com/infrawatch/sg-core-refactor/pkg/handler"
	"github.com/infrawatch/sg-core-refactor/pkg/transport"
	"github.com/pkg/errors"
)

var (
	transports   map[string]*transportPipe
	handlers     map[string]handler.Handler
	applications map[string]application.Application
	pluginPath   string
)

// handlers are associated with tranport plugins
type transportPipe struct {
	instance transport.Transport
	handlers []handler.Handler
}

func init() {
	transports = map[string]*transportPipe{}
	handlers = map[string]handler.Handler{}
	applications = map[string]application.Application{}
	pluginPath = "/usr/lib64/sg-core"
}

//SetPluginDir set directory path containing plugin binaries
func SetPluginDir(path string) {
	pluginPath = path
}

//SetTransport set up transport plugin with configuration.
func SetTransport(name string, mode string, handlers []string, config interface{}) error {
	n, path, err := initPlugin(name)
	if err != nil {
		return errors.Wrapf(err, "failed to open transport constructor 'New' in binary %s", path)
	}

	new, ok := n.(func() transport.Transport)
	if !ok {
		return fmt.Errorf("plugin %s constructor 'New' is not of type 'transport.NewFn'", name)
	}

	transports[name] = &transportPipe{
		instance: new(),
		handlers: []handler.Handler{},
	}

	for _, hName := range handlers {
		n, path, err := initPlugin(hName)
		if err != nil {
			return errors.Wrapf(err, "failed to open handler constructor 'New' in binary %s", path)
		}
		new, ok := n.(func() handler.Handler)
		if !ok {
			return fmt.Errorf("plugin %s constructor 'New' is not of type 'handler.NewFn'", name)
		}

		transports[name].handlers = append(transports[name].handlers, new())
	}

	if config == nil {
		return nil
	}

	err = transports[name].instance.Config(config)
	if err != nil {
		return err
	}
	return nil
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

package manager

import (
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"strings"

	"github.com/infrawatch/sg-core-refactor/pkg/transport"
)

var (
	binPaths map[string]string
	plugins  []interface{}
)

func init() {
	binPaths = map[string]string{}
	plugins = []interface{}{}
}

// LoadBinaries load binary shared objects from directory
func LoadBinaries(dir string) error {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) != ".so" {
			return nil
		}

		baseName := strings.Split(info.Name(), ".")
		binPaths[baseName[len(baseName)-2]] = path
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// SetPluginConfig pass config object to plugin
func SetPluginConfig(name string, config interface{}) error {
	if _, ok := binPaths[name]; !ok {
		return fmt.Errorf("could not load plugin '%s': binary does not exist", name)
	}

	p, err := plugin.Open(binPaths[name])
	if err != nil {
		return fmt.Errorf("could not load plugin %s: %s", name, err)
	}

	n, err := p.Lookup("New")
	if err != nil {
		return fmt.Errorf("could not load plugin 'New' function %s: %s", name, err)
	}

	o := n.(func() transport.Transport)()
	plugins = append(plugins, o)
	err = o.Config(config)
	if err != nil {
		return err
	}
	return nil
}

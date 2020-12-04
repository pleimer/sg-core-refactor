package manager

import (
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"strings"

	"github.com/infrawatch/sg-core-refactor/pkg/extension"
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
		return fmt.Errorf("could not load 'New' constructor for plugin %s: %s", name, err)
	}

	e := n.(func() extension.Extension)()
	err = e.Config(config)
	if err != nil {
		return err
	}
	plugins = append(plugins, e)
	return nil
}

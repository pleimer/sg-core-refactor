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
	extensions map[string]extension.Extension
)

func init() {
	extensions = map[string]extension.Extension{}
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

		p, err := plugin.Open(path)
		if err != nil {
			return fmt.Errorf("could not open plugin %s: %s", path, err)
		}

		baseName := strings.Split(info.Name(), ".")

		n, err := p.Lookup("New")
		if err != nil {
			return fmt.Errorf("could not load 'New' constructor for plugin %s: %s", baseName, err)
		}

		e := n.(func() extension.Extension)()
		extensions[baseName[len(baseName)-2]] = e
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// SetPluginConfig pass config object to plugin
func SetPluginConfig(name string, config interface{}) error {
	if _, ok := extensions[name]; !ok {
		return fmt.Errorf("could not load plugin '%s': binary does not exist", name)
	}

	err := extensions[name].Config(config)
	if err != nil {
		return err
	}
	return nil
}

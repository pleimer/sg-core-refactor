package extension

// package extension defines the top level plugin interface that can be managed
// by the plugin manager

// Extension all plugins must implement the extension interface
type Extension interface {
	Config(interface{}) error
}

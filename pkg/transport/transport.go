package transport

// package transport defines the interfaces for interacting with transport
// plugins

//Transport type listens on one interface and delivers data to core
type Transport interface {
	Config(interface{}) error
}

package transport

//Transport type listens on one interface and delivers data to core
type Transport interface {
	Run()
}

package runner

// package runner defines the plugin interface for plugins that should run in
// their own process

// Runner plugins that implement this interface will run in a separate process
type Runner interface {
	Run()
}

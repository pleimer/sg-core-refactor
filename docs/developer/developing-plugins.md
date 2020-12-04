# Plugin types

1. Transport
2. Handler
3. Application


# Plugin configurations
Plugins should not read from their own cofiguration files. The sg-core provides
all methods for reading configurations using the golang 
[yaml.v3](https://pkg.go.dev/gopkg.in/yaml.v3). However, validation must be 
implemented by the plugin.

Sg-core looks within it's own configuration for plugin configurations, matching 
`name` field with the name of the shared object file. `config` field is then
passed to the plugin's `Config()` function. Each plugin is
responsible for validating the passed in configuration and should return an
error in the case that it fails.
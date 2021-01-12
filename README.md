# Plugins


Default plugins exist in /plugins. Plugins can also be hosted as separate projects.

## Build
```bash
# build sg-core and plugins. Places plugin binaries in ./bin
./build.sh
```
## Example Configuration
This configuration assumes both a Quipid Dispatch Router and Prometheus instance
are running on the localhost and listens for incoming messages on a unix socket
at `/tmp/smartgateway`.

```yaml
plugindir: bin/
loglevel: debug
transports:
  - name: socket
    handlers: 
      - collectd-metrics
    config:
      address: /tmp/smartgateway
applications: 
  - name: prometheus 
    config:
      host: localhost
      port: 8081
      metrictimeout: 30
```

## Run
`./sg-core -config <path to config>

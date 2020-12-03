# Plugins


Default plugins exist in /plugins. Plugins can also be hosted as separate projects.

# Build

```bash
# build plugins
go build -o bin/ -buildmode=plugin ./plugins/...

# build core
go build -o sg-core cmd/*.go
```
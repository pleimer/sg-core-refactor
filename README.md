# Plugins


Default plugins exist in /plugins. Plugins can also be hosted as separate projects.

# Build

```bash
# build plugins
for i in plugins/transport/*; do go build -o bin/ -buildmode=plugin "./$i/..."; done
for i in plugins/handler/*; do go build -o bin/ -buildmode=plugin "./$i/..."; done
for i in plugins/application/*; do go build -o bin/ -buildmode=plugin "./$i/..."; done

# build core
go build -o sg-core cmd/*.go
```

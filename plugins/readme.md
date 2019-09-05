# Build
```bash
go clean -cache
go build -buildmode=plugin -o ./plugins-bin/ElasticSearch6.so ./plugins/ElasticSearch6.go
```
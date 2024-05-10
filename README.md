# System Monitoring
GRPC server for collecting host metrics

## How to build
For bulding server and client run
```bash
make build
```

## How to run
```bash
make run-server
make run-client
```
or
```bash
go run cmd/server/main.go --port 50051
go run cmd/client/main.go --port 50051 --scrape 5 average 15
```

## Configuration
Configuration file availabel in `configs/config.yaml`

Config Example
```yaml
metrics:
  loadaverage: true
```

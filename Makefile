generate:
	rm -rf pb
	mkdir -p pb

	protoc \
		--proto_path=proto/ \
		--go_out=pb/ \
		--go-grpc_out=pb/ \
		proto/*.proto

build.server:
	go build -v -o ./bin/server ./cmd/server

build.client:
	go build -v -o ./bin/client ./cmd/client

run-server:
	go run cmd/server/main.go

run-client:
	go run cmd/client/main.go

test:
	go test -race -count 100 ./pkg/...
	go test -race -count 100 ./internal/...

lint:
	golangci-lint run ./...
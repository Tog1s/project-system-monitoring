generate:
	rm -rf pb
	mkdir -p pb

	protoc \
		--proto_path=proto/ \
		--go_out=pb/ \
		--go-grpc_out=pb/ \
		proto/*.proto

run-server:
	go run cmd/server/main.go

run-client:
	go run cmd/client/main.go


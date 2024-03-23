generate:
	rm -rf pb
	mkdir -p pb

	protoc \
		--proto_path=proto/ \
		--go_out=pb/ \
		--go-grpc_out=pb/ \
		proto/*.proto

run:
	go run cmd/main.go

run-server:
	go run server/main.go

run-client:
	go run client/main.go


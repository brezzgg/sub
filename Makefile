build_out := ./sub

build:
	rm -f $(build_out)
	go build -o $(build_out) cmd/sub/main.go

gen:
	protoc --go_out=. --go_opt=paths=source_relative \
    	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
    	internal/transport/grpc/pb/service.proto \
        internal/models/payload/v1/payload.proto

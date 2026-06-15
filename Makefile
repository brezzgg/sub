build_out := ./sub
image_name := ghcr.io/brezzgg/sub

build:
	rm -f $(build_out)
	go build -o $(build_out) cmd/sub/main.go


build-image:
	docker buildx build . --tag $(image_name)

gen:
	protoc --go_out=. --go_opt=paths=source_relative \
    	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
        internal/transport/grpc/service.proto \
        internal/usecase/model.proto

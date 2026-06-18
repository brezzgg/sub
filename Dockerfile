FROM golang:1.26.4-alpine3.23 AS builder

# copy src
WORKDIR /build
COPY . .

# download & build
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download && \
    go build -o /bin/sub 'cmd/sub/main.go'

FROM alpine:3.23.4
WORKDIR /workspace

COPY --from=builder /bin/sub /bin/sub
CMD ["sub", "run", "--storage=badger", "--storage-file=/workspace/database"]

.PHONY: proto
proto:
	protoc --go_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative \
		internal/app/proto/song.proto

.PHONY: install-tools
install-tools:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

.PHONY: deps
deps:
	go get -u google.golang.org/grpc
	go get -u google.golang.org/protobuf
	go mod tidy

.PHONY: all
all: install-tools deps proto

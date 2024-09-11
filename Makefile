LOCAL_BIN:=$(CURDIR)/bin

PHONY: install protogen

install:
	GOBIN=$(LOCAL_BIN) go install github.com/bufbuild/buf/cmd/buf@v1.29.0
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.32
	GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2

PB_DIR := pkg/pb

protogen:
	rm -rf $(PB_DIR) && mkdir $(PB_DIR) && \
	buf generate $(CURDIR)/api \
	|| rm -rf $(PB_DIR)


generate_descriptor:
	buf build -o image.binpb


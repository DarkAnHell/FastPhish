.PHONY: proto clean

binary = a.out

proto:
	@echo "Building .proto files..."
	@protoc -I proto -I${GOPATH}/src --go_out=plugins=grpc:proto proto/*/*.proto

clean:
	@rm -f proto/*/*.pb.go
	@rm $(binary)

build:
	@go build ./... $(binary)

all: clean proto build

.PHONY: proto clean

binary = a.out

api:
	@echo "Building .proto files..."
	@protoc -I proto -I${GOPATH}/src --go_out=plugins=grpc:proto api/*/*.proto

clean:
	@echo "Cleaning..."
	@rm -f api/*/*.pb.go
	@rm $(binary)

build:
	@echo "Building binaries"
	@go build ./... $(binary)

all: clean api build

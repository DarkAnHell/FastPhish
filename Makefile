.PHONY: proto clean all build api

binary = a.out

all: clean api build

api:
	@echo "Building .proto files..."
	@protoc -I api -I${GOPATH}/src --go_out=plugins=grpc:api api/*/*.proto

clean:
	@echo "Cleaning..."
	-rm -f api/*/*.pb.go
	-rm $(binary)

build:
	@echo "Building binaries"
	@go build ./... $(binary)



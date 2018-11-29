SHELL := /bin/bash

.PHONY: proto clean all build api

binary = a.out
apifiles = $(shell ls api)

all: clean api build

api:
	@echo "Building .proto files..."
	$(shell for p in $(apifiles); do protoc -I api --go_out=plugins=grpc:api api/$$p/*.proto; done)

clean:
	@echo "Cleaning..."
	-@rm -f api/*/*.pb.go
	-@rm $(binary)

build:
	@echo "Building binaries"
	@go build ./... $(binary)



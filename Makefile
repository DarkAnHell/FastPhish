SHELL := /bin/bash

.PHONY: proto clean all build api run

binary = a.out
apifiles = $(shell ls api)
cmd_path = cmd
binaries = $(shell ls $(cmd_path))
bin_path = bin

all: clean api build

api:
	@echo "Building .proto files..."
	$(shell for p in $(apifiles); do protoc -I api --go_out=plugins=grpc:api api/$$p/*.proto; done)

clean:
	@echo "Cleaning..."
	-@rm -f api/*/*.pb.go
	-@rm -r $(bin_path)

build:
	@echo "Building binaries"
	$(shell for b in $(binaries); do go build -o $(bin_path)/$$b ./$(cmd_path)/$$b; chmod +x ./$(cmd_path)/$$b; done)


SHELL := /bin/bash

.PHONY: proto clean all build api run test certs


binary = a.out
apifiles = $(shell ls api)
cmd_path = cmd
binaries = $(shell ls $(cmd_path))
bin_path = bin
pkg_path = pkg

all: clean api build certs

api:
	@echo "Building .proto files..."
	@protoc -I api --go_out=plugins=grpc:api api/*.proto

clean:
	@echo "Cleaning..."
	-@rm -f api/*.pb.go
	-@rm -r $(bin_path)
	-@rm -fr certs

build:
	@echo "Building binaries"
	$(shell for b in $(binaries); do go build -o $(bin_path)/$$b ./$(cmd_path)/$$b; chmod +x ./$(cmd_path)/$$b; done)

test: all
	go test -v ./...

certs:
	@mkdir -p certs
	@openssl req -x509 -newkey rsa:4096 -keyout certs/server.key -out certs/server.crt -days 3650 -nodes -subj '/CN=localhost'


PROTO_FILE = ./api/hello.proto

all: proto server

server:
	-mkdir ./bin > /dev/null 2>&1
	go build -o ./bin/server cmd/server.go

clean:
	rm -rf ./bin

proto:
	docker run \
		--privileged --rm \
		-v $(shell pwd):/go/src/code \
		-e "GOPATH=/go" \
		-e "DOCKER_PROTO=yes" \
		-e "PROTO_USER=$(shell id -u)" \
		-e "PROTO_GROUP=$(shell id -g)" \
		-e "PATH=/bin:/usr/bin:/usr/local/bin:/go/bin:/usr/local/go/bin" \
		quay.io/openstorage/grpc-framework:v0.0.2 \
			make docker-proto

docker-proto:
ifndef DOCKER_PROTO
	$(error Do not run directly. Run 'make proto' instead.)
endif
	grpcfw $(PROTO_FILE)


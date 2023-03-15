
PROTO_FILE = ./api/hello.proto
TAG = v0.0.7

all: proto server client

./bin:
	mkdir ./bin

server: ./bin
	go build -o ./bin/server cmd/server.go

client: ./bin
	go build -o ./bin/client cmd/client/client.go

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
		quay.io/openstorage/grpc-framework:$(TAG)\
			make docker-proto

proto-lint:
	docker run \
		--privileged --rm \
		-v $(shell pwd):/go/src/code \
		-e "GOPATH=/go" \
		-e "DOCKER_PROTO=yes" \
		-e "PROTO_USER=$(shell id -u)" \
		-e "PROTO_GROUP=$(shell id -g)" \
		-e "PATH=/bin:/usr/bin:/usr/local/bin:/go/bin:/usr/local/go/bin" \
		quay.io/openstorage/grpc-framework:$(TAG)\
			grpcfw-lint \
				-o $(PROTO_FILE).lint \
				--disable-rule core::0191::proto-package \
				--set-exit-status \
				$(PROTO_FILE)

docker-proto:
ifndef DOCKER_PROTO
	$(error Do not run directly. Run 'make proto' instead.)
endif
	grpcfw $(PROTO_FILE)
	grpcfw-rest $(PROTO_FILE)
	grpcfw-doc $(PROTO_FILE)


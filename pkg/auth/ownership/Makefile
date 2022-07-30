PROTO_FILE = ownership.proto

all: proto

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
	rm -f ownership.swagger.json

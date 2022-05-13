
PROTO_FILE = ./api/hello.proto

PROTOC = protoc
PROTOS_PATH = $(GOPATH)/src
PROTOSRC_PATH = .

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
		-e "PROTO_USER=$(id -u)" \
		-e "PROTO_GROUP=$(id -g)" \
		-e "PATH=/bin:/usr/bin:/usr/local/bin:/go/bin:/usr/local/go/bin" \
		quay.io/openstorage/grpc-framework \
			make docker-proto

docker-proto:
ifndef DOCKER_PROTO
	$(error Do not run directly. Run 'make proto' instead.)
endif
	@echo ">>> Generating protobuf definitions"
	$(PROTOC) -I $(PROTOSRC_PATH) \
		-I /usr/local/include \
		-I $(PROTOS_PATH)/github.com/grpc-ecosystem/grpc-gateway \
		-I $(PROTOS_PATH)/github.com/googleapis/googleapis \
		--go_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative \
		$(PROTO_FILE)
	$(PROTOC) -I $(PROTOSRC_PATH) \
		-I /usr/local/include \
		-I $(PROTOS_PATH)/github.com/grpc-ecosystem/grpc-gateway \
		-I $(PROTOS_PATH)/github.com/googleapis/googleapis \
		--grpc-gateway_out=logtostderr=true:. \
		$(PROTO_FILE)
	$(PROTOC) -I $(PROTOSRC_PATH) \
		-I /usr/local/include \
		-I $(PROTOS_PATH)/github.com/grpc-ecosystem/grpc-gateway \
		-I $(PROTOS_PATH)/github.com/googleapis/googleapis \
		--openapiv2_out=logtostderr=true:. \
		$(PROTO_FILE)
	@echo ">>> Upgrading swagger 2.0 to openapi 3.0"
	mv api/hello.swagger.json api/20hello.swagger.json
	swagger2openapi api/20hello.swagger.json -o api/hello.swagger.json
	rm -f api/20hello.swagger.json

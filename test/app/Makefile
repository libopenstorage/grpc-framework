PROTOC = protoc
PROTOS_PATH = $(GOPATH)/src
PROTOSRC_PATH = .

all: proto

proto:
	docker run \
		--privileged --rm \
		-v $(shell pwd):/go/src/code \
		-e "GOPATH=/go" \
		-e "DOCKER_PROTO=yes" \
		-e "PATH=/bin:/usr/bin:/usr/local/bin:/go/bin:/usr/local/go/bin" \
		quay.io/openstorage/grpc-framework \
			make docker-proto

docker-proto:
ifndef DOCKER_PROTO
	$(error Do not run directly. Run 'make proto' instead.)
endif
	@echo ">>> Generating protobuf definitions from api/api.proto"
	$(PROTOC) -I $(PROTOSRC_PATH) \
		-I /usr/local/include \
		-I $(PROTOS_PATH)/github.com/grpc-ecosystem/grpc-gateway \
		-I $(PROTOS_PATH)/github.com/googleapis/googleapis \
		--go-grpc_out=hello \
		./hello.proto
	$(PROTOC) -I $(PROTOSRC_PATH) \
		-I /usr/local/include \
		-I $(PROTOS_PATH)/github.com/grpc-ecosystem/grpc-gateway \
		-I $(PROTOS_PATH)/github.com/googleapis/googleapis \
		--grpc-gateway_out=logtostderr=true:. \
		./hello.proto
	$(PROTOC) -I $(PROTOSRC_PATH) \
		-I /usr/local/include \
		-I $(PROTOS_PATH)/github.com/grpc-ecosystem/grpc-gateway \
		-I $(PROTOS_PATH)/github.com/googleapis/googleapis \
		--swagger_out=logtostderr=true:. \
		./hello.proto
	#@echo ">>> Upgrading swagger 2.0 to openapi 3.0"
	#mv api/server/sdk/api/api.swagger.json api/server/sdk/api/20api.swagger.json
	#swagger2openapi api/server/sdk/api/20api.swagger.json -o api/server/sdk/api/api.swagger.json
	#rm -f api/server/sdk/api/20api.swagger.json

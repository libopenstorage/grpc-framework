
PROTO_FILE = ./api/hello.proto
TAG = dev

all: protos server client

./bin:
	mkdir ./bin

.PHONY: server
server: ./bin
	go build -o ./bin/server cmd/server.go

.PHONY: client
client: ./bin
	go build -o ./bin/client cmd/client/client.go

.PHONY: clean
clean:
	rm -rf ./bin

.PHONY: protos
protos:
	$(MAKE) -C protos

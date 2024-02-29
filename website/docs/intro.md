# Introduction

The following will guide through some of the features provided 
by this framework

## Containerize tools
The framework provides the latest tools for generating gRPC as a single
container called `quay.io/openstorage/grpc-framework`. Here is a
sample set of targets for a `Makefile` on how the container can be
used to generate gRPC code from your protocol buffer files.

```Makefile
PROTO_FILE = ./api/hello.proto

proto:
	docker run \
		--privileged --rm \
		-v $(shell pwd):/go/src/code \
		-e "GOPATH=/go" \
		-e "DOCKER_PROTO=yes" \
		-e "PROTO_USER=$(shell id -u)" \
		-e "PROTO_GROUP=$(shell id -g)" \
		-e "PATH=/bin:/usr/bin:/usr/local/bin:/go/bin:/usr/local/go/bin" \
		quay.io/openstorage/grpc-framework:v0.1.2 \
			make docker-proto

docker-proto:
ifndef DOCKER_PROTO
	$(error Do not run directly. Run 'make proto' instead.)
endif
	grpcfw $(PROTO_FILE)
	grpcfw-rest $(PROTO_FILE)
	grpcfw-doc $(PROTO_FILE)
```

The framework provides a set of `grpcfw*` programs in a container,
to help with the generation of the sources from the protocol buffers file.

## Generate REST and swagger APIs
The framework utilizes the [grpc-gateway] to generate a REST interface
for your application. The framework will setup and start the HTTP server
and automatically connect it to your gRPC server. The REST APIs will be
served by the HTTP server which are then forwared to the gRPC server
to be handled.

The framework will also generate a [swagger] API file which can be
provided to REST client developers.

[grpc-gateway]: https://grpc-ecosystem.github.io/grpc-gateway/
[swagger]: https://swagger.io/

The following is an example of how a gRPC service can be used from a REST
client. 

```proto
// The greeting service definition.
service HelloGreeter {
  // Sends a greeting
  rpc SayHello (HelloGreeterSayHelloRequest)
    returns (HelloGreeterSayHelloResponse) {
      option(google.api.http) = {
        post: "/v1/greeter/sayhello"
        body: "*"
      };
  }
}

// The request message containing the user's name.
message HelloGreeterSayHelloRequest {
  string name = 1;
}

// The response message containing the greetings
message HelloGreeterSayHelloResponse {
  string message = 1;
}
```

Here, the _grpc-gateway_ would use the information located in the _option(google.api.http)_
section of your gRPC rpc to generate a REST call. To enable the REST server
for your application, you would need to enable it in your configuration:

```go
import(
	"github.com/libopenstorage/grpc-framework/server"
)
...
	grpcConfig := &server.ServerConfig{
		Name:         "hello",
		Address:      "127.0.0.1:9009",
		Socket:       "/tmp/hello-server.sock",
		AuditOutput:  os.Stdout,
		AccessOutput: os.Stdout,
	}
    restPort := "9010"
	grpcConfig.WithDefaultRestServer(restPort)
...
```

Once the server is started, you can then use a REST client to send commands
to your application:

```bash
$ curl  -X POST -d '{ "name": "Luis" }' \
    --silent http://localhost:9010/v1/greeter/sayhello | jq
{
  "message": "Hello, Luis"
}
```

## Generate Markdown documentation
The framework also provides [protoc-doc] to generate Markdown documentation from
the comments on your protocol buffers files. 

## Security
The framework makes it simple to add authentication, authorization, and TLS
to secure your application. Authentication and RBAC authorization are provided
by a set of interceptors in the gRPC server. 

### Authentication (OIDC and JWT supported)
The framework support shared secret, public-private key, or OpenID Connect
authentication models.

### Authorization
The framework provides (RBAC) role based access control for gRPC services as
well as a generic resource authorization model.

### TLS support
The framework provides TLS server support.

### Auditing
The framework logs access to the APIs by recording identifying information read
from the authentication token of the caller. This is done by an interceptor
that is automatically installed by the framework when authentication
is enabled on the gRPC server.

## Rate Limiter
The framework provides rate limiter support with a plan for future releases
tor provide the rate limit per user.

## Metrics for Prometheus
Support for Prometheus is provided by [go-grpc-prometheus].

## API Logging
Like auditing, a logging interceptor is provided which can provide rountrip
information for API services.

## proto/gRPC build container
All tools and updated libraries are all provided by a container to make it
simple to utilize on your projects.

[protoc-doc]: https://github.com/pseudomuto/protoc-gen-doc
[go-grpc-prometheus]: https://github.com/grpc-ecosystem/go-grpc-prometheus


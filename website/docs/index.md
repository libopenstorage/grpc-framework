# Welcome to grpc-framework

The grpc-framework enables Golang developers to create secure gRPC applications
easily. The project provides developers with the following features:

- Generate REST and swagger APIs
- Generate Markdown documentation
- Security
    - Authentication (OIDC and JWT supported)
    - Authorization (RBAC / Role Based Access Control)
    - TLS support
    - Auditing
- Rate Limiter
- Metrics for Prometheus
- API Logging
- proto/gRPC build container
- And more...

## Usage

To add the library to your Golang application use the following command:

```bash
go get github.com/libopenstorage/grpc-framework@v0.0.5
```

We are working on a tutorial, but in the meantime, please check out
the example [test program].

[test program]: https://github.com/libopenstorage/grpc-framework/tree/master/test/app

## Projects Used

grpc-framework uses the following excellent projects in the framework:

* [gRPC Golang](https://grpc.io/docs/languages/go/basics/)
* [gRPC REST Gateway](https://grpc-ecosystem.github.io/grpc-gateway/)
* [Golang JWT](https://github.com/golang-jwt/jwt)
* [Logging with logrus](https://github.com/sirupsen/logrus)
* [Generate Markdown documentation](https://github.com/pseudomuto/protoc-gen-doc)
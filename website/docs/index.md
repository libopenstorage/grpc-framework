# Welcome to grpc-framework

The grpc-framework enables Golang developers to create secure gRPC applications
easily. The project provides developers with the following features:

- Security
    - Authentication (OIDC and JWT supported)
    - Authorization (RBAC / Role Based Access Control)
    - TLS support
    - Auditing
- proto/gRPC build container
- Rate Limiter
- Metrics for Prometheus
- API Logging
- And more...
## Dependencies

grpc-framework uses the following excellent projects in the framework:

* [gRPC Golang](https://grpc.io/docs/languages/go/basics/)
* [gRPC REST Gateway](https://grpc-ecosystem.github.io/grpc-gateway/)
* [Golang JWT](https://github.com/golang-jwt/jwt)
* [Logging with logrus](https://github.com/sirupsen/logrus)
* [Generate Markdown documentation](https://github.com/pseudomuto/protoc-gen-doc)
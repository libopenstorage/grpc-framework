# Tutorial

This tutorial will guide you to create your sample application on your
system.

You will need [Docker](https://docker.com) (or complient runtime) and [Go](https://go.dev) installed

## Container Runtime

### Linux

First, make sure to have a container runtime installed like [Docker](https://docker.com)
or [podman](https://podman.io).

### MacOS

If you use a MacOS system, then it is recommended to
use [Docker Desktop](https://www.docker.com/products/docker-desktop/).

### Windows

If you use Windows, it is highly recommended to install
[WSL](https://learn.microsoft.com/en-us/windows/wsl/install)
and [Docker Desktop](https://www.docker.com/products/docker-desktop/).

## Sample Application

Open a command prompt and make a new directory:

```
mkdir hello
cd hello
```

Populate with the sample application from the grpc-framework:

```
curl -L \
  https://github.com/libopenstorage/grpc-framework/archive/refs/heads/master.tar.gz | \
  tar xz --strip=3 "grpc-framework-master/test/app"
```

Run `go mod init`. See [Golang: Getting Started](https://go.dev/doc/tutorial/getting-started) for
more information:

```
go mod init hello
```

Now add the grpc-framework as a dependency:

```
go get github.com/libopenstorage/grpc-framework@v0.1.3
```

Let golang determine the rest of the dependencies:

```
go mod tidy
```

Build:

```
make
```

Run the server:

```
./bin/server
```


On another terminal run the client:

```
./bin/client
```

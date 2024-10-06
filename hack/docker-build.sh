#!/bin/bash

function fail()
{
    echo $@ >&2
    exit 1
}

function RUN()
{
     eval $@ || fail "Failed to run [$@]: $?"
}

# output: ${PLATFORM}
PLATFORM=$(uname -m)
if [ "$PLATFORM" = "x86_64" ] ; then
    ARCH="amd64"
elif [ "$PLATFORM" = "aarch64" ] ; then
	ARCH="arm64"
else
    echo "Unknown platform: $PLATFORM"
    exit 1
fi

## VERSIONS
## Confirm that the links are correct. Some tools change the links on newer versions
# https://go.dev/dl/
GFGOLANG=1.23.2
# https://github.com/grpc-ecosystem/grpc-gateway/releases
GFGRPCGATEWAY=2.22.0
# https://github.com/pseudomuto/protoc-gen-doc/releases
GFPROTOCGENDOC=1.5.1
# https://github.com/protocolbuffers/protobuf/releases
GFPROTOC=28.2

# Get gRPC golang versions from here: https://grpc.io/docs/languages/go/quickstart/
# Also see: https://pkg.go.dev/google.golang.org/protobuf/cmd/protoc-gen-go
GFPROTOCGENGO=1.34.2
# Also see: https://pkg.go.dev/google.golang.org/grpc/cmd/protoc-gen-go-grpc
GFPROTOCGENGOGRPC=1.5.1

# Install tools from Ubuntu
RUN apt-get -y -qq install \
	make \
	nodejs \
	npm \
	wget \
	curl \
	python3 \
	python3-venv \
	unzip \
	git && \
	apt-get clean && \
	apt-get autoclean

# Install latest golang
RUN rm -rf /usr/local/go
RUN wget -nv https://dl.google.com/go/go${GFGOLANG}.linux-${ARCH}.tar.gz && \
	tar -xf go${GFGOLANG}.linux-${ARCH}.tar.gz && mv go /usr/local
RUN rm -f go${GFGOLANG}.linux-${ARCH}.tar.gz


# Install protoc
if [ "${PLATFORM}" = "aarch64" ] ; then
	PROTOCPLATFORM="aarch_64"
else
	PROTOCPLATFORM=$PLATFORM
fi
RUN curl -LO https://github.com/protocolbuffers/protobuf/releases/download/v${GFPROTOC}/protoc-${GFPROTOC}-linux-${PROTOCPLATFORM}.zip
RUN unzip protoc-${GFPROTOC}-linux-${PROTOCPLATFORM}.zip -d /usr/local
RUN rm -f protoc-${GFPROTOC}-linux-${PROTOCPLATFORM}.zip

##
## gRPC Gateway
##
RUN go install \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v${GFGRPCGATEWAY} \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v${GFGRPCGATEWAY}
# Get the proto files from the grpc-gateway
RUN mkdir -p /go/src/github.com/grpc-ecosystem && \
	cd /go/src/github.com/grpc-ecosystem && \
	git clone -b v${GFGRPCGATEWAY} https://github.com/grpc-ecosystem/grpc-gateway.git
# Install swagger 2.0 to OpenApi 3.0 converter
RUN npm install -g swagger2openapi

##
## protobuf and golang gRPC compilers
##
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v${GFPROTOCGENGO}
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v${GFPROTOCGENGOGRPC}
# Install Google Api proto files
RUN mkdir -p /go/src/github.com/googleapis && \
	cd /go/src/github.com/googleapis && \
	git clone https://github.com/googleapis/googleapis.git

##
## proto-gen-doc
##
RUN wget https://github.com/pseudomuto/protoc-gen-doc/releases/download/v${GFPROTOCGENDOC}/protoc-gen-doc_${GFPROTOCGENDOC}_linux_${ARCH}.tar.gz && \
	tar xzvf protoc-gen-doc_${GFPROTOCGENDOC}_linux_${ARCH}.tar.gz && \
	mv protoc-gen-doc /usr/local/bin
RUN rm -f protoc-gen-doc_${GFPROTOCGENDOC}_linux_${ARCH}.tar.gz

##
## Install Google AIP api-linter
##
RUN go install github.com/googleapis/api-linter/cmd/api-linter@latest

##
## Cleanup
##
RUN rm -rf /usr/share/doc
RUN rm -rf /go/pkg/*
RUN apt remove -y unzip $(dpkg -l | grep X11 | awk '{print $2}' | cut -d: -f1)
RUN apt autoremove -y


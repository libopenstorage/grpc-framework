FROM golang
LABEL org.opencontainers.image.authors="lpabon@purestorage.com"

ENV GOPATH=/go
RUN apt update

##
## grpc-framework additions
##
# Install tools
COPY ./tools/grpcfw* /usr/local/bin/
# Add protofiles
RUN mkdir -p /go/src/github.com/libopenstorage/grpc-framework
COPY . /go/src/github.com/libopenstorage/grpc-framework

##
## Install software specific to this arch
##
RUN bash /go/src/github.com/libopenstorage/grpc-framework/hack/docker-build.sh

##
## Set working directory
##
RUN mkdir -p /go/src/code
WORKDIR /go/src/code

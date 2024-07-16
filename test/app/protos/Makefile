export TOP=$(shell pwd)

ifndef L
L=go
endif

ifndef LINT_OUTPUT
LINT_OUTPUT=false
endif
DOCKERCMD=docker run \
		--privileged --rm \
		-v $(shell pwd):/go/src/code \
		-e "LINT_OUTPUT=$(LINT_OUTPUT)" \
		-e "GOPATH=/go" \
		-e "L=$(L)" \
		-e "DOCKER_PROTO=yes" \
		-e "PROTO_USER=$(shell id -u)" \
		-e "PROTO_GROUP=$(shell id -g)" \
		-e "PATH=/bin:/usr/bin:/usr/local/bin:/go/bin:/usr/local/go/bin" \
		$(GRPC_FRAMEWORK_CONTAINER)

SUBDIRS=apis

all: build docs

python:
	$(DOCKERCMD) make -C sdk/python setup
	$(MAKE) L=python build
	$(DOCKERCMD) make -C sdk/python sdk

.PHONY: publish
publish: build docs python
	$(DOCKERCMD) make -C publish

docker-shell:
	docker run \
		-ti \
		--privileged --rm \
		-v $(shell pwd):/go/src/code \
		-e "LINT_OUTPUT=$(LINT_OUTPUT)" \
		-e "GOPATH=/go" \
		-e "L=$(L)" \
		-e "DOCKER_PROTO=yes" \
		-e "PROTO_USER=$(shell id -u)" \
		-e "PROTO_GROUP=$(shell id -g)" \
		-e "PATH=/bin:/usr/bin:/usr/local/bin:/go/bin:/usr/local/go/bin" \
		$(GRPC_FRAMEWORK_CONTAINER) /bin/bash

ifndef DOCKER_PROTO
build:
	$(DOCKERCMD) make build
else
.PHONY: build
build: lint $(SUBDIRS)
endif

.PHONY: lint
lint:
	@SCRIPTSDIR=$(TOP)/lint/repo-scripts $(TOP)/lint/run.sh

include $(TOP)/make/common.mk

.PHONY: docs
ifndef DOCKER_PROTO
docs:
	$(DOCKERCMD) make docs
else
docs:
	$(MAKE) -C website
endif

.PHONY: serve
serve:
	cd docs ; python3 -m http.server

.PHONY: clean
clean:
	$(MAKE) -C website clean


TAG := dev
HAS_ERRCHECK := $(shell command -v errcheck 2> /dev/null)
PKGS := $(shell go list ./... | grep -v vendor | grep -v examples)

all: build

build:
	@echo ">>> go build"
	go build $(PKGS)

fmt:
	@echo ">>> go fmt"
	@echo "-- ignoring fmt checks due to golang 1.19 changes"
	-go fmt $(PKGS) | wc -l | xargs | grep "^0"

vet:
	@echo ">>> go vet"
	@go vet $(PKGS)

errcheck:
ifndef HAS_ERRCHECK
	-GO111MODULE=off go get -u github.com/kisielk/errcheck
endif
	@echo ">>> errcheck"
	errcheck $(PKGS)

pr-verify:
	git-validation -run DCO,short-subject
	go mod vendor && git grep -rw GPL vendor | grep LICENSE | egrep -v "yaml.v2" | wc -l | grep "^0"

test: build
	@echo ">>> go test"
	go test $(PKGS)

testapp:
	$(MAKE) -C test/app

verify: vet fmt test testapp

travis-verify: pr-verify verify

# Run this after creating and pushing a release tag into the repo
go-mod-publish:
	GOPROXY=proxy.golang.org go list -m github.com/libopenstorage/grpc-framework@$(shell git describe --tags)

proto:
	$(MAKE) -C pkg proto

clean:
	$(MAKE) clean -C test/app

container:
	docker build -f Dockerfile.proto -t quay.io/openstorage/grpc-framework:$(TAG) .

./venv:
	python3 -m venv venv
	bash -c "source venv/bin/activate && \
		pip3 install -r requirements.txt"
	@echo "Type: 'source venv/bin/active' to get access to mkdocs"

doc-env: ./venv

doc-build: doc-env
	bash -c "source venv/bin/activate && \
		cd website && \
		mkdocs build"

doc-serve: doc-env
	bash -c "source venv/bin/activate && \
		cd website && \
		mkdocs serve"

.PHONY: clean proto go-mod-publish travis-verify verify \
	testapp test pr-verify errcheck vet fmt build \
	doc-env doc-build doc-serve container


HAS_ERRCHECK := $(shell command -v errcheck 2> /dev/null)
PKGS := $(shell go list ./... | grep -v vendor | grep -v examples)

all: build

build:
	@echo ">>> go build"
	go build $(PKGS)

fmt:
	@echo ">>> go fmt"
	go fmt $(PKGS) | wc -l | xargs | grep "^0"

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

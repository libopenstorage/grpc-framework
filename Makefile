
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
	$(MAKE) -C test/app proto-lint
	$(MAKE) -C test/app

testapp-verify: testapp
	./hack/client-server-test.sh

verify: vet fmt test testapp-verify

travis-verify: pr-verify verify

# Run this after creating and pushing a release tag into the repo
go-mod-publish:
	GOPROXY=proxy.golang.org go list -m github.com/libopenstorage/grpc-framework@$(shell git describe --tags)

proto:
	$(MAKE) -C pkg proto

clean:
	$(MAKE) clean -C test/app

container:
	docker build -t quay.io/openstorage/grpc-framework:$(TAG) .

container-buildx-install:
	@echo "Setting up multiarch emulation"
	docker run --privileged --rm tonistiigi/binfmt --install all
	@echo "Setting up multiarch builder"
	docker buildx create --name gfwbuilder --driver docker-container --bootstrap
	docker buildx use gfwbuilder

# Run: make container-buildx-install first to install the emulation
container-release:
	@echo "This will automatically push. Must be logged in to quay.io"
	docker buildx build \
		--push \
		--platform linux/amd64,linux/arm64  \
		--tag quay.io/openstorage/grpc-framework:$(TAG) .

container-buildx-uninstall:
	docker buildx stop gfwbuilder
	docker buildx rm gfwbuilder

./venv:
	python3 -m venv venv
	bash -c "source venv/bin/activate && \
		pip3 install --upgrade pip && \
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
	doc-env doc-build doc-serve container testapp-verify \
	container-buildx-install container-release container-buildx-uninstall


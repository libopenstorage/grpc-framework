
HAS_ERRCHECK := $(shell command -v errcheck 2> /dev/null)
PKGS := $(shell go list ./... | grep -v vendor | grep -v examples)

all:

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

test:
	@echo ">>> go test"
	go test $(PKGS)

verify: vet fmt test

travis-verify: pr-verify verify


HAS_ERRCHECK := $(shell command -v errcheck 2> /dev/null)
PKGS := $(shell go list ./... | grep -v vendor | grep -v examples)

all:

fmt:
	go fmt $(PKGS) | grep -v "api.pb.go" | wc -l | grep "^0";

vet:
	go vet $(PKGS)

errcheck:
ifndef HAS_ERRCHECK
	-GO111MODULE=off go get -u github.com/kisielk/errcheck
endif
	errcheck $(PKGS)

pr-verify:
	git-validation -run DCO,short-subject
	go mod vendor && git grep -rw GPL vendor | grep LICENSE | egrep -v "yaml.v2" | wc -l | grep "^0"

test:
	go test $(PKGS)

verify: vet fmt test

travis-verify: pr-verify verify

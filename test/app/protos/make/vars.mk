
VERSION = $(shell git config --global --add safe.directory /go/src/code && git describe --always --tags)
PYTHON_VERSION = $(shell echo $(VERSION) | sed -e 's/-[[:digit:]]\+-g/+/')


TOOLS=$(TOP)/tools
PWD=$(shell pwd)
GRPC_FRAMEWORK_TAG=v0.1.3
GRPC_FRAMEWORK_CONTAINER=quay.io/openstorage/grpc-framework:$(GRPC_FRAMEWORK_TAG)

ifndef LINTRULES
LINTRULES=$(TOP)/lint/rules.yml
endif

ifdef SUBDIRS
.PHONY: $(SUBDIRS)
$(SUBDIRS):
	$(MAKE) -C $@
endif

ifdef PROTO_FILES
PROTO_PATH:=$(PWD:$(TOP)/%=%)
BUILD_PROTO_FILES:=$(addprefix $(PROTO_PATH)/, $(PROTO_FILES))

# DEBUG
# $(info TOP = $(TOP) PWD = $(PWD) PROTO_PATH = $(PROTO_PATH))
# $(info BUILD_PROTO_FILES = $(BUILD_PROTO_FILES))

.PHONY: $(PROTO_FILES)
$(PROTO_FILES): $(BUILD_PROTO_FILES)

.PHONY: $(BUILD_PROTO_FILES)
$(BUILD_PROTO_FILES): lint-scripts

ifndef NOBUILD
ifeq ($(L),go)
	@echo ">>> Building golang code for $@"
	cd $(TOP) && $(TOOLS)/grpcfw $@

	@echo ">>> Building golang REST gateway for $@"
	cd $(TOP) && $(TOOLS)/grpcfw-rest $@
endif
endif # NOBUILD

ifndef NODOC
	@echo ">>> Building docs $@"
	cd $(TOP) && $(TOOLS)/grpcfw-doc $@
endif # NODOC

ifndef NOLINT
	@echo ">>> Linting $@"
ifeq ($(LINT_OUTPUT),true)
	cd $(TOP) && $(TOOLS)/grpcfw-lint --config $(LINTRULES) \
				--set-exit-status \
				$@
else
	cd $(TOP) && $(TOOLS)/grpcfw-lint --config $(LINTRULES) \
				--set-exit-status \
				--output-path=$@.lint \
				$@
endif # LINT_OUTPUT
endif # NOLINT

endif # PROTO_FILES

.PHONY: lint-scripts
lint-scripts:
	@SCRIPTSDIR=$(TOP)/lint/proto-scripts $(TOP)/lint/run.sh


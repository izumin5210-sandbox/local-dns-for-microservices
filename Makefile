.DEFAULT_GOAL := all

#  dep
#----------------------------------------------------------------
DEP_BIN_DIR := ./vendor/.bin/
DEP_SRCS := \
	github.com/mattn/goreman

DEP_BINS := $(addprefix $(DEP_BIN_DIR),$(notdir $(DEP_SRCS)))

define dep-bin-tmpl
$(eval OUT := $(addprefix $(DEP_BIN_DIR),$(notdir $(1))))
$(OUT): dep
	@cd vendor/$(1) && GOBIN="$(shell pwd)/$(DEP_BIN_DIR)" go install .
endef

$(foreach src,$(DEP_SRCS),$(eval $(call dep-bin-tmpl,$(src))))


#  app
#----------------------------------------------------------------
.PHONY: all
all: bin/discoverer bin/server

.PHONY: run
run: all
	@$(DEP_BIN_DIR)/goreman start

bin/discoverer: cmd/discoverer/*.go
	@go build -v -o bin/discoverer cmd/discoverer/*.go

bin/server: cmd/server/*.go
	@go build -v -o bin/server cmd/server/*.go

.PHONY: setup
setup: dep $(DEP_BINS)

.PHONY: dep
dep: Gopkg.toml Gopkg.lock
	@dep ensure -v -vendor-only

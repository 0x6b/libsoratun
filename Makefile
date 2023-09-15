NAME = soratun
GO = go
LIB_DIR = lib

SRC = libsoratun.go

LIB_SHARED = $(LIB_DIR)/shared/lib$(NAME).so
LIB_ARCHIVE = $(LIB_DIR)/archive/lib$(NAME).a

BINDING_DIR_RUST = rust
BINDING_RUST = $(BINDING_DIR_RUST)/src/$(NAME).rs

help:
	@echo "usage: make <\033[36mtarget\033[0m>"
	@echo
	@echo "available targets:"
	@grep -E '^[a-zA-Z._-]+:.*?## .*$$' Makefile | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2}'

all: libs bindings ## Build all

libs: $(LIB_SHARED) $(LIB_ARCHIVE) ## Build libraries

$(LIB_SHARED): $(SRC) go.mod ## Build shared library
	go build -buildmode=c-shared -o $(LIB_DIR)/shared/lib$(NAME).so $(SRC)

$(LIB_ARCHIVE): $(SRC) go.mod ## Build archive library
	go build -buildmode=c-archive -o $(LIB_DIR)/archive/lib$(NAME).a $(SRC)

bindings: $(BINDING_RUST) ## Build bindings

$(BINDING_RUST): $(LIB_ARCHIVE) ## Build Rust bindings
	bindgen --no-layout-tests $(LIB_DIR)/archive/lib$(NAME).h -o $(BINDING_DIR_RUST)/src/$(NAME).rs

clean: ## Clean up
	rm -rf $(LIB_DIR)/{archive,shared}/*

.PHONY: help all libs bindings

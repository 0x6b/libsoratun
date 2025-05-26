NAME = soratun
GO = go

SRC = $(wildcard libsoratun/*.go) $(LIB_ENTRY)
LIB_DIR = lib

LIB_ENTRY = libsoratun.go

LIB_SHARED = $(LIB_DIR)/shared/lib$(NAME).so
LIB_ARCHIVE = $(LIB_DIR)/archive/lib$(NAME).a

BINDING_DIR_RUST = examples/rust
BINDING_RUST = $(BINDING_DIR_RUST)/src/$(NAME).rs

LDFLAGS = -ldflags="-X 'github.com/0x6b/libsoratun/libsoratun.Revision=$(shell git rev-parse --short HEAD)'"

help:
	@echo "usage: make <\033[36mtarget\033[0m>"
	@echo
	@echo "available targets:"
	@grep -E '^[a-zA-Z._-]+:.*?## .*$$' Makefile | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2}'

libs: $(LIB_SHARED) $(LIB_ARCHIVE) ## Build libraries

$(LIB_SHARED): $(SRC) go.mod ## Build shared library
	go build -buildmode=c-shared $(LDFLAGS) -o $(LIB_DIR)/shared/lib$(NAME).so $(LIB_ENTRY)
	go build -buildmode=c-shared $(LDFLAGS) -o $(LIB_DIR)/shared/lib$(NAME).dylib $(LIB_ENTRY)
	go build -buildmode=c-shared $(LDFLAGS) -o $(LIB_DIR)/shared/lib$(NAME).dll $(LIB_ENTRY)
ifeq ($(shell uname -s),Linux)
	GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CXX=x86_64-w64-mingw32-g++ CC=x86_64-w64-mingw32-gcc go build -buildmode=c-shared $(LDFLAGS) -o $(LIB_DIR)/shared/lib$(NAME).dll $(LIB_ENTRY)
endif

$(LIB_ARCHIVE): $(SRC) go.mod ## Build archive library
	go build -buildmode=c-archive $(LDFLAGS) -o $(LIB_DIR)/archive/lib$(NAME).a $(LIB_ENTRY)

bindings: $(BINDING_RUST) ## Build Rust binding

$(BINDING_RUST): $(LIB_ARCHIVE) ## Build Rust bindings
	bindgen --no-layout-tests $(LIB_DIR)/archive/lib$(NAME).h -o $(BINDING_DIR_RUST)/src/$(NAME).rs

all: libs bindings ## Build libraries and Rust binding

clean: ## Clean up
	rm -rf $(LIB_DIR)/archive/*
	rm -rf $(LIB_DIR)/shared/*

.PHONY: help all libs bindings

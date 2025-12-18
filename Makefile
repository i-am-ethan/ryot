GOCACHE ?= $(CURDIR)/.gocache

BIN_DIR := $(CURDIR)/bin
RYOT := $(BIN_DIR)/ryot

.PHONY: all build install clean

all: build

build:
	mkdir -p $(BIN_DIR)
	GOCACHE=$(GOCACHE) go build -o $(RYOT) ./cmd/ryot

install: build
	install $(RYOT) $(HOME)/bin/

clean:
	rm -rf $(BIN_DIR) .gocache .dircache temp_git_file_*

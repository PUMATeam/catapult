# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

PREFIX=.
ARTIFACT_DIR ?= .

VERSION?=$(shell git describe --tags --always --match "v[0-9]*" | awk -F'-' '{print substr($$1,2) }')
RELEASE?=$(shell git describe --tags --always --match "v[0-9]*" | awk -F'-' '{if ($$2 != "") {print $$2 "." $$3} else {print 1}}')
VERSION_RELEASE=$(VERSION)$(if $(RELEASE),-$(RELEASE))

COMMIT=$(shell git rev-parse HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)

COMMON_ENV=CGO_ENABLED=0 GOOS=linux GOARCH=amd64
#COMMON_GO_BUILD_FLAGS=-ldflags '-extldflags "-static"'
COMMON_GO_BUILD_FLAGS=

TARBALL=catapult-$(VERSION_RELEASE).tar.gz

all: clean generate build test

binaries = \
	catapult \
	catapult-node

$(binaries): check test
	go vet ./cmd/$@ && \
	$(COMMON_ENV) $(GOBUILD) \
    	$(COMMON_GO_BUILD_FLAGS) \
    	-o $(PREFIX)/$@ \
    	-v cmd/$@/*.go

.PHONY: check
check:
	go vet ./...

.PHONY: build
build: $(binaries)

.PHONY: test
test:
	$(GOTEST) -v ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o /tmp/coverage.html

.PHONY: clean
clean:
	$(GOCLEAN)
	git clean -dfx -e .idea*

.PHONY: generate
generate:
	protoc -I pkg/node --go_out=plugins=grpc:pkg/node pkg/node/node.proto

.PHONY: tarball
tarball: $(TARBALL)

$(TARBALL):
	/bin/git archive --format=tar.gz HEAD > $(TARBALL)

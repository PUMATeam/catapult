# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY=catapult
BINARYPKG=./cmd

all: clean generate build test

.PHONY: check
check:
	go vet ./...

.PHONY: build
build: 
	$(GOBUILD) -o $(BINARY) $(BINARYPKG)
.PHONY: test
test:
	$(GOTEST) -v ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o /tmp/coverage.html

.PHONY: clean
clean:
	$(GOCLEAN) $(BINARYPKG)
	git clean -dfx -e .idea* -e *os.img -e db.toml

.PHONY: generate
generate:
	protoc -I pkg/node --go_out=plugins=grpc:pkg/node pkg/node/node.proto

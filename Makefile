BINARY_NAME=terraform-provider-slack
VERSION=0.0.1

.PHONY: all build test lint doc

all: build

build:
	go build -o $(BINARY_NAME)_v$(VERSION)

test: build
	TF_ACC=1 TF_CLI_CONFIG_FILE=$(CURDIR)/.terraformrc go test -v ./...

lint:
	go vet ./...

doc:
	go generate ./...

BUILD_DIR 		:= build
NAME := loxone-homekit
WORKDIR := $(dir $(realpath $(firstword $(MAKEFILE_LIST))))
.DEFAULT_GOAL := build

.PHONY: init
init:
	export GO111MODULE=on
	go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.17.1
	go get -u github.com/onsi/ginkgo/ginkgo
	go get -u github.com/modocache/gover
	go mod download

.PHONY: lint
lint:
	golangci-lint run --config golangci.yml

.PHONY: build
build:
	env GOOS=linux CGO_ENABLED=0 go build -o $(BUILD_DIR)/$(NAME) .

.PHONY: test
test:
	$(GOPATH)/bin/ginkgo -r --randomizeAllSpecs --randomizeSuites --failOnPending --cover --trace --progress --compilers=2

.PHONY: cover
cover:
	$(GOPATH)/bin/gover . coverage.txt

.PHONY: format
format:
	go fmt $(go list ./... | grep -v /vendor/ | grep -v generated)
	goimports -e -w -d $(shell find . -path "./vendor/*" -prune -o -path "./pkg/generated/*" -prune -o -type f -name '*.go' -print)

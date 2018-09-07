GO_PKG = gitlab.jetstack.net/jetstack/vault-unsealer

REGISTRY := quay.io/jetstack
IMAGE_NAME := vault-unsealer
BUILD_TAG := build
IMAGE_TAGS := canary

BINDIR ?= $(CURDIR)/bin

BUILD_IMAGE_NAME := golang:1.10.4

GOPATH ?= /tmp/go

CI_COMMIT_TAG ?= unknown
CI_COMMIT_SHA ?= unknown

UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Linux)
	SHASUM := sha256sum -c
	DEP_URL := https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64
	DEP_HASH := 31144e465e52ffbc0035248a10ddea61a09bf28b00784fd3fdd9882c8cbb2315
endif
ifeq ($(UNAME_S),Darwin)
	SHASUM := shasum -a 256 -c
	DEP_URL := https://github.com/golang/dep/releases/download/v0.4.1/dep-darwin-amd64
	DEP_HASH := 1544afdd4d543574ef8eabed343d683f7211202a65380f8b32035d07ce0c45ef
endif

help:
	# all 		- runs verify, build and docker_build targets
	# test 		- runs go_test target
	# build 	- runs go_build target
	# verify 	- verifies generated files & scripts

# Util targets
##############
.PHONY: all build verify

all: verify build docker_build

build: go_build

verify: depend verify_vendor go_verify

.builder_image:
	docker pull ${BUILD_IMAGE_NAME}

# Builder image targets
#######################
docker_%: .builder_image
	docker run -it \
		-v ${GOPATH}/src:/go/src \
		-v $(shell pwd):/go/src/${GO_PKG} \
		-w /go/src/${GO_PKG} \
		-e GOPATH=/go \
		${BUILD_IMAGE_NAME} \
		/bin/sh -c "make $*"

# Docker targets
################
docker_build:
	docker build -t $(REGISTRY)/$(IMAGE_NAME):$(BUILD_TAG) .

docker_push: docker_build
	set -e; \
		for tag in $(IMAGE_TAGS); do \
		docker tag $(REGISTRY)/$(IMAGE_NAME):$(BUILD_TAG) $(REGISTRY)/$(IMAGE_NAME):$${tag} ; \
		docker push $(REGISTRY)/$(IMAGE_NAME):$${tag}; \
	done

# Go targets
#################
go_verify: go_fmt go_vet go_test

verify_vendor:
	$(BINDIR)/dep ensure -no-vendor -dry-run -v

go_build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w -X main.version=$(CI_COMMIT_TAG) -X main.commit=$(CI_COMMIT_SHA) -X main.date=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)' -o vault-unsealer_linux_amd64

go_test:
	go test $$(go list ./... | grep -v '/vendor/')

go_fmt:
	@set -e; \
	GO_FMT=$$(git ls-files *.go | grep -v 'vendor/' | xargs gofmt -d); \
	if [ -n "$${GO_FMT}" ] ; then \
		echo "Please run go fmt"; \
		echo "$$GO_FMT"; \
		exit 1; \
	fi

go_vet:
	go vet $$(go list ./... | grep -v '/vendor/')

depend: $(BINDIR)/dep

$(BINDIR)/dep:
	mkdir -p $(BINDIR)
	curl -sL -o $@ $(DEP_URL)
	echo "$(DEP_HASH)  $@" | $(SHASUM)
	chmod +x $@


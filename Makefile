REGISTRY := quay.io/jetstack
IMAGE_NAME := vault-unsealer
BUILD_TAG := build
IMAGE_TAGS := canary

CI_COMMIT_TAG ?= unknown
CI_COMMIT_SHA ?= unknown

help:
	# all 		- runs verify, build and docker_build targets
	# test 		- runs go_test target
	# build 	- runs go_build target
	# docker  - build local version of vault-unsealer
	# release - build application release in container and push it to repo
	# verify 	- verifies generated files & scripts

# Util targets
##############
.PHONY: all build verify

all: verify build docker_build

test: go_test

build: go_build

docker: docker_build

release: verify docker_build docker_push

verify: go_verify

# Docker targets
################
docker_build:
	docker build -t $(REGISTRY)/$(IMAGE_NAME):$(BUILD_TAG) --build-arg CI_COMMIT_TAG=${CI_COMMIT_TAG} \
	--build-arg CI_COMMIT_SHA=${CI_COMMIT_SHA} --build-arg CI_DATE=$(shell date -u +%Y-%m-%dT%H:%M:%SZ) \
	-f Dockerfile .

docker_push: docker_build
	set -e; \
		for tag in $(IMAGE_TAGS); do \
		docker tag $(REGISTRY)/$(IMAGE_NAME):$(BUILD_TAG) $(REGISTRY)/$(IMAGE_NAME):$${tag} ; \
		docker push $(REGISTRY)/$(IMAGE_NAME):$${tag}; \
	done

# Go targets
#################
go_verify: go_fmt go_vet go_test

go_build:
	CGO_ENABLED=0 go build -a -tags netgo -ldflags '-w -X main.version=$(CI_COMMIT_TAG) \
		-X main.commit=$(CI_COMMIT_SHA) -X main.date=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)' \
		-o vault-unsealer

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

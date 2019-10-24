IMG_NAME = rafter-upload-service
IMG = $(DOCKER_PUSH_REPOSITORY)$(DOCKER_PUSH_DIRECTORY)/$(IMG_NAME)
TAG = $(DOCKER_TAG)

ROOT := ../..
BUILD_DIR := ${ROOT}/hack/ci/rafter-upload-service
COVERAGE_OUTPUT_PATH := ${BUILD_DIR}/cover.out
PROJECT_FILES := $(shell find ${ROOT} -type f -name "*.go" | egrep -v "\/vendor\/|_*/automock/|_*/testdata/|/pkg\/|_*export_test.go")
COVERAGE = $(go tool cover -func=${COVERAGE_OUTPUT_PATH} | grep total | awk '{print $3}'

RED=\033[0;31m
GREEN=\033[0;32m
INVERTED=\033[7m
NC=\033[0m	# No Color

GO_VET_CMD = $(shell go list ${ROOT}/... | grep -v /vendor/ | xargs -L1 go vet && echo $$? || echo $$?)
GO_FMT_CMD = $(shell echo ${PROJECT_FILES} | tr ' ' '\n' | xargs -L1 go fmt && echo $$? || echo $$?)

build-local: test build
.PHONY: build-local

clean:
	rm -rf ${BUILD_DIR}/**/*
.PHONY: clean

go-vet:
	@result=$(GO_VET_CMD);\
	if [ "$${result}" =	"0" ]; then echo "${GREEN}√\tgo vet${NC}";\
	else echo "${RED}✗\tgo vet${NC}"; exit 1;\
	fi
.PHONY: go-vet

go-fmt:
	@result=$(GO_FMT_CMD);\
	if [ "$${result}" =	"0" ]; then echo "${GREEN}√\tgo fmt${NC}";\
	else echo "${RED}✗\tgo fmt${NC}"; exit 1;\
	fi
.PHONY: go-fmt

test: clean go-vet go-fmt
	@result=$$(go test -short -coverprofile=${COVERAGE_OUTPUT_PATH} ${ROOT}/... 1>&2 && echo $$?);\
	if [ "$${result}" =	"0" ]; then echo "${GREEN}√\tgo test total statement coverage: \
		$$(go tool cover -func=${COVERAGE_OUTPUT_PATH} | grep total | awk '{print $3}' | tr -s '\t' | cut -f3) ${NC}";\
	else echo "${RED}✗\tgo test exit code: $${result}${NC}"; exit $${result};\
	fi
.PHONY: test

pull-licenses:
ifdef LICENSE_PULLER_PATH
	bash $(LICENSE_PULLER_PATH)
else
	mkdir -p ${ROOT}/licenses
endif
.PHONY: pull-licenses

build: pull-licenses
	go mod vendor
	docker build -t $(IMG_NAME) -f ${ROOT}/deploy/uploader/Dockerfile ${ROOT}
.PHONY: build

push-image:
	docker tag $(IMG_NAME) $(IMG):$(TAG)
	docker push $(IMG):$(TAG)
.PHONY: push-image

ci-pr: build push-image
.PHONY: ci-pr

ci-master: build push-image
.PHONY: ci-master

ci-release: build push-image
.PHONY: ci-master
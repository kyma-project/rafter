ROOT :=  $(shell pwd)
COVERAGE_OUTPUT_PATH := ${ROOT}/cover.out
LICENSES_PATH := ${ROOT}/licenses

# Image URL to use all building/pushing image targets
UPLOADER_IMG_NAME := rafter-upload-service
MANAGER_IMG_NAME := rafter-controller-manager
FRONTMATTER_IMG_NAME := rafter-frontmatter-service
ASYNCAPI_IMG_NAME := rafter-asyncapi-service

IMG-CI-NAME-PREFIX := $(DOCKER_PUSH_REPOSITORY)$(DOCKER_PUSH_DIRECTORY)

UPLOADER-CI-IMG-NAME := $(IMG-CI-NAME-PREFIX)/$(UPLOADER_IMG_NAME):$(DOCKER_TAG)
MANAGER-CI-IMG-NAME :=  $(IMG-CI-NAME-PREFIX)/$(MANAGER_IMG_NAME):$(DOCKER_TAG)
FRONTMATTER-CI-IMG-NAME := $(IMG-CI-NAME-PREFIX)/$(FRONTMATTER_IMG_NAME):$(DOCKER_TAG)
ASYNCAPI-CI-IMG-NAME :=  $(IMG-CI-NAME-PREFIX)/$(ASYNCAPI_IMG_NAME):$(DOCKER_TAG)

# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd:trivialVersions=true"

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

all: docker-build
.PHONY: all

build-uploader:
	docker build -t $(UPLOADER_IMG_NAME) -f ${ROOT}/deploy/uploader/Dockerfile ${ROOT}
.PHONY: build-uploader

push-uploader:
	docker tag $(UPLOADER_IMG_NAME) $(UPLOADER-CI-IMG-NAME)
	docker push $(UPLOADER-CI-IMG-NAME)
.PHONY: push-uploader

build-manager:
	docker build -t $(MANAGER_IMG_NAME) -f ${ROOT}/deploy/manager/Dockerfile ${ROOT}
.PHONY: build-manager

push-manager:
	docker tag $(MANAGER_IMG_NAME) $(MANAGER-CI-IMG-NAME)
	docker push $(MANAGER-CI-IMG-NAME)
.PHONY: push-manager

build-frontmatter:
	docker build -t $(FRONTMATTER_IMG_NAME) -f ${ROOT}/deploy/extension/frontmatter/Dockerfile ${ROOT}
.PHONY: build-frontmatter

push-frontmatter:
	docker tag $(FRONTMATTER_IMG_NAME) $(FRONTMATTER-CI-IMG-NAME)
	docker push $(FRONTMATTER-CI-IMG-NAME)
.PHONY: push-frontmatter

build-asyncapi:
	docker build -t $(ASYNCAPI_IMG_NAME) -f ${ROOT}/deploy/extension/asyncapi/Dockerfile ${ROOT}
.PHONY: build-asyncapi

push-asyncapi:
	docker tag $(ASYNCAPI_IMG_NAME) $(ASYNCAPI-CI-IMG-NAME)
	docker push $(ASYNCAPI-CI-IMG-NAME)
.PHONY: push-asyncapi

clean:
	rm -f ${COVERAGE_OUTPUT_PATH}
	rm -rf ${LICENSE_PATH}
.PHONY: clean

pull-licenses:
ifdef LICENSE_PULLER_PATH
	bash $(LICENSE_PULLER_PATH)
else
	mkdir -p ${LICENSE_PATH}
endif
.PHONY: pull-licenses

fmt:
	find ${ROOT} -type f -name "*.go" \
	| egrep -v '_*/automock|_*/testdata|_*export_test.go' \
	| xargs -L1 go fmt
.PHONY: fmt

vet:
	@go list ${ROOT}/... \
	| grep -v "automock" \
	| xargs -L1 go vet
.PHONY: vet

# Run tests
test: clean manifests vet fmt
	go test -short -coverprofile=${COVERAGE_OUTPUT_PATH} ${ROOT}/...
	@go tool cover -func=${COVERAGE_OUTPUT_PATH} \
		| grep total \
		| awk '{print "Total test coverage: " $$3}'
.PHONY: test

# Generate manifests e.g. CRD, RBAC etc.
manifests: controller-gen
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths="${ROOT}/..." \
		object:headerFile=${ROOT}/hack/boilerplate.go.txt \
		output:crd:artifacts:config=${ROOT}/config/crd/bases \
		output:rbac:artifacts:config=${ROOT}/config/rbac \
		output:webhook:artifacts:config=${ROOT}/config/webhook
.PHONY: manifests

docker-build: \
	test \
	pull-licenses \
	build-uploader \
	build-frontmatter \
	build-asyncapi \
	build-manager
.PHONY: docker-build

# Push the docker image
docker-push: \
	push-uploader \
	push-frontmatter \
	push-asyncapi \
	push-manager
.PHONY: docker-push

# find or download controller-gen
# download controller-gen if necessary
controller-gen:
ifeq (, $(shell which controller-gen))
	go get sigs.k8s.io/controller-tools/cmd/controller-gen@v0.2.0
CONTROLLER_GEN=$(GOBIN)/controller-gen
else
CONTROLLER_GEN=$(shell which controller-gen)
endif
.PHONY: controller-gen

ci-pr: docker-build docker-push
.PHONY: ci-pr

ci-master: docker-build docker-push
.PHONY: ci-master

ci-release: docker-build docker-push
.PHONY: ci-release
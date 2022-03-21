ifeq (private.mk,$(shell ls -1 private.mk 2>/dev/null))
include private.mk
endif
WATCH_NAMESPACE ?= $(shell kubectl config view --minify --output 'jsonpath={..namespace}')

IMG_NAME := freeipa-operator

ifneq (,$(shell command -v podman 2>/dev/null))
DOCKER := podman
else
ifneq (,$(shell command -v docker 2>/dev/null))
DOCKER := docker
else
DOCKER :=
endif
endif


# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN := $(shell go env GOPATH)/bin
else
GOBIN := $(shell go env GOBIN)
endif


# https://docs.github.com/en/free-pro-team@latest/actions/reference/environment-variables#default-environment-variables
ifneq (,$(GITHUB_SHA))
COMMIT_SHA := $(GITHUB_SHA)
endif
# https://docs.travis-ci.com/user/environment-variables/#default-environment-variables
ifneq (,$(TRAVIS_COMMIT))
COMMIT_SHA := $(TRAVIS_COMMIT)
endif
# https://docs.gitlab.com/ee/ci/variables/predefined_variables.html
ifneq (,$(CI_COMMIT_SHA))
COMMIT_SHA := $(CI_COMMIT_SHA)
endif

ifeq (,$(COMMIT_SHA))
COMMIT_SHA := $(shell git rev-parse HEAD)
endif

CONTAINER_IMAGE_FILE ?= $(IMG_NAME).tar

IMG_TAG := dev-$(COMMIT_SHA)
IMG_BASE ?= quay.io/freeipa

# Image URL to use all building/pushing image targets
IMG ?= $(IMG_BASE)/$(IMG_NAME):$(IMG_TAG)
# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd:crdVersions=v1"
TEMPLATES_PATH ?= $(PWD)/config/templates
SAMPLES_PATH ?= $(PWD)/config/samples
SAMPLE ?= ephemeral-storage
DEFAULT_STORAGE ?= ephemeral
CONFIG ?= default


# Install kind by:
# GO111MODULE="on" go get sigs.k8s.io/kind@v0.10.0
KIND_CLUSTER_NAME ?= idmcontroller
K8S_NODE_IMAGE ?= v1.19.0
PROMETHEUS_INSTANCE_NAME ?= prometheus-operator
# CONFIG_MAP_NAME ?= initcontainer-configmap
RELATED_IMAGE_FREEIPA ?= "quay.io/freeipa/freeipa-openshift-container:latest"

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN := $(shell go env GOPATH)/bin
else
GOBIN := $(shell go env GOBIN)
endif

export RELATED_IMAGE_FREEIPA

all: manager

# Empty rule to allow force other rules. The name of the rule should not
# match any file.
# https://www.gnu.org/software/make/manual/html_node/Force-Targets.html
.PHONY: FORCE
FORCE:

.PHONY: help
help:
	@cat HELP

# USE_EXISTING_CLUSTER={1,0}
# Run tests
ENVTEST_ASSETS_DIR:=$(shell pwd)/testbin
.PHONY: test
# idmocp-243 Workarounded by using the tool from v0.8.3 tag
test: generate fmt vet manifests
	mkdir -p "$(ENVTEST_ASSETS_DIR)"
	test -f "$(ENVTEST_ASSETS_DIR)/setup-envtest.sh" \
	|| curl -sSLo "$(ENVTEST_ASSETS_DIR)/setup-envtest.sh" "https://raw.githubusercontent.com/kubernetes-sigs/controller-runtime/v0.8.3/hack/setup-envtest.sh"
	source "$(ENVTEST_ASSETS_DIR)/setup-envtest.sh"; \
	fetch_envtest_tools "$(ENVTEST_ASSETS_DIR)"; \
	setup_envtest_env "$(ENVTEST_ASSETS_DIR)"; \
	go test ./... -coverprofile cover.out

# Build manager binary
# https://www.reddit.com/r/golang/comments/9ai79z/correct_usage_of_go_modules_vendor_still_connects/
.PHONY: manager
manager: generate fmt vet
# manager: generate fmt
	go build -mod vendor -o bin/manager main.go

# Run against the configured Kubernetes cluster in ~/.kube/config
.PHONY: run
run: generate fmt vet manifests
ifneq (,$(DEFAULT_STORAGE))
	DEFAULT_STORAGE=$(DEFAULT_STORAGE) \
	WATCH_NAMESPACE=$(WATCH_NAMESPACE) \
	RELATED_IMAGE_FREEIPA=$(RELATED_IMAGE_FREEIPA) \
	go run ./main.go $(CONTROLLER_ARGS)
else
	WATCH_NAMESPACE=$(WATCH_NAMESPACE) \
	RELATED_IMAGE_FREEIPA=$(RELATED_IMAGE_FREEIPA) \
	go run ./main.go $(CONTROLLER_ARGS)
endif

KUSTOMIZE = $(shell pwd)/bin/kustomize
.PHONY: kustomize
kustomize: ## Download kustomize locally if necessary.
	$(call go-get-tool,$(KUSTOMIZE),sigs.k8s.io/kustomize/kustomize/v3@v3.8.7)

# Install CRDs into a cluster
.PHONY: install-crds
install-crds: kustomize manifests
	$(KUSTOMIZE) build config/crd | kubectl apply -f -

# Uninstall CRDs from a cluster
.PHONY: uninstall-crds
uninstall-crds: kustomize manifests
	$(KUSTOMIZE) build config/crd | kubectl delete -f -

# Redeploy cluster updated
.PHONY: redeploy-cluster
redeploy-cluster: undeploy-cluster container-build container-push deploy-cluster

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
.PHONY: deploy-cluster
deploy-cluster: kustomize manifests
	cd config/manager && $(KUSTOMIZE) edit set image controller=$(IMG)
	cd config/default && $(KUSTOMIZE) edit set namespace $(WATCH_NAMESPACE)
	oc project $(WATCH_NAMESPACE) 2>/dev/null || oc new-project $(WATCH_NAMESPACE)
	$(KUSTOMIZE) build config/$(CONFIG) | kubectl create -f -

# Undeploy controller in the configured Kubernetes cluster in ~/.kube/config
.PHONY: deploy-cluster
undeploy-cluster: kustomize manifests
	cd config/manager && $(KUSTOMIZE) edit set image controller=$(IMG)
	-$(KUSTOMIZE) build config/$(CONFIG) | kubectl delete -f -

# Generate manifests e.g. CRD, RBAC etc.
.PHONY: manifests
manifests: controller-gen
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases

# Launch lint
.PHONY: lint
lint:
	./devel/lint.sh *.go $(shell find controllers -name '*.go') $(shell find api -name '*.go')

# Run go fmt against code
.PHONY: fmt
fmt:
	go fmt ./...

# Run go vet against code
.PHONY: vet
vet:
	go vet ./...

# Generate code
.PHONY: generate
generate: controller-gen
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."

# Check image size. It needs to run firstly container-build
.PHONY: dive
dive: $(CONTAINER_IMAGE_FILE)
	./devel/dive.sh

.PHONY: check-container-runtime
ifeq (,$(DOCKER))
check-container-runtime: FORCE
	@echo "ERROR: No docker nor podman were found"; exit 1
else
check-container-runtime:
endif

# Build the docker image
.PHONY: container-build
container-build: check-container-runtime
	$(DOCKER) build . -t $(IMG)

.PHONY: container-build-root
container-build-root: container-build container-save
	cat $(CONTAINER_IMAGE_FILE) | sudo $(DOCKER) load $(IMG)

.PHONY: container-delete-root
container-delete-root:
	sudo -E --preserve-env=HOME,PATH,GOPATH $(DOCKER) image rm $(IMG)

.PHONY: container-save
container-save: check-container-runtime $(CONTAINER_IMAGE_FILE)
$(CONTAINER_IMAGE_FILE): FORCE
	$(DOCKER) save $(IMG) > $(CONTAINER_IMAGE_FILE)

.PHONY: container-save-gz
container-save-gz: check-container-runtime $(CONTAINER_IMAGE_FILE).gz
$(CONTAINER_IMAGE_FILE).gz: FORCE
	$(DOCKER) save $(IMG) | gzip --best --force --stdout - > $(CONTAINER_IMAGE_FILE).gz

.PHONY: container-load
container-load: check-container-runtime
	@if [ ! -e "$(CONTAINER_IMAGE_FILE)" ]; then echo "No image file found. Run 'make container-build container-save' to generate a fresh image file before load it"; exit 2; fi
	cat $(CONTAINER_IMAGE_FILE) | $(DOCKER) load $(IMG)

.PHONY: container-load-gz
container-load-gz: check-container-runtime
	@if [ ! -e "$(CONTAINER_IMAGE_FILE).gz" ]; then echo "No image file found. Run 'make container-build container-save-gz' to generate a fresh image file before load it"; exit 2; fi
	gunzip $(CONTAINER_IMAGE_FILE).gz -c | $(DOCKER) load $(IMG)

# Push the docker image
.PHONY: container-push
container-push: check-container-runtime
	$(DOCKER) push $(IMG)

CONTROLLER_GEN = $(shell pwd)/bin/controller-gen
.PHONY: controller-gen
controller-gen: ## Download controller-gen locally if necessary.
	$(call go-get-tool,$(CONTROLLER_GEN),sigs.k8s.io/controller-tools/cmd/controller-gen@v0.6.2)

# https://itnext.io/testing-kubernetes-operators-with-ginkgo-gomega-and-the-operator-runtime-6ad4c2492379
.PHONY: deploy-kind
deploy-kind: kind-create kind-load-img deploy-cluster

.PHONY: kind
kind:
ifeq (, $(shell which kind))
	@(cd && GO111MODULE="on" go get sigs.k8s.io/kind@v0.10.0)
KIND=$(GOBIN)/kind
else
KIND=$(shell which kind)
endif

# make deploy-kind; make kind-tests
.PHONY: kind-create
ifeq (podman,$(DOCKER))
kind-create: kind
	@if (sudo -E --preserve-env=HOME,PATH,GOPATH $(KIND) get clusters 2>/dev/null | grep -q ^$(KIND_CLUSTER_NAME)\$$); \
	then \
	  echo "Cluster '$(KIND_CLUSTER_NAME)' already exists"; \
	else \
	  sudo -E --preserve-env=HOME,PATH,GOPATH $(KIND) create cluster --name $(KIND_CLUSTER_NAME) --image=kindest/node:$(K8S_NODE_IMAGE); \
	fi
else ifeq (docker,$(DOCKER))
kind-create:
	@if ($(KIND) get clusters 2>/dev/null | grep -q ^$(KIND_CLUSTER_NAME)\$$); \
	then \
	  echo "Cluster '$(KIND_CLUSTER_NAME)' already exists"; \
	else \
	  $(KIND) create cluster --name $(KIND_CLUSTER_NAME) --image=kindest/node:$(K8S_NODE_IMAGE); \
	fi
else
kind-create:
	@echo container enginer not supported; exit 1
endif

.PHONY: kind-delete
ifeq (podman,$(DOCKER))
kind-delete:
	@if (sudo -E --preserve-env=HOME,PATH,GOPATH $(KIND) get clusters 2>/dev/null | grep -q ^$(KIND_CLUSTER_NAME)\$$); \
	then \
	  sudo -E --preserve-env=HOME,PATH,GOPATH $(KIND) delete cluster --name $(KIND_CLUSTER_NAME); \
	else \
	  echo "Cluster '$(KIND_CLUSTER_NAME)' does not exist"; \
	fi
else ifeq (docker,$(DOCKER))
kind-delete:
	@if ($(KIND) get clusters 2>/dev/null | grep -q ^$(KIND_CLUSTER_NAME)\$$); \
	then \
	  $(KIND) delete cluster --name $(KIND_CLUSTER_NAME); \
	else \
	  echo "Cluster '$(KIND_CLUSTER_NAME)' does not exist"; \
	fi
else
kind-delete:
	@echo container enginer not supported; exit 1
endif

.PHONY: kind-load-img
ifeq (podman,$(DOCKER))
kind-load-img: container-build-root
	@echo "Loading image into kind"
	sudo -E --preserve-env=HOME,PATH,GOPATH -- $(KIND) load docker-image $(IMG) --name $(KIND_CLUSTER_NAME) --loglevel "trace"
else ifeq (docker,$(DOCKER))
kind-load-img: container-build
	@echo "Loading image into kind"
	$(KIND) load docker-image $(IMG) --name $(KIND_CLUSTER_NAME) --loglevel "trace"
else
kind-load-img: container-build-root
	@echo container enginer not supported; exit 1
endif

# Run integration tests in KIND
.PHONY: kind-tests
kind-tests:
	ginkgo --skip="LONG TEST:" --nodes 6 --race --randomizeAllSpecs --cover --trace --progress --coverprofile controllers.coverprofile ./controllers
	-kubectl delete IDM --all -n idm-system

.PHONY: kind-long-tests
kind-long-tests:
	ginkgo --focus="LONG TEST:" -nodes 6 --randomizeAllSpecs --trace --progress ./controllers
	-kubectl delete IDM --all -n idm-system

.PHONY: check-password-is-provided
check-password-is-provided:
ifeq (,$(IPA_ADMIN_PASSWORD))
	@echo "IPA_ADMIN_PASSWORD must be provided; IPA_ADMIN_PASSWORD=MySecretPassword make ..."; exit 1
endif
ifeq (,$(IPA_DM_PASSWORD))
	@echo "IPA_DM_PASSWORD must be provided; IPA_DM_PASSWORD=MySecretPassword make ..."; exit 1
endif


.PHONY: sample-build
sample-build:
	-kustomize build $(SAMPLES_PATH)/$(SAMPLE)/

.PHONY: sample-delete
sample-delete:
	-@kubectl delete secrets/idm-sample
	-kustomize build $(SAMPLES_PATH)/$(SAMPLE)/ | kubectl delete -f -

.PHONY: sample-create
sample-create: check-password-is-provided
	@-kubectl create secret generic idm-sample \
	          --from-literal=IPA_ADMIN_PASSWORD=$(IPA_ADMIN_PASSWORD) \
	          --from-literal=IPA_DM_PASSWORD=$(IPA_DM_PASSWORD)
	kustomize build $(SAMPLES_PATH)/$(SAMPLE)/ | kubectl create -f -

.PHONY: sample-recreate
sample-recreate: sample-delete sample-create

# go-get-tool will 'go get' any package $2 and install it to $1.
PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
define go-get-tool
@[ -f $(1) ] || { \
set -e ;\
TMP_DIR=$$(mktemp -d) ;\
cd $$TMP_DIR ;\
go mod init tmp ;\
echo "Downloading $(2)" ;\
GOBIN=$(PROJECT_DIR)/bin go get $(2) ;\
rm -rf $$TMP_DIR ;\
}
endef

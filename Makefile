IMG_NAME = freeipa-operator

ifneq (,$(shell command -v podman 2>/dev/null))
DOCKER=podman
else
ifneq (,$(shell command -v docker 2>/dev/null))
DOCKER=docker
else
DOCKER=
endif
endif

# https://docs.github.com/en/free-pro-team@latest/actions/reference/environment-variables#default-environment-variables
ifneq (,$(GITHUB_SHA))
COMMIT_SHA=$(GITHUB_SHA)
endif
# https://docs.travis-ci.com/user/environment-variables/#default-environment-variables
ifneq (,$(TRAVIS_COMMIT))
COMMIT_SHA=$(TRAVIS_COMMIT)
endif
# https://docs.gitlab.com/ee/ci/variables/predefined_variables.html
ifneq (,$(CI_COMMIT_SHA))
COMMIT_SHA=$(CI_COMMIT_SHA)
endif

ifeq (,$(COMMIT_SHA))
COMMIT_SHA=$(shell git rev-parse HEAD)
endif

CONTAINER_IMAGE_FILE ?= $(IMG_NAME).tar

IMG_TAG = dev-$(COMMIT_SHA)
IMG_BASE ?= quay.io/freeipa

# Image URL to use all building/pushing image targets
IMG ?= $(IMG_BASE)/$(IMG_NAME):$(IMG_TAG)
# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd:trivialVersions=true"

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

all: manager

# Empty rule to allow force other rules. The name of the rule should not
# match any file.
# https://www.gnu.org/software/make/manual/html_node/Force-Targets.html
FORCE:

# Run tests
test: generate fmt vet manifests
	go test ./... -coverprofile cover.out

# Build manager binary
manager: generate fmt vet
	go build -o bin/manager main.go

# Run against the configured Kubernetes cluster in ~/.kube/config
run: generate fmt vet manifests
	go run ./main.go

# Install CRDs into a cluster
install: manifests kustomize
	$(KUSTOMIZE) build config/crd | kubectl apply -f -

# Uninstall CRDs from a cluster
uninstall: manifests kustomize
	$(KUSTOMIZE) build config/crd | kubectl delete -f -

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
deploy: manifests kustomize
	cd config/manager && $(KUSTOMIZE) edit set image controller=${IMG}
	$(KUSTOMIZE) build config/default | kubectl apply -f -

# Generate manifests e.g. CRD, RBAC etc.
manifests: controller-gen
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases

# Launch lint
.PHONY: lint
lint:
	./devel/lint.sh

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

# Generate code
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
	$(DOCKER) build . -t ${IMG}

.PHONY: container-save
container-save: check-container-runtime $(CONTAINER_IMAGE_FILE)
$(CONTAINER_IMAGE_FILE): FORCE
	$(DOCKER) save ${IMG} > $(CONTAINER_IMAGE_FILE)

.PHONY: container-save-gz
container-save-gz: check-container-runtime $(CONTAINER_IMAGE_FILE).gz
$(CONTAINER_IMAGE_FILE).gz: FORCE
	$(DOCKER) save ${IMG} | gzip --best --force --stdout - > $(CONTAINER_IMAGE_FILE).gz

.PHONY: container-load
container-load: check-container-runtime
	@if [ ! -e "$(CONTAINER_IMAGE_FILE)" ]; then echo "No image file found. Run 'make container-build container-save' to generate a fresh image file before load it"; exit 2; fi
	cat $(CONTAINER_IMAGE_FILE) | $(DOCKER) load $(IMG)

.PHONY: container-load-gz
container-load-gz: check-container-runtime
	@if [ ! -e "$(CONTAINER_IMAGE_FILE).gz" ]; then echo "No image file found. Run 'make container-build container-save-gz' to generate a fresh image file before load it"; exit 2; fi
	gunzip $(CONTAINER_IMAGE_FILE).gz -c | $(DOCKER) load $(IMG)

# Push the docker image
container-push: check-container-runtime
	$(DOCKER) push ${IMG}

# find or download controller-gen
# download controller-gen if necessary
controller-gen:
ifeq (, $(shell which controller-gen))
	@{ \
	set -e ;\
	CONTROLLER_GEN_TMP_DIR=$$(mktemp -d) ;\
	cd $$CONTROLLER_GEN_TMP_DIR ;\
	go mod init tmp ;\
	go get sigs.k8s.io/controller-tools/cmd/controller-gen@v0.3.0 ;\
	rm -rf $$CONTROLLER_GEN_TMP_DIR ;\
	}
CONTROLLER_GEN=$(GOBIN)/controller-gen
else
CONTROLLER_GEN=$(shell which controller-gen)
endif

kustomize:
ifeq (, $(shell which kustomize))
	@{ \
	set -e ;\
	KUSTOMIZE_GEN_TMP_DIR=$$(mktemp -d) ;\
	cd $$KUSTOMIZE_GEN_TMP_DIR ;\
	go mod init tmp ;\
	go get sigs.k8s.io/kustomize/kustomize/v3@v3.5.4 ;\
	rm -rf $$KUSTOMIZE_GEN_TMP_DIR ;\
	}
KUSTOMIZE=$(GOBIN)/kustomize
else
KUSTOMIZE=$(shell which kustomize)
endif

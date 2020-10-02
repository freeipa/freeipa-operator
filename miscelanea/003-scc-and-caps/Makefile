# This makefile make life easier for playing with the different
# proof of concepts.

# Set the container runtime interface
ifneq (,$(shell bash -c "command -v podman 2>/dev/null"))
DOCKER?=podman
else
ifneq (,$(shell bash -c "command -v docker 2>/dev/null"))
DOCKER?=docker
else
ifeq (,$(DOCKER))
$(error DOCKER is not set)
endif
endif
endif

NAMESPACE=$(shell oc project -q)
ifeq (,$(NAMESPACE))
$(error Not project selected or not logged in an cluster)
endif

AWAIT_TIMEOUT?=10
APP_list=poc-003
POC_003_containers=init-container-uuid init-uid-gid-base init-volume main

APP=poc-003
CONTAINERS="$($(shell echo $(APP) | tr '[:lower:]' '[:upper:]' | tr '-' '_' )_containers)"
# Change DOCKER_IMAGE on your pipeline settings to point to your upstream
DOCKER_IMAGE?=quay.io/freeipa/freeipa-openshift-container:$(APP)


.PHONY: FORCE
FORCE:

.PHONY: help
help: FORCE
	@echo "Available commands:"
	@echo "    container-build"
	@echo "    container-push"
	@echo "    container-remove"
	@echo "    app-validate"
	@echo "    app-deploy"
	@echo "    app-delete"
	@echo "    get-info"

# Check logged in OpenShift cluster
.PHONY: check-logged-in-openshift
ifeq (,$(shell oc whoami 2>/dev/null))
check-logged-in-openshift: FORCE
	@echo "ERROR: You must be logged in OpenShift cluster. Try 'oc login https://mycluster' matching your cluster API endpoint"
	@exit 1
else
check-logged-in-openshift:
endif

# Check APP is not empty
.PHONY: check-app-not-empty
ifeq (,$(APP))
check-app-not-empty: FORCE
	@echo "'APP' must be specified: $(APP_list); Try 'APP=my-app-id make $0'"
	@exit 1
else
check-app-not-empty:
endif

# Check DOCKER_IMAGE is not empty
.PHONY: check-docker-image-not-empty
ifeq (,$(DOCKER_IMAGE))
check-docker-image-not-empty: FORCE
	@echo "'DOCKER_IMAGE' must be defined. Eg: 'export DOCKER_IMAGE=quay.io/myusername/freeipa-server:dev-test-003'"
	@exit 1
else
check-docker-image-not-empty:
endif

# Build the container image
.PHONY: container-build
container-build: check-docker-image-not-empty Dockerfile
	$(DOCKER) build -t $(DOCKER_IMAGE) -f Dockerfile .

# Push the container image to the container registry
.PHONY: container-push
container-push: check-docker-image-not-empty FORCE
	$(DOCKER) push $(DOCKER_IMAGE)

# Remove container image from the local storage
.PHONY: container-remove
container-remove: check-docker-image-not-empty FORCE
	$(DOCKER) image rm $(DOCKER_IMAGE)

# Validate kubernetes object for the app
.PHONY: app-validate
app-validate: check-logged-in-openshift check-app-not-empty FORCE
	oc create -f $(APP).yaml --dry-run --validate
	oc create -f $(APP)-admin.yaml --dry-run --validate

# Deploy the application
.PHONY: app-deploy
app-deploy: check-logged-in-openshift check-app-not-empty $(APP).yaml app-validate FORCE
	oc create -f $(APP)-admin.yaml
	oc adm policy add-scc-to-user $(APP) system:serviceaccount:$(NAMESPACE):$(APP)
	oc create -f $(APP).yaml

# Delete the application from the cluster
.PHONY: app-delete
app-delete: check-logged-in-openshift check-app-not-empty FORCE
	oc delete all,cm,sa,role,rolebinding,psp,scc,clusterrole,clusterrolebinding -l app=$(APP)

.PHONY: get-info
get-info: check-logged-in-openshift check-app-not-empty
	@echo ">>> oc describe pod/$(APP)"
	oc describe pod/$(APP)
	@for container in $(shell echo -n $(CONTAINERS)); do make get-info-container APP=$(APP) container=$${container}; done

.PHONY: check-container-not-empty
ifeq (,$(container))
check-container-not-empty: FORCE
	@[ "$(container)" != "" ] || echo "'container' must be specified: $(CONTAINERS); Try 'make APP=$(APP) container=my-container-id get-info-container'"
	@[ "$(container)" != "" ] || exit 1
else
check-container-not-empty:
endif

.PHONY: check-logged-in-openshift get-info-container
get-info-container: check-app-not-empty check-container-not-empty
	@echo ">>> container: $(container)"
	@if oc wait --for=condition=ready --timeout=$(AWAIT_TIMEOUT)s pod/$(APP); then oc logs $(APP) -c $(container); else oc logs -f $(APP) -c $(container) --insecure-skip-tls-verify-backend=true; fi


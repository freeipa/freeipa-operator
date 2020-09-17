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

APP_list=poc-01 poc-03 poc-04 poc-05-a poc-05-b poc-05-c poc-05-d freeipa
POC_01_containers=init-uid-gid-base main
POC_02_containers=init-uid-gid-base init-uid-gid main
POC_03_containers=init-my-pod-id main
POC_04_containers=init-container-uuid main
POC_05_A_containers=init-container-uuid init-uid-gid-base init-volume main
POC_05_B_containers=init-container-uuid init-uid-gid-base init-volume main
POC_05_C_containers=prepare main
POC_05_D_containers=prepare main
POC_05_E_containers=init-ping-pod main main-curl-pod

#APP?=poc-01
CONTAINERS="$($(shell echo $(APP) | tr '[:lower:]' '[:upper:]' | tr '-' '_' )_containers)"
#DOCKER_IMAGE?=quay.io/avisied0/freeipa-openshift-container:dev-test


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
	@echo "'DOCKER_IMAGE' must be defined. Eg: 'export DOCKER_IMAGE=quay.io/myusername/freeipa-server:dev-test'"
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
app-validate: check-app-not-empty FORCE
	oc create -f $(APP).yaml --dry-run --validate

# Deploy the application
.PHONY: app-deploy
app-deploy: check-app-not-empty $(APP).yaml app-validate FORCE
	oc create -f $(APP).yaml

# poc-05-a.yaml: poc-05-a.yaml.envsubst
# 	@envsubst < "$<" > "$@"

# Delete the application from the cluster
.PHONY: app-delete
app-delete: check-app-not-empty FORCE
	oc delete all,cm,sa,role,rolebinding -l app=$(APP)

.PHONY: get-info
get-info: check-app-not-empty
	@echo ">>> oc describe pod/$(APP)"
	oc describe pod/$(APP)
	@for container in $(shell echo -n $(CONTAINERS)); do make get-info-container APP=$(APP) container=$${container}; done

.PHONY: get-info-container
get-info-container: check-app-not-empty
	@[ "$(container)" != "" ] || echo "'container' must be specified: $(CONTAINERS); Try 'make APP=$(APP) container=my-container-id get-info-container'"
	@[ "$(container)" != "" ] || exit 1
	@echo ">>> container: $(container)"
	@if oc wait --for=condition=ready --timeout=10s pod/$(APP); then oc logs $(APP) -c $(container); else oc logs -f $(APP) -c $(container); fi

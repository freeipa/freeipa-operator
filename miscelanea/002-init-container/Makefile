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

AWAIT_TIMEOUT?=10
APP_list=poc-01 poc-03 poc-04 poc-05-a poc-05-b poc-05-c poc-05-d poc-05-e poc-05-f freeipa
POC_01_containers=init-uid-gid-base main
POC_02_containers=init-uid-gid-base init-uid-gid main
POC_03_containers=init-my-pod-id main
POC_04_containers=init-container-uuid main
POC_05_A_containers=init-container-uuid init-uid-gid-base init-volume main
POC_05_B_containers=init-container-uuid init-uid-gid-base init-volume main
POC_05_C_containers=prepare main
POC_05_D_containers=prepare main
POC_05_E_containers=init-ping-pod main main-curl-pod
POC_05_F_containers=

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
app-validate: check-logged-in-openshift check-app-not-empty FORCE
	oc create -f $(APP).yaml --dry-run --validate

# Deploy the application
.PHONY: app-deploy
app-deploy: check-logged-in-openshift check-app-not-empty $(APP).yaml app-validate FORCE
ifeq (poc-05-a,$(APP))
	oc create -f $(APP)-admin.yaml
	oc adm policy add-scc-to-user poc-05-a system:serviceaccount:avisiedo-init-container:poc-05-a
endif
	oc create -f $(APP).yaml

# poc-05-a.yaml: poc-05-a.yaml.envsubst
# 	@envsubst < "$<" > "$@"

# Delete the application from the cluster
.PHONY: app-delete
app-delete: check-logged-in-openshift check-app-not-empty FORCE
	oc delete all,cm,sa,role,rolebinding,psp,scc,clusterrole,clusterrolebinding -l app=$(APP)

ifeq (poc-05-f,$(APP))
.PHONY: get-info
get-info: check-logged-in-openshift check-app-not-empty
	@echo ">>> oc describe pod/$(APP)-default"
	oc describe pod/$(APP)-default
	@make get-info-container APP=$(APP)-default container=main-default
	@make get-info-container APP=$(APP)-default container=main-default-privileged
	@echo ">>> oc describe pod/$(APP)-default-with-dns-config"
	oc describe pod/$(APP)-default-with-dns-config
	@make get-info-container APP=$(APP)-default-with-dns-config container=main-default-with-dns-config
	@echo ">>> oc describe pod/$(APP)-default-with-host-aliases"
	oc describe pod/$(APP)-default-with-host-aliases
	@make get-info-container APP=$(APP)-default-with-host-aliases container=main-default-with-host-aliases
	@echo ">>> oc describe pod/$(APP)-anyuid"
	oc describe pod/$(APP)-anyuid
	@make get-info-container APP=$(APP)-anyuid container=main-anyuid
	@make get-info-container APP=$(APP)-anyuid container=main-anyuid-privileged
	@echo ">>> oc describe pod/$(APP)-restricted"
	oc describe pod/$(APP)-restricted
	@make get-info-container APP=$(APP)-restricted container=main-restricted
	@make get-info-container APP=$(APP)-restricted container=main-restricted-privileged
	@echo ">>> oc describe pod/$(APP)-privileged"
	oc describe pod/$(APP)-privileged
	@make get-info-container APP=$(APP)-privileged container=main-privileged
	@make get-info-container APP=$(APP)-privileged container=main-privileged-privileged
	oc describe pod/$(APP)-scc
	@make get-info-container APP=$(APP)-scc container=main-scc
	@make get-info-container APP=$(APP)-scc container=main-scc-privileged
else
.PHONY: get-info
get-info: check-logged-in-openshift check-app-not-empty
	@echo ">>> oc describe pod/$(APP)"
	oc describe pod/$(APP)
	@for container in $(shell echo -n $(CONTAINERS)); do make get-info-container APP=$(APP) container=$${container}; done
endif

.PHONY: check-logged-in-openshift get-info-container
get-info-container: check-app-not-empty
	@[ "$(container)" != "" ] || echo "'container' must be specified: $(CONTAINERS); Try 'make APP=$(APP) container=my-container-id get-info-container'"
	@[ "$(container)" != "" ] || exit 1
	@echo ">>> container: $(container)"
	@if oc wait --for=condition=ready --timeout=$(AWAIT_TIMEOUT)s pod/$(APP); then oc logs $(APP) -c $(container); else oc logs -f $(APP) -c $(container); fi

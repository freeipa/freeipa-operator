# This makefile make life easier for playing with the different
# proof of concepts.

# Customizable cluster domain for deploying wherever we want
CLUSTER_DOMAIN?=permanent.idmocp.lab.eng.rdu2.redhat.com
REALM?=APPS.$(shell echo $(CLUSTER_DOMAIN) | tr  '[:lower:]' '[:upper:]')
TIMESTAMP?=$(shell date +%Y%m%d%H%M%S)

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
$(error Not project indicated)
endif

AWAIT_TIMEOUT?=10
FREEIPA_containers=init-container-uuid init-uid-gid-base main

APP?=$(NAMESPACE)
CONTAINERS="$(FREEIPA_containers)"
# Change DOCKER_IMAGE on your pipeline settings to point to your upstream
DOCKER_IMAGE?=quay.io/freeipa/freeipa-openshift-container:freeipa-server


default: help


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

# Check that cluster domain is not empty
.PHONY: check-cluster-domain-not-empty
ifeq (,$(CLUSTER_DOMAIN))
check-cluster-domain-not-empty: FORCE
	@echo "'CLUSTER_DOMAIN' must be specified; Try 'CLUSTER_DOMAIN=my.cluster.domain.com make $$0'"
	@exit 1
else
check-cluster-domain-not-empty:
endif

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
	@echo "'APP' must be specified; Try 'APP=my-app-id make $0'"
	@exit 1
else
check-app-not-empty:
endif

# Check DOCKER_IMAGE is not empty
.PHONY: check-docker-image-not-empty
ifeq (,$(DOCKER_IMAGE))
check-docker-image-not-empty: FORCE
	@echo "'DOCKER_IMAGE' must be defined. Eg: 'export DOCKER_IMAGE=quay.io/myusername/freeipa-server:latest'"
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
app-validate: check-logged-in-openshift check-app-not-empty freeipa.yaml freeipa-admin.yaml FORCE
	oc create -f freeipa.yaml --dry-run=client --validate=true \
	&& oc create -f freeipa-admin.yaml --dry-run=client --validate=true

# Deploy the application
.PHONY: app-deploy
app-deploy: check-cluster-domain-not-empty check-app-not-empty check-docker-image-not-empty check-logged-in-openshift freeipa.yaml freeipa-admin.yaml app-validate FORCE
	oc create user freeipa
	oc create -f freeipa-admin.yaml
	# oc adm policy add-scc-to-user freeipa -z system:serviceaccount:$(NAMESPACE):freeipa
	oc adm policy add-scc-to-user freeipa -z freeipa
	oc create -f freeipa.yaml --as freeipa
	oc get routes/freeipa

# Force manifest regeneration
freeipa-admin.yaml: freeipa-admin.yaml.envsubst FORCE
	DOLLAR='$$' CLUSTER_DOMAIN="$(CLUSTER_DOMAIN)" NAMESPACE=$(NAMESPACE) DOCKER_IMAGE=$(DOCKER_IMAGE) APP=$(APP) envsubst < "$<" > "$@"

# Force manifest regeneration
freeipa.yaml: freeipa.yaml.envsubst FORCE
	DOLLAR='$$' TIMESTAMP=$(TIMESTAMP) REALM="$(REALM)" CLUSTER_DOMAIN="$(CLUSTER_DOMAIN)" NAMESPACE=$(NAMESPACE) DOCKER_IMAGE=$(DOCKER_IMAGE) APP=$(APP) envsubst < "$<" > "$@"

# Delete the application from the cluster
.PHONY: app-delete
app-delete: check-logged-in-openshift check-app-not-empty FORCE
	oc delete all,cm,sa,role,rolebinding,psp,scc,clusterrole,clusterrolebinding -l app=freeipa
	oc delete user freeipa || true

.PHONY: get-info
get-info: check-logged-in-openshift check-app-not-empty
	@echo ">>> oc describe pod/freeipa"
	oc describe pod/freeipa
	@for container in $(shell echo -n $(CONTAINERS)); do make get-info-container container=$${container}; done

.PHONY: check-container-not-empty
ifeq (,$(container))
check-container-not-empty: FORCE
	@[ "$(container)" != "" ] || echo "'container' must be specified: $(CONTAINERS); Try 'make container=my-container-id get-info-container'"
	@[ "$(container)" != "" ] || exit 1
else
check-container-not-empty:
endif

.PHONY: check-logged-in-openshift get-info-container
get-info-container: check-app-not-empty check-container-not-empty
	@echo ">>> container: $(container)"
	@if oc wait --for=condition=ready --timeout=$(AWAIT_TIMEOUT)s pod/freeipa; then oc logs freeipa -c $(container); else oc logs -f freeipa -c $(container) --insecure-skip-tls-verify-backend=true; fi

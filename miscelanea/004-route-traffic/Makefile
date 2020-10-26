
ifneq (,$(shell command -v podman 2>/dev/null))
DOCKER?=podman
else
ifneq (,$(shell command -v docker 2>/dev/null))
DOCKER?=docker
else
DOCKER?=
endif
endif

NAMESPACE=$(shell oc project -q)
ifeq (,$(NAMESPACE))
$(error Not project selected or not logged in an cluster)
endif

AWAIT_TIMEOUT?=10
APP_list=poc-004-1 poc-004-2 poc-004-3 poc-004-4a poc-004-4b poc-004-4c poc-004-4d
POC_004_containers=http https ldap ldaps kerberos-tcp kerberos-udp kerberos-admin-tcp kerberos-admin-udp

# APP=poc-004-1
CONTAINERS="$($(shell echo $(APP) | tr '[:lower:]' '[:upper:]' | tr '-' '_' )_containers)"
# Change DOCKER_IMAGE on your pipeline settings to point to your upstream
DOCKER_IMAGE?=quay.io/freeipa/freeipa-operator:poc-004

# Environment variables used by the container
ECHO_ADDRESS?=0.0.0.0
ECHO_PORT?=8007
ECHO_MAX_CLIENTS?=1

.PHONY: help
help:
	@echo "build             Build the DOCKER_IMAGE container image."
	@echo "push              Push the DOCKER_IMAGE container image."
	@echo "run-udp           Run a UDP server locally."
	@echo "run-tcp           Run a TCP server locally."
	@echo "app-deploy        Deploy the proof of concept application."
	@echo "app-delete        Remove the application from the cluster."
	@echo "------------------------------------------------------------------------------"
	@echo "APP               $(APP_list)"
	@echo "ECHO_ADDRESS      The address which the server bind to (0.0.0.0 by default)"
	@echo "ECHO_PORT         The port which the server bind to (8007 by default)"
	@echo "ECHO_MAX_CLIENTS  For tcp server, the maximum clients accepted (1 by default)"

.PHONY: FORCE
FORCE:

# Check the container frontend was found
.PHONY: check-docker
ifeq (,$(DOCKER))
check-docker:
	@echo "No docker nor podman found in the system"; exit 1
else
check-docker:
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

.PHONY: check-container-not-empty
ifeq (,$(container))
check-container-not-empty: FORCE
	@[ "$(container)" != "" ] || echo "'container' must be specified: $(CONTAINERS); Try 'make APP=$(APP) container=my-container-id get-info-container'"
	@[ "$(container)" != "" ] || exit 1
else
check-container-not-empty:
endif

.PHONY: build
build: check-docker check-docker-image-not-empty
	$(DOCKER) build -t $(DOCKER_IMAGE) .

.PHONY: push
push: check-docker check-docker-image-not-empty
	$(DOCKER) push $(DOCKER_IMAGE)

.PHONY: run-udp
run-udp: check-docker
	$(DOCKER) run --rm -it -p$(ECHO_PORT):$(ECHO_PORT)/udp -e ECHO_PROTOCOL=udp -e ECHO_ADDRESS=$(ECHO_ADDRESS) -e ECHO_PORT=$(ECHO_PORT) $(DOCKER_IMAGE)

.PHONY: run-tcp
run-tcp: check-docker check-docker-image-not-empty
	$(DOCKER) run --rm -it -p$(ECHO_PORT):$(ECHO_PORT)/tcp -e ECHO_PROTOCOL=tcp -e ECHO_ADDRESS=$(ECHO_ADDRESS) -e ECHO_PORT=$(ECHO_PORT) -e ECHO_MAX_CLIENTS=$(ECHO_MAX_CLIENTS) $(DOCKER_IMAGE)

# Validate kubernetes object for the app
.PHONY: app-validate
app-validate: check-logged-in-openshift check-app-not-empty FORCE
	oc create -f $(APP).yaml --dry-run --validate
	# oc create -f $(APP)-admin.yaml --dry-run --validate
	oc create -f poc-004-admin.yaml --dry-run --validate

# Deploy the application
.PHONY: app-deploy
app-deploy: check-logged-in-openshift check-app-not-empty $(APP).yaml app-validate FORCE
	# oc create -f $(APP)-admin.yaml
	oc apply -f poc-004-admin.yaml || true
	# oc adm policy add-scc-to-user $(APP) system:serviceaccount:$(NAMESPACE):$(APP)
	oc create -f $(APP).yaml

# Delete the application from the cluster
.PHONY: app-delete
app-delete: check-logged-in-openshift check-app-not-empty FORCE
	oc delete all,cm,sa,role,rolebinding,psp,scc,clusterrole,clusterrolebinding -l app=$(APP)

# Read logs
.PHONY: get-info
get-info: check-logged-in-openshift check-app-not-empty
	@echo ">>> oc describe pod/$(APP)"
	oc describe pod/$(APP)
	@for container in $(shell echo -n $(CONTAINERS)); do make get-info-container APP=$(APP) container=$${container}; done

.PHONY: check-logged-in-openshift get-info-container
get-info-container: check-app-not-empty check-container-not-empty
	@echo ">>> container: $(container)"
	if oc wait --for=condition=ready --timeout=$(AWAIT_TIMEOUT)s pod/$(APP); then oc logs $(APP) -c $(container); else oc logs -f $(APP) -c $(container); fi

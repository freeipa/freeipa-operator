# Default goal
.DEFAULT_GOAL := help

# This allow to run locally the controller without webhooks by default
# which let you debug the code directly. A more complex configuration
# is needed to debug with webhooks.
ENABLE_WEBHOOKS ?= false
export ENABLE_WEBHOOKS

# Set the workload image to be used to deploy Freeipa when a new
# IDM custom resource is created. This let you test different
# images locally without deploy into the cluster.
RELATED_IMAGE_FREEIPA ?= quay.io/freeipa/freeipa-openshift-container:latest
export RELATED_IMAGE_FREEIPA

# The namespace to be watched by the controller; by default it is set to
# the current namespace
WATCH_NAMESPACE ?= $(shell oc project -q 2>/dev/null)
export WATCH_NAMESPACE

# Include sample rules
include mk/scorecard.mk
include mk/checks.mk
include mk/samples.mk
include mk/cert-manager.mk
include mk/miscelanea.mk
include mk/deprecated.mk

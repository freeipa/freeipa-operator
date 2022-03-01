SAMPLE ?= ./config/samples/ephemeral-storage

##@ Sample idm custom resources management

.PHONY: sample-build
sample-build: ## Print out the resulting IDM resource
	-$(KUSTOMIZE) build $(SAMPLE)

.PHONY: sample-delete
sample-delete:  ## Delete the IDM sample resource
	-@kubectl delete secrets/idm-sample
	-$(KUSTOMIZE) build $(SAMPLE) | kubectl delete --wait=true -f -

.PHONY: sample-create
sample-create: check-password-is-provided  ## Create the IDM sample resource
	@-kubectl create secret generic idm-sample \
	          --from-literal=IPA_ADMIN_PASSWORD=$(IPA_ADMIN_PASSWORD) \
	          --from-literal=IPA_DM_PASSWORD=$(IPA_DM_PASSWORD)
	$(KUSTOMIZE) build $(SAMPLE) | kubectl create -f -

.PHONY: sample-recreate
sample-recreate: sample-delete sample-create #

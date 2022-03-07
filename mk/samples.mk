SAMPLE ?= ./config/samples/ephemeral-storage

##@ Samples

.PHONY: sample-build
sample-build:
	-$(KUSTOMIZE) build $(SAMPLE)

.PHONY: sample-create
sample-create: check-password-is-provided  ## Create the IDM sample resource
	@kubectl get secret/idm-sample &>/dev/null \
	 || kubectl create secret generic idm-sample \
	          --from-literal=IPA_ADMIN_PASSWORD='$(IPA_ADMIN_PASSWORD)' \
	          --from-literal=IPA_DM_PASSWORD='$(IPA_DM_PASSWORD)'
	$(KUSTOMIZE) build $(SAMPLE) | kubectl create -f -

.PHONY: sample-delete
sample-delete:  ## Delete the IDM sample resource
	@kubectl get secret/idm-sample 2>/dev/null \
	 || kubectl delete secrets/idm-sample
	-$(KUSTOMIZE) build $(SAMPLE) | kubectl delete --wait=true -f -

.PHONY: sample-recreate
.NOTPARALLEL: sample-recreate
sample-recreate: sample-delete sample-create

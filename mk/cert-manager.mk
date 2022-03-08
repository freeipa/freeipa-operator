
##@ Cert Manager

.PHONY: cert-manager-install
cert-manager-install:  ## Install cert-manager operator
	kubectl create -f config/certmanager/subscription.yaml
	kubectl wait Subscription/cert-manager -n openshift-operators --for=condition=CatalogSourcesUnhealthy=False
	$(call kubectl-wait-for-value,Succeeded,csv/cert-manager.v1.6.1,.status.phase)

.PHONY: cert-manager-uninstall
cert-manager-uninstall:  ## Delete cert-manager operator
	kubectl delete -f config/certmanager/subscription.yaml --wait=true


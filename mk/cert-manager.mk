
CERT_MANAGER_CVS_NAME=$(shell kubectl get Subscription/cert-manager -n openshift-operators -o jsonpath='{.status.installedCSV}' 2>/dev/null)

##@ Cert Manager

.PHONY: cert-manager-install
cert-manager-install:  ## Install cert-manager operator
	kubectl create -f config/certmanager/subscription.yaml

.PHONY: cert-manager-uninstall
cert-manager-uninstall:  ## Delete cert-manager operator
ifneq (,$(CERT_MANAGER_CVS_NAME))
	kubectl delete ClusterServiceVersion/$(CERT_MANAGER_CVS_NAME) -n openshift-operators
endif
	kubectl delete -f config/certmanager/subscription.yaml

.PHONY: cert-manager-self-signed-issuer-create
cert-manager-self-signed-issuer-create:  ## Create the self-signed certificate issuer
	kubectl create -f config/certmanager/clusterissuer-selfsigned.yaml

.PHONY: cert-manager-self-signed-issuer-delete
cert-manager-self-signed-issuer-delete:  ## Delete the self-signed certificate issuer
	kubectl delete -f config/certmanager/clusterissuer-selfsigned.yaml


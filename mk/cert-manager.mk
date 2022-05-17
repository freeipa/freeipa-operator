
##@ Cert Manager

.PHONY: cert-manager-install
cert-manager-install:  ## Install cert-manager operator
	oc create -f config/certmanager/subscription.yaml
	oc wait Subscription/cert-manager -n openshift-operators --for=condition=CatalogSourcesUnhealthy=False

.PHONY: cert-manager-uninstall
cert-manager-uninstall:  ## Delete cert-manager operator
	oc delete -f config/certmanager/subscription.yaml --wait=true

.PHONY: cert-manager-self-signed-issuer-create
cert-manager-self-signed-issuer-create:  ## Create a cluster self signed issuer
	oc create -f config/certmanager/clusterissuer-selfsigned.yaml

.PHONY: cert-manager-self-signed-issuer-delete
cert-manager-self-signed-issuer-delete:  ## Delete the cluster self signed issuer
	oc delete -f config/certmanager/clusterissuer-selfsigned.yaml

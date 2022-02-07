
CERT_MANAGER_CVS_NAME=$(shell kubectl get Subscription/cert-manager -n openshift-operators -o jsonpath='{.status.installedCSV}')

.PHONY: cert-manager-install
cert-manager-install:
	kubectl create -f config/certmanager/subscription.yaml

.PHONY: cert-manager-uninstall
cert-manager-uninstall:
ifneq (,$(CERT_MANAGER_CVS_NAME))
	kubectl delete ClusterServiceVersion/$(CERT_MANAGER_CVS_NAME) -n openshift-operators
endif
	kubectl delete -f config/certmanager/subscription.yaml

.PHONY: cert-manager-self-signed-issuer-create
cert-manager-self-signed-issuer-create:
	kubectl create -f config/certmanager/clusterissuer-selfsigned.yaml

.PHONY: cert-manager-self-signed-issuer-delete
cert-manager-self-signed-issuer-delete:
	kubectl delete -f config/certmanager/clusterissuer-selfsigned.yaml


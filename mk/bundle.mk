.PHONY: bundle-install
bundle-install: ## Install the operator into the current cluster and namespace by using OLM
	$(OPERATOR_SDK) run bundle $(BUNDLE_IMG) --namespace "$(WATCH_NAMESPACE)"

# TODO Review this rule and find a better way to retrieve the information
.PHONY: bundle-delete
bundle-uninstall: SUBSCRIPTION_NAME=$(shell echo "freeipa-operator.v$(VERSION)-sub" | sed -e 's/\./-/g')
bundle-uninstall: CLUSTER_SERVICE_VERSION=$(shell echo "freeipa-operator.v$(VERSION)")
bundle-uninstall: ## Uninstall the Operator from the current cluster and namespace
	$(OPERATOR_SDK) cleanup freeipa-operator --namespace "$(WATCH_NAMESPACE)" --delete-all

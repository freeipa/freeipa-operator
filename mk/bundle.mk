.NOTPARALLEL: bundle-install
.PHONY: bundle-install
bundle-install: docker-build docker-push bundle-push catalog-push ## Install the operator into the current cluster and namespace by using OLM
	$(OPERATOR_SDK) run bundle $(BUNDLE_IMG) --namespace "$(WATCH_NAMESPACE)"

.PHONY: bundle-delete
bundle-uninstall: ## Uninstall the Operator from the current cluster and namespace
	$(OPERATOR_SDK) cleanup freeipa-operator --namespace "$(WATCH_NAMESPACE)"

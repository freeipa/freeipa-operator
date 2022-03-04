define obsolete-rule
$(if $1,$(if $2,$(error Rule '$1' is obsolete; Use '$2' instead),$(error Rule '$1' is obsolete)))
endef

define deprecated-rule
$(if $1,$(warning Rule '$1' is deprecated; Use '$2' instead),)
endef

## >>NOTE<< Deprecated rules, do not use them

.NOTPARALLEL: redeploy-cluster
.PHONY: redeploy-cluster
redeploy-cluster: undeploy deploy
	$(call deprecated-rule,$@,$<)

.PHONY: deploy-cluster
deploy-cluster: deploy
	$(call deprecated-rule,$@,$<)

.PHONY: undeploy-cluster
undeploy-cluster: deploy
	$(call deprecated-rule,$@,$<)

.PHONY: container-build
container-build: docker-build
	$(call deprecated-rule,$@,$<)

.PHONY: container-push
container-push: docker-push  ## Push the controller image to the image registry
	$(call deprecated-rule,$@,$<)

## >>NOTE<< Obsolete rules, they will fails

.PHONY: container-build-root
container-build-root:
	$(call obsolete-rule,$@,docker-build)

.PHONY: container-delete-root
container-delete-root:
	$(call obsolete-rule,$@)

.PHONY: check-container-runtime
check-container-runtime:
	$(call obsolete-rule,$@)

.PHONY: kind
kind:
	$(call obsolete-rule,$@)

.PHONY: kind-create
kind-create:
	$(call obsolete-rule,$@)

.PHONY: kind-delete
kind-delete:
	$(call obsolete-rule,$@)

.PHONY: kind-load-img
kind-load-img:
	$(call obsolete-rule,$@)

.PHONY: kind-tests
kind-tests:
	$(call obsolete-rule,$@)

.PHONY: kind-long-tests
kind-long-tests:
	$(call obsolete-rule,$@)

.PHONY: deploy-kind
deploy-kind:
	$(call obsolete-rule,$@)

.PHONY: undeploy-kind
undeploy-kind:
	$(call obsolete-rule,$@)

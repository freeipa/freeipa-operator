##@ Miscelanea

# https://www.cmcrossroads.com/article/dumping-every-makefile-variable
.PHONY: printvars
printvars: ## Print variable name and values
	@$(foreach V, $(sort $(.VARIABLES)),$(if $(filter-out environment% default automatic,$(origin $V)),$(info $V=$($V))))

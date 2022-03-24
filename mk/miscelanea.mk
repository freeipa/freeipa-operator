##@ Miscelanea

# https://www.cmcrossroads.com/article/dumping-every-makefile-variable
.PHONY: printvars
printvars: ## Print variable name and values
	@$(foreach V, $(sort $(.VARIABLES)),$(if $(filter-out environment% default automatic,$(origin $V)),$(info $V=$(value $V))))

$(PROJECT_DIR)/bin:
	[ -e bin ] || mkdir bin

OPERATOR_SDK := $(PWD)/bin/operator-sdk
.PHONY: operator-sdk
operator-sdk: $(PROJECT_DIR)/bin ## Donwload operator-sdk
ifeq (,$(wildcard $(OPERATOR_SDK)))
	set -e ; \
	export ARCH=$$(case $$(uname -m) in x86_64) echo -n amd64 ;; aarch64) echo -n arm64 ;; *) echo -n $$(uname -m) ;; esac) ; \
	export OS=$$(uname | awk '{print tolower($$0)}') ; \
	export OPERATOR_SDK_DL_URL=https://github.com/operator-framework/operator-sdk/releases/download/v1.18.0 ; \
	curl -sSLo "$(OPERATOR_SDK)" "$$OPERATOR_SDK_DL_URL/operator-sdk_$${OS}_$${ARCH}" && \
	[ -x "$(OPERATOR_SDK)" ] || chmod a+x "$(OPERATOR_SDK)"
endif

.PHONY: lint
lint:  ## Run linters
	./devel/lint.sh *.go $(shell find controllers -name '*.go') $(shell find api -name '*.go')

.PHONY: tidy
tidy:  ## Update golang dependencies
	go mod tidy

.PHONY: vendor
vendor:  ## Update vendor directory
	go mod vendor

.PHONY: .venv
.venv:
	python3 -m venv .venv
	source .venv/bin/activate; pip install --upgrade pip
	source .venv/bin/activate; pip install -r requirements-dev.txt

ifeq (/.cache,$(shell go env GOCACHE))
GOCACHE=/tmp/.cache
export GOCACHE
endif


##@ Scorecard

IMG_SCORECARD=$(IMAGE_TAG_BASE)-scorecard:v$(VERSION)

# https://github.com/operator-framework/operator-sdk/blob/09c3aa14625965af9f22f513cd5c891471dbded2/Makefile#L78
.PHONY: scorecard-build
scorecard-build:  ## Build the image for the custom scorecard tests
	docker build -t $(IMG_SCORECARD) -f images/scorecard-tests/Dockerfile .

.PHONY: scorecard-push
scorecard-push:  ## Push the image for the custom scorecard tests
	docker image push $(IMG_SCORECARD)

.PHONY: scorecard-bundle
scorecard-bundle: operator-sdk ## Execute scorecard from bundle/ directory; You can specify selector eg 'make scorecard-bundle scorecard_selector=test=olm-status-descriptors-test'
ifneq (,$(scorecard_selector))
	$(OPERATOR_SDK) scorecard bundle --selector='$(scorecard_selector)'
else
	$(OPERATOR_SDK) scorecard bundle
endif

.PHONY: scorecard-image
scorecard-image: operator-sdk ## Execute scorecard from the generated image (BUNDLE_IMG); You can specify selector eg 'make scorecard-image scorecard_selector=test=olm-status-descriptors-test'
ifneq (,$(scorecard_selector))
	$(OPERATOR_SDK) scorecard $(BUNDLE_IMG) --selector $(scorecard_selector)
else
	$(OPERATOR_SDK) scorecard $(BUNDLE_IMG)
endif

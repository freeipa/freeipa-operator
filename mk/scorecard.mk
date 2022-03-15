
##@ Scorecard

scorecard-bundle: operator-sdk ## Execute scorecard from bundle/ directory; You can specify selector eg 'make scorecard-bundle scorecard_selector=test=olm-status-descriptors-test'
ifneq (,$(scorecard_selector))
	$(OPERATOR_SDK) scorecard bundle --selector $(scorecard_selector)
else
	$(OPERATOR_SDK) scorecard bundle
endif

scorecard-image: operator-sdk ## Execute scorecard from the generated image (BUNDLE_IMG); You can specify selector eg 'make scorecard-image scorecard_selector=test=olm-status-descriptors-test'
ifneq (,$(scorecard_selector))
	$(OPERATOR_SDK) scorecard $(BUNDLE_IMG) --selector $(scorecard_selector)
else
	$(OPERATOR_SDK) scorecard $(BUNDLE_IMG)
endif

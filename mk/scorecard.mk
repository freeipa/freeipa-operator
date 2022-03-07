
##@ Scorecard

scorecard-bundle: operator-sdk ## Execute scorecard from bundle/ directory
	$(OPERATOR_SDK) scorecard bundle

scorecard-image: operator-sdk ## Execute scorecard from the generated image (BUNDLE_IMG)
	$(OPERATOR_SDK) scorecard $(BUNDLE_IMG)

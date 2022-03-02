
##@ Scorecard

scorecard-bundle: ## Execute scorecard 
	operator-sdk scorecard bundle

scorecard-image: ## Execure scorecard from the image
	operator-sdk scorecard $(BUNDLE_IMG)

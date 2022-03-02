

PHONY: check-password-is-provided
check-password-is-provided:
ifeq (,$(IPA_ADMIN_PASSWORD))
	@echo IPA_ADMIN_PASSWORD was not provided; exit 1
endif
ifeq (,$(IPA_DM_PASSWORD))
	@echo IPA_DM_PASSWORD was not provided; exit 1
endif

DIVE_VERSION=0.10.0
DIVE_BIN=$(shell pwd)/bin/dive
.PHONY: dive
dive: $(DIVE_BIN)
	TMP_FILE="$(shell mktemp dive.XXXXXX.tar)"; \
	  docker image save $(IMG) > "$${TMP_FILE}" \
	  && $(DIVE_BIN) "docker-archive://$${TMP_FILE}" \
	  && rm "$${TMP_FILE}"

$(DIVE_BIN):
	[ -e "bin" ] || mkdir "bin"
	curl --silent -L "https://github.com/wagoodman/dive/releases/download/v${DIVE_VERSION}/dive_${DIVE_VERSION}_$(shell go env GOOS)_$(shell go env GOARCH).tar.gz" | tar xz -C "./bin/" dive

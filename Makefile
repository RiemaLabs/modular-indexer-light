VERSION := $(shell git describe --tags 2>/dev/null)
GIT_HASH := $(shell git rev-parse --short HEAD)

GOOS :=
ENV := GOOS=${GOOS}

LDFLAGS := \
	-X main.version=${VERSION} \
	-X main.gitHash=${GIT_HASH}
FLAGS := -ldflags='${LDFLAGS}'

modular-indexer-light:
	env ${ENV} go build ${FLAGS} ./cmd/$@

.PHONY: clean
clean:
	rm -rf ./modular-indexer-light

VERSION := $(shell git describe --tags 2>/dev/null)
GIT_HASH := $(shell git rev-parse --short HEAD)

GOOS :=
ENV := GOOS=${GOOS}

LDFLAGS := \
	-X ./cmd/modular-indexer-light/main.version=${VERSION} \
	-X ./cmd/modular-indexer-light/main.gitHash=${GIT_HASH}
FLAGS := -ldflags='${LDFLAGS}'

modular-indexer-light:
	env ${ENV} go build ${FLAGS} ./cmd/$@

VERSION := $(shell git describe --tags 2>/dev/null)
GIT_HASH := $(shell git rev-parse --short HEAD)

GOOS :=
GOARCH :=
ENV := GOOS=${GOOS} GOARCH=${GOARCH}

LDFLAGS := \
	-X main.version=${VERSION} \
	-X main.gitHash=${GIT_HASH}
FLAGS := -ldflags='${LDFLAGS}'

EXEC := modular-indexer-light

${EXEC}:
	env ${ENV} go build ${FLAGS} ./cmd/$@

packages/${EXEC}/wasm_exec.js:
	cp $(shell go env GOROOT)/misc/wasm/wasm_exec.js $@

packages/${EXEC}/${EXEC}.wasm: packages/${EXEC}/wasm_exec.js
	env ${ENV} go build ${FLAGS} -o $@ ./internal/wasm/${EXEC}

.PHONY: clean
clean:
	rm -rf ./${EXEC} packages/${EXEC}/*.wasm

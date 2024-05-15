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
SDK := ${EXEC}-sdk

${EXEC}:
	env ${ENV} go build ${FLAGS} ./cmd/$@

packages/${SDK}/wasm_exec.js:
	cp $(shell go env GOROOT)/misc/wasm/wasm_exec.js $@

packages/${SDK}/${EXEC}.wasm: packages/${SDK}/wasm_exec.js
	env ${ENV} go build ${FLAGS} -o $@ ./internal/wasm/${EXEC}

.PHONY: fmt
fmt:
	goimports -w .

.PHONY: ci
ci:
	go run github.com/RiemaLabs/nubit-ci/cmd/nubitci-lint@latest

.PHONY: clean
clean:
	rm -rf *.exe ./${EXEC} packages/${SDK}/*.wasm

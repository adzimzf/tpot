
VERSION ?= $(shell git describe --tags 2>/dev/null || git rev-parse --short HEAD)

GO_LDFLAGS := -X main.Version=$(VERSION) $(GO_LDFLAGS)

build:
	@go build -trimpath -ldflags "$(GO_LDFLAGS)" -o "$@" main.go && \
	mv build tpot
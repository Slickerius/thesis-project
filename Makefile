GOFILES!=find . -name '*.go'
VERSION!=git describe --tags --dirty 2>/dev/null | grep . || echo "devel"
COMMIT!=git rev-parse --short HEAD 2>/dev/null
LDFLAGS=-X main.Commit=$(COMMIT) -X main.Version=$(VERSION)
GO=go
TAGS=

communiqué: go.mod go.sum $(GOFILES)
	$(GO) build \
		-trimpath \
		-tags "$(TAGS)" \
		-o $@ \
		-ldflags "$(LDFLAGS)"

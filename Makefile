GO := go
GOTEST = $(GO) test
GOVET = $(GO) vet
INSTALL := install
LINT_TOOL := golangci-lint

BINARY_NAME := taskgo
PREFIX ?= /usr/local

.PHONY = all
all: build

.PHONY: build
build: go.mod go.sum
	$(GO) build ./cmd/taskgo

.PHONY: run
run: taskgo
	./taskgo

.PHONY: install
install: $(BINARY_NAME)
	$(INSTALL) -d $(PREFIX)/bin/
	$(INSTALL) -m 755 taskgo $(PREFIX)/bin/taskgo

.PHONY = uninstall
uninstall:
	$(RM) -f $(PREFIX)/bin/taskgo

# Development Stuff
.PHONY: vendor
vendor:
	go mod vendor

.PHONY = test
test:
	$(GOTEST) -race ./...

.PHONY = bench
bench:
	$(GOTEST) -bench=. -benchtime=5s -run=Bench -v ./...

.PHONY = fmt
fmt:
	$(GO) fmt ./...

.PHONY = vet
vet:
	go vet ./...

.PHONY = lint
lint:
	$(LINT_TOOL) run

.PHONY = clean
clean:
	rm -f taskgo
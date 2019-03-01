LINT_FLAGS := run --deadline=120s
LINTER := ./bin/golangci-lint
TESTFLAGS := -v -cover

GO111MODULE := on
all: $(LINTER) deps test lint build

$(LINTER):
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s v1.15.0

.PHONY: lint
lint: $(LINTER)
	$(LINTER) $(LINT_FLAGS) ./...

.PHONY: deps
deps:
	go get .

.PHONY: build
build:
	go build .

.PHONY: test
test:
	go test $(TESTFLAGS) ./...

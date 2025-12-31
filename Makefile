VERSION ?= $(shell git describe --tags --always)
GOBIN ?= $(shell go env GOPATH)/bin

bin/bnf: build

run:
	go run \
		-ldflags="-X main.BuildVersion=$(VERSION)" \
		. \
		./examples/numbers.bnf ./examples/numbers.test

build:
	go build \
		-ldflags="-X main.BuildVersion=$(VERSION)" \
		-o bin/bnf .

test:
	go test -v ./...

install: bin/bnf
	@mkdir -p $(GOBIN)
	@cp bin/bnf $(GOBIN)/bnf
	@echo "Installed bnf to $(GOBIN)/bnf"

$(GOBIN)/goimports:
	@go install golang.org/x/tools/cmd/goimports@v0.35.0

$(GOBIN)/gocyclo:
	@go install github.com/fzipp/gocyclo/cmd/gocyclo@v0.6.0

$(GOBIN)/golangci-lint:
	@go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.3.0

$(GOBIN)/gocritic:
	@go install github.com/go-critic/go-critic/cmd/gocritic@v0.13.0

install-linters: $(GOBIN)/goimports $(GOBIN)/gocyclo $(GOBIN)/golangci-lint $(GOBIN)/gocritic
	@echo "Linters installed successfully."

lint: install-linters
	@pre-commit run -a

clean:
	@rm -rfv bin

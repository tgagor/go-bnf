VERSION ?= $(shell git describe --tags --always)
GOBIN ?= $(shell go env GOPATH)/bin

bin/bnf: build

run:
	go run \
		-ldflags="-X main.BuildVersion=$(VERSION)" \
		. \
		-g examples/numbers.bnf -i examples/numbers.test
	cat examples/numbers.test | go run \
		-ldflags="-X main.BuildVersion=$(VERSION)" \
		. \
		-g examples/numbers.bnf

build:
	go build \
		-ldflags="-w -s -X main.BuildVersion=$(VERSION)" -gcflags=all="-l -B" \
		-o bin/bnf .
	@du -sh bin/bnf

test:
	go test -v ./...

integration-test: bin/bnf
	./bin/bnf -l -g examples/numbers.bnf -i examples/numbers.test
	./bin/bnf -l -g examples/hour.bnf -i examples/hour.test
	cat examples/postal1.txt | ./bin/bnf -g ./examples/postal.bnf

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

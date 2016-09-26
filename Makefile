COMMIT = $$(git describe --always)

.PHONY: all

all: build

help: ## This help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: deps ## Build CLI (./sqsc binary)
	go build -o sqsc -ldflags "-X main.GitCommit=\"$(COMMIT)\""

clean: ## Clean repository
	go clean
	rm -f sqsc coverage.out

deps: ## Install dependencies inside $GOPATH
	go get ./...


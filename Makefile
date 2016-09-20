help: ## This help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: deps ## Build CLI (./sqsc binary)
	go build -o sqsc

clean: ## Clean repository
	go clean
	rm -f sqsc coverage.out

#coverage: deps ## Generate test coverage
#	go test -coverprofile=coverage.out ./...
#	go tool cover -html=coverage.out

deps: ## Install dependencies inside $GOPATH
	go get ./...

.PHONY: all help

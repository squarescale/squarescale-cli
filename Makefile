COMMIT = $$(git describe --always)

.PHONY: all

all: build

help: ## This help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: deps ## Build CLI (./sqsc binary)
	go build -o sqsc -ldflags "-X main.GitCommit=\"$(COMMIT)\""

docker:
	GOOS=linux CGO_ENABLED=0 go build -o sqsc-docker -ldflags "-X main.GitCommit=\"$(COMMIT)\""
	docker build -t sqsc-cli .

clean: ## Clean repository
	go clean
	rm -f sqsc sqsc-docker

deps: ## Install dependencies inside $GOPATH
	go get ./...

lint: ## Lint Docker
	docker run --rm -v $$PWD:/root/ projectatomic/dockerfile-lint dockerfile_lint
	hadolint Dockerfile

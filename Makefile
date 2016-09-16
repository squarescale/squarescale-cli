help: ## This help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

cli: deps ## Build CLI
	GOPATH=$$PWD go install github.com/squarescale/squarescale-cli/sqsc
linux:
	GOPATH=$$PWD CGO_ENABLED=0 GOOS=linux go install github.com/squarescale/squarescale-cli/sqsc
docker: linux
	docker build -t sqsc .
deps:
	GOPATH=$$PWD go get github.com/Bowery/prompt
	GOPATH=$$PWD go get github.com/squarescale/go-netrc/netrc

.PHONY: all help

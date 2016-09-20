help: ## This help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

cli: ## Build CLI
	GOPATH=$$PWD go install github.com/squarescale/squarescale-cli/sqsc

.PHONY: all help

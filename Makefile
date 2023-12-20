default: help

.PHONY: help
help: ## Print this help message
	@echo "Available make commands:"; grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: lint
lint: ## Run golangci-lint
	golangci-lint run --out-format=tab

.PHONY: lint-fix
lint-fix: ## Run golangci-lint with the --fix flag to fix linter errors
	golangci-lint run --fix --out-format=tab --issues-exit-code=0

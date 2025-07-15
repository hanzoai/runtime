# Hanzo Runtime Makefile

.PHONY: help
help: ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

.PHONY: build
build: build-python build-typescript ## Build all packages

.PHONY: build-python
build-python: ## Build Python package
	@echo "Building Python package..."
	@cd libs/sdk-python && python -m build

.PHONY: build-typescript
build-typescript: ## Build TypeScript package
	@echo "Building TypeScript package..."
	@cd libs/sdk-typescript && npm pack

.PHONY: clean
clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	@rm -rf libs/sdk-python/dist libs/sdk-python/build libs/sdk-python/*.egg-info
	@rm -f libs/sdk-typescript/*.tgz
	@rm -f *.tgz

.PHONY: publish
publish: publish-check build publish-python publish-npm ## Build and publish all packages

.PHONY: publish-check
publish-check: ## Check if ready to publish
	@echo "Checking publishing prerequisites..."
	@which twine > /dev/null || (echo "Error: twine not installed. Run: pip install twine" && exit 1)
	@npm whoami > /dev/null 2>&1 || (echo "Error: Not logged in to npm. Run: npm login" && exit 1)
	@echo "✓ All publishing tools available"

.PHONY: publish-python
publish-python: ## Publish Python package to PyPI
	@echo "Publishing Python package to PyPI..."
	@cd libs/sdk-python && python -m twine upload dist/* --repository pypi || echo "Note: You may need to configure PyPI credentials"

.PHONY: publish-npm
publish-npm: ## Publish TypeScript package to npm
	@echo "Publishing TypeScript package to npm..."
	@cd libs/sdk-typescript && npm publish --access public --tag beta || echo "Note: You may need to provide OTP with --otp=<code>"

.PHONY: publish-dry-run
publish-dry-run: build ## Dry run of publishing (no actual upload)
	@echo "Dry run - checking Python package..."
	@cd libs/sdk-python && python -m twine check dist/*
	@echo "Dry run - checking TypeScript package..."
	@cd libs/sdk-typescript && npm pack --dry-run

.PHONY: test
test: ## Run tests
	@echo "Running tests..."
	@echo "Note: Tests not yet configured"

.PHONY: install-deps
install-deps: ## Install development dependencies
	@echo "Installing Python build dependencies..."
	@pip install --upgrade pip build twine
	@echo "Installing Node.js dependencies..."
	@npm install

.PHONY: version
version: ## Show current package versions
	@echo "Package versions:"
	@echo -n "  Python: "
	@grep "^version = " libs/sdk-python/pyproject.toml | cut -d'"' -f2
	@echo -n "  TypeScript: "
	@grep '"version":' libs/sdk-typescript/package.json | cut -d'"' -f4
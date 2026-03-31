.PHONY: check test build lint lint-fix format format-check install deploy_prod

# Run svelte-check in both consumer and admin (type-check only, no artifacts)
check:
	npm run check

# Build both apps for production
build:
	npm run build

# Run unit tests in both consumer and admin
test:
	npm run test

# Lint both apps
lint:
	npm run lint

# Lint with auto-fix
lint-fix:
	npm run lint:fix

# Format all code
format:
	npm run format

# Check formatting only
format-check:
	npm run format:check

# Install dependencies (from frontend root)
install:
	npm install

# Deploy frontend (consumer + admin) to prod. Run from repo root: make deploy_prod_frontend
deploy_prod:
	$(MAKE) -C .. deploy_prod_frontend

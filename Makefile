GOLANGCI_LINT_VERSION ?= v2.11.4

run:
	go run ./cmd/server

test:
	go test ./...

test-integration:
	go test -tags=integration ./internal/modules/order/infrastructure/postgres

fmt:
	gofmt -w ./cmd ./internal

fmt-check:
	@files=$$(gofmt -l ./cmd ./internal); \
	if [ -n "$$files" ]; then \
		echo "unformatted files:"; \
		echo "$$files"; \
		exit 1; \
	fi

vet:
	go vet ./...

lint:
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION) run

new-module:
	@test -n "$(name)" || (echo "usage: make new-module name=invoice" && exit 1)
	@./scripts/new_module.sh $(name)

up:
	docker compose up -d

down:
	docker compose down

logs:
	docker compose logs -f postgres

migrate-up:
	go run ./cmd/migrate up

migrate-down:
	go run ./cmd/migrate down

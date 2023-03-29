.PHONY: lint
lint:
	golangci-lint -v run --fix ./...

.PHONY: gen
gen:
	swag i -g swagger.go --parseInternal --parseDependency --dir ./internal

.PHONY: mock
mock:
	mockery --all --output internal/mocks

.PHONY: test
test:
	go test --race --vet= ./... -v
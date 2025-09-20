SERVER_MAIN=cmd/main.go

.PHONY: run
run: ## runs the server
	go run ${SERVER_MAIN}

.PHONY: fmt
fmt: ## runs go formatter
	go fmt ./...

.PHONY: tidy
tidy: ## runs tidy to fix go.mod dependencies
	go mod tidy

.PHONY: test
test: ## runs tests
	go test -v ./...
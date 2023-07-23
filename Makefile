.PHONY: help go unary server install
help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
open: ## Open DevContaienr
	@devcontainer open .
up: ## Start Container
	@docker-compose up -d
down: ## Stop Container
	@docker-compose down
exec: ## Login Container
	@docker-compose exec go bash
install: ## Install Go Plugin for Protobuf
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28 && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
ps: ## Check Container Status
	@docker-compose ps
protoc: ## Proto-Compiler for Go 
	@protoc -I. --go_out=. proto/*.proto
mod: ## Install modules
	@go mod tidy
unary: ## Compile Unary RPC file
	@protoc -I. --go_out=. --go-grpc_out=. proto/*.proto
server: ## Start Server
	@go run server/main.go
client: ## Start Client
	@go run client/main.go

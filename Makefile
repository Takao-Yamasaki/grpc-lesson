.PHONY: go unary
open:# Open DevContaienr
	@devcontainer open .
up: # Start Container
	@docker-compose up -d
down: # Stop Container
	@docker-compose down
exec: # Login Container
	@docker-compose exec go bash
install: # Install Go Plugin for Protobuf
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28 && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
ps: # Check Container Status
	@docker-compose ps
protoc: #Proto-Compiler for Go 
	@protoc -I. --go_out=. proto/*.proto
mod: # Install modules
	@go mod tidy
unary: # Compile Unary RPC file
	@protoc -I. --go_out=. --go-grpc_out=. proto/*.proto

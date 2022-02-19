run:
	@go run .

tools:
	go install github.com/vektra/mockery/v2@latest

generate:
	go generate ./...
	mockery --all

test:
	go test -v --cover ./...
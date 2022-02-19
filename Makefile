run:
	@go run .

tools:
	go install github.com/vektra/mockery/v2@latest

generate:
	go generate ./...

test:
	go test -v --cover ./...
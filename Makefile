format:
	go fmt ./...

linter: 
	golangci-lint run

test:
	go test -v ./...
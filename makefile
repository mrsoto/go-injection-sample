vet:
	go vet ./...
	shadow ./...

.PHONY:vet

linter:
	golangci-lint run

.PHONY:linter

run: linter
	go run .

.PHONY: run
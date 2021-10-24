vet:
	go vet ./...
	shadow ./...

.PHONY:vet

linter:
	golangci-lint run

.PHONY:linter

build: linter
	go build -o server main.go

.PHONY: build

run: linter build
	./server

.PHONY: run

watch: 
	ulimit -n 100
	reflex -s -g '**/*.go' -G '**/*_test.go' make run

installReflect:
	go install github.com/cespare/reflex@latest

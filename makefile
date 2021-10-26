OUT=server

all: clean vet build linter test

.PHONY: all

vet:
	go vet ./...
	shadow ./...

.PHONY: vet

linter:
	golangci-lint run

.PHONY: linter

test:
	go test -coverprofile=cover.out -timeout 10s ./...

.PHONY: test

cover: test
	go tool cover -html=cover.out

.PHONY: cover

build: 
	go build -o ${OUT} main.go

.PHONY: build

run: build
	./server

.PHONY: run

clean:
	go clean
	rm -f ${OUT}
	rm -f cover.out

watch: 
	ulimit -n 100
	reflex -s -g '**/*.go' -G '**/*_test.go' make run

watch_test: 
	ulimit -n 100
	reflex -s -g '**/*.go' -G '**/*_test.go' make test

deps: installReflect

installReflect:
	go install github.com/cespare/reflex@latest

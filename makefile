OUT=server

all: ${OUT}

.PHONY: all clean run

vet:
	go vet ./...
	shadow ./...

linter:
	golangci-lint run

test:
	go test -coverprofile=cover.out -timeout 10s ./...

race:
	go test -race -timeout 10s ./...

cover: test
	go tool cover -html=cover.out
	
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

${OUT}: build

build: 
	go build -o ${OUT} cmd/server/main.go

run: ${OUT}
	./${OUT}

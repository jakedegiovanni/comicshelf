LINT_VERSION := 1.54.2

.PHONY: all
all: clean install lint generate test build

.PHONY: install
install:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v${LINT_VERSION}

.PHONY: lint
lint:
	golangci-lint run -v ./...

.PHONY: clean
clean:
	rm -rf bin

.PHONY: build
build: clean
	go build -v -x -o bin/ ./cmd/comicshelf

.PHONY: run
run:
	go run ./cmd/comicshelf server $(ARGS)

.PHONY: run-cli
run-cli:
	go run ./cmd/comicshelf $(ARGS)

.PHONY: test
test:
	go test -v ./...

.PHONY: generate
generate:
	go generate ./...

.PHONY: tidy
tidy:
	go mod tidy -v
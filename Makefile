ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
$(eval $(ARGS):;@:)

BUILD_DIR := bin
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
	rm -rf ${BUILD_DIR}

.PHONY: build
build: clean
	go build -v -x -o ${BUILD_DIR}/ ./cmd/...

.PHONY: dev
dev: 
	docker compose build
	docker compose up -d

.PHONY: down
down:
	docker compose down --remove-orphans

.PHONY: run
run:
	go run ./cmd/comicshelf server ${ARGS}

.PHONY: run-cli
run-cli:
	go run ./cmd/comicshelf ${ARGS}

.PHONY: test
test:
	go test -v ./...

.PHONY: generate
generate:
	go generate ./...

.PHONY: tidy
tidy:
	go mod tidy -v
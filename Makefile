.PHONY: clean
clean:
	rm -rf bin

.PHONY: build
build: clean
	go build -v -x -o bin/ .

.PHONY: run
run:
	go run .

.PHONY: generate
generate:
	go generate ./...

.PHONY: tidy
tidy:
	go mod tidy -v
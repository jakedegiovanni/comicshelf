clean:
	rm -rf bin

build:
	go build -o bin/ .

run:
	go run .

generate:
	go generate ./...
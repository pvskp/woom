all: build

build:
	go build ./...

run:
	go run ./...

clean:
	rm ./woom

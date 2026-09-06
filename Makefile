all: build

build:
	go build -o woom ./cmd

run: build
	./woom

clean:
	rm ./woom

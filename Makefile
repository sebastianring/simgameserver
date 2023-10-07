build: test
	go build -o bin/simgameserver

run: build
	./bin/simgameserver

test:
	go test -v ./...

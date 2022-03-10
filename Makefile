all: ./examples/simple/main

./examples/simple/main: ./*.go ./examples/simple/*.go
	cd ./examples/simple; \
	go build main.go

test: ./examples/simple/main
	./examples/simple/main

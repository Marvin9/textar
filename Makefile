bench:
	go test -bench=. ./pkg

.PHONY: build
build:
	go build -o ./bin/textar

bin: build
	./bin/textar
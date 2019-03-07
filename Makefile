.PHONY: build
build:
	go build -o ./cmd/gohealthz/gohealthz ./cmd/gohealthz

.PHONY: clean
clean:
	rm -rf ./cmd/gohealthz/gohealthz

.PHONY: run
run:
	./cmd/gohealthz/gohealthz -timeout=800ms -interval=5m

.PHONY: test
test:
	go test ./...
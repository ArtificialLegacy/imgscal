
.PHONY: build
build-windows:
	rm -rf build/
	go build -o build/imgscal.exe ./cmd/imgscal

.PHONY: build
build-linux:
	rm -rf build/
	go build -o build/imgscal ./cmd/imgscal

.PHONY: examples
examples:
	go run ./cmd/imgscal-examples

start:
	go run ./cmd/imgscal

dev:
	make examples
	make start

.PHONY: doc
doc:
	go run ./cmd/imgscal-doc

.PHONY: doc-open
doc-open:
	go run ./cmd/imgscal-doc
	open ./docs/index.html

.PHONY: types
types:
	go run ./cmd/imgscal-types

.PHONY: log
log:
	go run ./cmd/imgscal-log

.PHONY: new
new:
	go run ./cmd/imgscal-new

install-slim:
	go install ./cmd/imgscal

install:
	go install ./cmd/imgscal
	go install ./cmd/imgscal-new
	go install ./cmd/imgscal-entrypoint
	go install ./cmd/imgscal-log


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

.PHONY: types
types:
	go run ./cmd/imgscal-types

.PHONY: log
log:
	go run ./cmd/imgscal-log

.PHONY: new
new:
	go run ./cmd/imgscal-new


.PHONY: build
build-windows:
	rm -rf build/
	go build -o build/imgscal.exe ./cmd/cli/

.PHONY: build
build-linux:
	rm -rf build/
	go build -o build/imgscal ./cmd/cli

.PHONY: install-examples
install-examples:
	go run ./cmd/install_examples

start:
	go run ./cmd/cli/

dev:
	make install-examples
	make start

doc:
	go run ./cmd/doc

.PHONY: log
log:
	go run ./cmd/log
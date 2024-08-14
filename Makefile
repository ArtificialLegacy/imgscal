
.PHONY: build
build-windows:
	rm -rf build/
	go build -o build/imgscal.exe ./cmd/cli/
	make doc
	mkdir build/docs/
	cp docs/*.md build/docs/

.PHONY: build
build-linux:
	rm -rf build/
	go build -o build/imgscal ./cmd/cli
	make doc
	mkdir build/docs/
	cp docs/*.md build/docs/

.PHONY: install-examples
install-examples:
	go run ./cmd/install_examples

start:
	go run ./cmd/cli/

doc:
	go run ./cmd/doc

.PHONY: log
log:
	go run ./cmd/log
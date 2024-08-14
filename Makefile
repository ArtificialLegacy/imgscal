
.PHONY: build
build-windows:
	rm -rf build/
	go build -o build/imgscal.exe ./cmd/cli/
	mkdir build/workflows/
	cp workflows/*.lua build/workflows/
	cp assets/*.png build/assets/
	make doc
	mkdir build/docs/
	cp docs/*.md build/docs/

.PHONY: build
build-linux:
	rm -rf build/
	go build -o build/imgscal ./cmd/cli
	mkdir build/workflows/
	cp workflows/*.lua build/workflows/
	mkdir build/assets/
	cp assets/*.png build/assets/
	make doc
	mkdir build/docs/
	cp docs/*.md build/docs/

start:
	go run ./cmd/cli/

doc:
	go run ./cmd/doc

.PHONY: log
log:
	go run ./cmd/log

.PHONY: build
build-windows:
	rm -rf build/
	go build -o build/imgscal.exe ./cmd/imgscal/
	mkdir build/workflows/
	cp workflows/*.lua build/workflows/

start:
	go run ./cmd/imgscal/
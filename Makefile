
.PHONY: build
build-windows:
	rm -rf build/
	go build -o build/imgscal.exe
	mkdir build/workflows/
	cp workflows/*.lua build/workflows/

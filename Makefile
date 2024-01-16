
.PHONY: build
build-windows:
	go build -o build/imgscal.exe
	mkdir build/workflows/
	cp workflows/*.lua build/workflows/

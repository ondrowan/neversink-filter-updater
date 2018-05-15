clean:
	-rm -r ./build

build: clean
	mkdir ./build
	GOOS=windows GOARCH=amd64 go build -a -o build/neversink-filter-updater-64-bit.exe
	GOOS=windows GOARCH=386 go build -a -o build/neversink-filter-updater-32-bit.exe

.PHONY: build clean
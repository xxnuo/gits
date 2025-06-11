build:
	mkdir -p dist
	GOOS=linux GOARCH=amd64 go build -o dist/gits-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build -o dist/gits-linux-arm64 .
	GOOS=darwin GOARCH=amd64 go build -o dist/gits-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -o dist/gits-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build -o dist/gits-windows-amd64.exe .
	GOOS=windows GOARCH=arm64 go build -o dist/gits-windows-arm64.exe .
	find dist -type f -name "gits-*" -not -name "*darwin*" -exec upx {} \;
	chmod +x dist/*

default: build
.PHONY: build default
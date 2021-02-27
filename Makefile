.PHONY: dev
dev:
	rm -rf bin
	mkdir -p bin
	go build -ldflags "-X github.com/ghokun/appletv3-iptv/internal/config.Version=Development" -o bin/appletv3-iptv-darwin

.PHONY: release
release:
	VERSION := $(shell git describe --tags --dirty)
	LDFLAGS := "-X github.com/ghokun/appletv3-iptv/internal/config.Version=$(VERSION)"
	mkdir -p bin
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags $(LDFLAGS) -a -o bin/appletv3-iptv
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -ldflags $(LDFLAGS) -a -o bin/appletv3-iptv-armhf
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags $(LDFLAGS) -a -o bin/appletv3-iptv-arm64
	CGO_ENABLED=0 GOOS=darwin go build -ldflags $(LDFLAGS) -a -o bin/appletv3-iptv-darwin
	CGO_ENABLED=0 GOOS=windows go build -ldflags $(LDFLAGS) -a -o bin/appletv3-iptv.exe
	cd bin && shasum -a 256 * > checksum
VERSION := $(shell git describe --tags --dirty)
LDFLAGS := "-X github.com/ghokun/appletv3-iptv/internal/config.Version=$(VERSION)"
SOURCE_DIRS := internal main.go

.PHONY: check-fmt
check-fmt:
	@test -z $(shell gofmt -l -s $(SOURCE_DIRS) ./ | tee /dev/stderr) || (echo "[WARN] Fix formatting issues with 'go fmt'" && exit 1)

.PHONY: build
build:
	rm -rf bin
	mkdir -p bin
	go build -o bin/appletv3-iptv

.PHONY: copy-sample-config
copy-sample-config:
	cp sample/config.yaml bin/config.yaml

.PHONY: dev
dev: check-fmt build copy-sample-config

.PHONY: release
release:
	mkdir -p bin
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags $(LDFLAGS) -a -o bin/appletv3-iptv
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -ldflags $(LDFLAGS) -a -o bin/appletv3-iptv-armhf
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags $(LDFLAGS) -a -o bin/appletv3-iptv-arm64
	CGO_ENABLED=0 GOOS=darwin go build -ldflags $(LDFLAGS) -a -o bin/appletv3-iptv-darwin
	CGO_ENABLED=0 GOOS=windows go build -ldflags $(LDFLAGS) -a -o bin/appletv3-iptv.exe
	cd bin && shasum -a 256 * > checksum
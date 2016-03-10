BIN = agent aggregator graph hbs judge nodata query sender task transfer fe
TARGET = open-falcon

GOTOOLS = github.com/mitchellh/gox golang.org/x/tools/cmd/stringer \
	github.com/jteeuwen/go-bindata/... github.com/elazarl/go-bindata-assetfs/...
PACKAGES=$(shell go list ./... | grep -v '^github.com/Cepave/open-falcon/')
VERSION?=$(shell awk -F\" '/^const Version/ { print $$2; exit }' version.go)

all: $(BIN)
	go build -o open-falcon

$(BIN)):
	go build -o ./bin/$@ github.com/cepave/open-falcon/modules/$@

# dev creates binaries for testing locally - these are put into ./bin and $GOPATH
dev: format
	@CONSUL_DEV=1 sh -c "'$(CURDIR)/scripts/build.sh'"

test: format
	@./scripts/test.sh

format:
	@echo "--> Running go fmt"
	@go fmt $(PACKAGES)

tools:
	go get -u -v $(GOTOOLS)

pack:
	tar -zcvf open-falcon-v$(VERSION).tar.gz ./bin ./config ./open-falcon ./cfg.json

clean:
	rm -rf ./bin
	rm -rf ./$(TARGET)

.PHONY: agent aggregator graph hbs judge nodata query sender task transfer fe clean all

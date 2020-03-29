NOW = $(shell date +%s)

GOCMD=go
GOBUILD=$(GOCMD) build -mod vendor
GOTEST=$(GOCMD) test -mod vendor
GORUN=$(GOCMD) run main.go
THRIFT_CMD=thrift
THRIFT_GO_BASE_CMD=$(THRIFT_CMD) -out thrift/ --gen go:thrift_import=github.com/apache/thrift/lib/go/thrift,package_prefix={{package}}/thrift/


.PHONY: build build-prod test run thrift


build: thrift
	go mod vendor
	$(GOBUILD) ./...

clean:
	rm -f {{project}}
	rm -rf thrift/{{owner}}

build-prod:
	$(GOBUILD) -trimpath -o /build/service .

test: build
	$(GOTEST) ./...

run:
	$(GORUN) -debug -config dev.yaml

thrift:
	mkdir -p thrift
	curl -o idl/baseplate.thrift https://raw.githubusercontent.com/reddit/baseplate.py/master/baseplate/thrift/baseplate.thrift
	$(THRIFT_GO_BASE_CMD) idl/baseplate.thrift
	$(THRIFT_GO_BASE_CMD) idl/{{project}}.thrift

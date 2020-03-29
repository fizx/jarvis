NOW = $(shell date +%s)
THRIFT_CMD=thrift
THRIFT_GO_BASE_CMD=$(THRIFT_CMD) -out generated/thrift/ --gen go:thrift_import=github.com/apache/thrift/lib/go/thrift,package_prefix=github.com/fizx/jarvis/generated/thrift/

default: thrift
	mkdir -p generated/assets
	go generate .
	go build .
	go build ./pkg/jarvis

thrift:
	mkdir -p generated/thrift
	$(THRIFT_GO_BASE_CMD) idl/baseplate.thrift

	
test: default
	rm -rf iron_man
	./jarvis github.com/fizx/iron-man
	echo 'replace github.com/fizx/jarvis => ../' >> iron_man/go.mod
	echo 'replace github.com/fizx/iron-man => ./' >> iron_man/go.mod
	cd iron* && make test

clean:
	rm -rf generated
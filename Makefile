NOW = $(shell date +%s)

default:
	mkdir -p generated/assets
	go generate .
	go build .
	
test: default
	rm -rf iron_man
	./jarvis github.com/fizx/iron-man
	echo 'replace github.com/fizx/jarvis => ../' >> iron_man/go.mod
	echo 'replace github.com/fizx/iron-man => ./' >> iron_man/go.mod
	cd iron* && make test

clean:
	rm -rf generated
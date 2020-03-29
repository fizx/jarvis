NOW = $(shell date +%s)

default:
	mkdir -p generated/assets
	go generate .
	go build .
	
test: default
	rm -rf iron*
	./jarvis github.com/fizx/iron-man
	cd iron* && make test

clean:
	rm -rf generated
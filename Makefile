NOW = $(shell date +%s)

default:
	mkdir -p generated/assets
	go generate .
	go build .
	
test: default
	rm -rf iron-man
	./jarvis create iron-man
	cd iron-man && make test

clean:
	rm -rf generated
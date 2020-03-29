NOW = $(shell date +%s)

default:
	mkdir -p generated/assets
	go generate .
	go build .
	
test: default
	rm -rf iron_man
	./jarvis iron-man
	cd iron_man && go mod init github.com/fizx/iron_man && make test

clean:
	rm -rf generated
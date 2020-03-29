NOW = $(shell date +%s)

default:
	go build ./...
	
test: default
	go test ./...
	
clean:
	rm -f {{project}}
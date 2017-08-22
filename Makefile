all:  dep build 
	@echo "All done. \o/"

build: test
	@rm -rf build
	@go build -o "build/docklog" src/docklog.go
	
dep:
	@echo "getting dependency tool"
	@go get -u github.com/golang/dep/cmd/dep
	@echo "update dependencies"
	@dep ensure

lint: 
	@echo "verify src with go vet"
	@go tool vet -composites=false -shadow=true src/**/*.go

test: lint
	@echo "let's doing some tests"
	@go test -race ./src/...

.PHONY: test dep build lint
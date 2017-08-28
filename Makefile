BIN_DIR := $(GOPATH)/bin
GOMETALINTER := $(BIN_DIR)/gometalinter

.PHONY: test dep build lint megalint

all:  dep build 
	@echo "All done. \o/"

build: test
	@rm -rf build
	@go build -o "build/docklog" docklog.go
	
dep:
	@echo "getting dependency tool"
	@go get -u github.com/golang/dep/cmd/dep
	@rm -rf ./vendor
	@echo "update dependencies"
	@dep ensure

lint: 
	@echo "verify src with go vet"
	@go tool vet -composites=false -shadow=true *.go
	@go tool vet -composites=false -shadow=true tools/*.go

megalint: $(GOMETALINTER)
	gometalinter ./... --vendor

test: lint
	@echo "let's doing some tests"
	@go test -race ./...


$(GOMETALINTER):
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install &> /dev/null




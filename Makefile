all: deps cmd

deps:
	- go get "github.com/docker/docker/api/types"
	- go get "github.com/docker/docker/client"
	- go get "github.com/google/go-jsonnet/cmd/jsonnet"
	- go get "gopkg.in/yaml.v2"
	- go get "github.com/docker/go-connections/nat"
	- rm -rf $(HOME)/go/src/rm -rf ../github.com/docker/docker/vendor/github.com/docker/go-connections/nat

cmd:
	- go build -o build/combu main/main.go

run:
	- ./build/combu -f config/config.jsonnet run

run-test:
	- go test combu/test

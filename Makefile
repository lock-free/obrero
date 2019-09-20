GOPATH := $(shell cd ../../../.. && pwd)
export GOPATH

init-dep:
	@dep init

dep:
	@dep ensure

status-dep:
	@dep status

update-dep:
	@dep ensure -update

test:
	@go test -v -race
	@cd ./utils && go test -v -race
	@cd ./utils/cq && go test -v -race

cover:
	@go test -coverprofile=coverage.out && go tool cover -html=coverage.out
	@cd ./utils && go test -coverprofile=coverage.out && go tool cover -html=coverage.out
	@cd ./utils/cq && go test -coverprofile=coverage.out && go tool cover -html=coverage.out

.PHONY: test

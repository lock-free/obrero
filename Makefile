GO111MODULE := on
export GO111MODULE

init:
	@go mod init

clean:
	@go mod tidy

update:
	@go get -u

test:
	@go test ./... -v -race

cover:
	@go test -coverprofile=coverage.out && go tool cover -html=coverage.out
	@cd ./utils && go test -coverprofile=coverage.out && go tool cover -html=coverage.out
	@cd ./utils/cq && go test -coverprofile=coverage.out && go tool cover -html=coverage.out
	@cd ./utils/dlb && go test -coverprofile=coverage.out && go tool cover -html=coverage.out
	@cd ./mids && go test -coverprofile=coverage.out && go tool cover -html=coverage.out
	@cd ./mids/httpmids && go test -coverprofile=coverage.out && go tool cover -html=coverage.out
	@cd ./box/cqbox && go test -coverprofile=coverage.out && go tool cover -html=coverage.out

.PHONY: test

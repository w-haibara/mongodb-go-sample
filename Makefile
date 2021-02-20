mongodb-go-sample: *.go */*.go go.mod
	gofmt -w *.go */*.go
	go build -o mongodb-go-sample main.go

.PHONY: init
init:
	go mod init mongodb-go-sample

.PHONY: test
test:
	gofmt -w *.go
	go test

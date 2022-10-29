
 
.PHONY: proto
proto:
	sudo docker run --rm -v $(pwd):/go -w /go poabob/protoc-builder:latest sh -c 'protoc --proto_path=./protos --micro_out=./protos --go_out=./protos ./protos/base/base.proto'
# 要修改
.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o base *.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t base:latest

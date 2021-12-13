test:
	go test ./...

fmt:
	go fmt ./...

build: fmt
	CGO_ENABLED=0 go build -mod vendor -o ./zaplogw contrib-linter/cmd/zap/logw


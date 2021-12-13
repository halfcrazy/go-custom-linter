fmt:
	go fmt ./...

build: fmt
	CGO_ENABLED=0 GOOS=linux go build -mod vendor -o ./logw ./main.go

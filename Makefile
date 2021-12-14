test-dep:
	cp -r ./vendor/* ./pkg/passes/gorm/addressable/testdata/src/

clean:
	git clean -fd

test: test-dep
	go test ./...

fmt:
	go fmt ./...

build: fmt
	CGO_ENABLED=0 go build -mod vendor -o ./zaplogw go-custom-linter/cmd/zap/logw
	CGO_ENABLED=0 go build -mod vendor -o ./gormaddressable go-custom-linter/cmd/gorm/addressable


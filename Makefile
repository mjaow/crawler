.PHONY: build
build: fmt vet
	go build -o crawler main.go

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

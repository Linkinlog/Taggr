build:
	CGO_ENABLED=0 GOOS=linux go build -o tag -ldflags "-s -w" ./...

.PHONY: build

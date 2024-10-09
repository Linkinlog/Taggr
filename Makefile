build:
	CGO_ENABLED=0 GOOS=linux go build -o tag -ldflags "-s -w"
gen:
	go generate ./...
lint:
	golangci-lint run

.PHONY: build gen lint

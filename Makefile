IMAGE = 7phs/area-51
VERSION = latest

.PHONY: build
build:
	go build -o ./bin/ ./cmd/...

.PHONY: build-race
build-race:
	go build --race -o ./bin/ ./cmd/...

.PHONY: run-server
run:
	go run ./cmd/server

.PHONY: dep
dep:
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.42.0

.PHONY: test
test:
	go test --race ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: image
image:
	docker build -t $(IMAGE):$(VERSION)  .

.PHONY: push
push:
	docker push $(IMAGE):$(VERSION)

SWAG_CMD = swag init --generalInfo cmd/server/main.go --output docs

.PHONY: all
all: build

.PHONY: build
build:
	go build -o bin/app ./cmd/server

.PHONY: serve
run:
	go run ./cmd/server

.PHONY: dev
dev:
	air

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: docs
docs:
	$(SWAG_CMD)

.PHONY: clean
clean:
	rm -rf bin/
	rm -rf tmp/
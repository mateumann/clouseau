.PHONY: all
all: lint test build

.PHONY: clean
clean: clean-bin

.PHONY: clean-bin
clean-bin:
	rm -f collector

.PHONY: fmt
fmt: fmt-cmd fmt-pkg fmt-test

.PHONY: fmt-cmd
fmt-cmd:
	gofumpt -s -w cmd/collector/*.go

.PHONY: fmt-pkg
fmt-pkg:
	gofumpt -s -w pkg/collector/*.go
	gofumpt -s -w pkg/encoding/*.go
	gofumpt -s -w pkg/storage/*.go

.PHONY: fmt-test
fmt-test:
	gofumpt -s -w test/*.go

.PHONY: lint
lint:
	golangci-lint run

.PHONY: analyse
analyse:
	gocyclo -over 9 -ignore "_test.go$$" .
	gocognit -over 7 \
		cmd/collector/main.go \
		pkg/collector/collector.go \
		pkg/collector/udp.go \
		pkg/encoding/decode.go \
		pkg/encoding/element.go \
		pkg/encoding/helpers.go \
		pkg/encoding/message.go \
		pkg/storage/storage.go \
		pkg/storage/rethinkdb.go
	go list -json -m all | nancy sleuth

.PHONY: test
test:
	go-acc ./...

.PHONY: integration-test
integration-test:
	go test -tags=integration ./test

.PHONY: build
build: build-collector

.PHONY: build-collector
build-collector:
	go build -o collector cmd/collector/main.go

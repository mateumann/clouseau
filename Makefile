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
	gofumpt -s -w pkg/encoding/ipfix/*.go
	gofumpt -s -w pkg/storage/*.go
	gofumpt -s -w cmd/collector/*.go

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
		pkg/collector/collector.go \
		pkg/collector/udp.go \
		pkg/encoding/ipfix/decode.go \
		pkg/encoding/ipfix/element.go \
		pkg/encoding/ipfix/fieldspecifier.go \
		pkg/encoding/ipfix/helpers.go \
		pkg/encoding/ipfix/message.go \
		pkg/encoding/ipfix/record.go \
		pkg/encoding/ipfix/set.go \
		pkg/encoding/ipfix/template.go \
		pkg/storage/storage.go \
		pkg/storage/rethinkdb.go \
		cmd/collector/main.go
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

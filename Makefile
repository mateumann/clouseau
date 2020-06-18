.PHONY: all
all: fmt lint test build

.PHONY: clean
clean: clean-bin

.PHONY: clean-bin
clean-bin:
	rm -f collector

.PHONY: fmt
fmt:
	gofmt -s -w cmd/collector/*.go
	gofmt -s -w encoding/ipfix/*.go
	goimports -w cmd/collector/*.go
	goimports -w encoding/ipfix/*.go


.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test: export TMPDIR=${XDG_RUNTIME_DIR}/tmp
test: mktmpdir test-ipfix

.PHONY: test-ipfix
test-ipfix:
	go test ./encoding/ipfix/

.PHONY: build
build: build-collector

.PHONY: build-collector
build-collector:
	go build -o collector cmd/collector/main.go

.SILENT: mktmpdir
mktmpdir:
	if [ ! -d ${TMPDIR} ] ; then mkdir ${TMPDIR} ; fi

.PHONY: all
all: fmt lint test build

.PHONY: clean
clean: clean-bin

.PHONY: clean-bin
clean-bin:
	rm -f collector

.PHONY: fmt
fmt: fmt-cmd fmt-ipfix

.PHONY: fmt-cmd
fmt-cmd:
	gofmt -s -w cmd/collector/*.go
	goimports -w cmd/collector/*.go

.PHONY: fmt-ipfix
fmt-ipfix:
	gofmt -s -w encoding/ipfix/*.go
	goimports -w encoding/ipfix/*.go

.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test: export TMPDIR=${XDG_RUNTIME_DIR}/tmp
test: make-tmp-dir test-ipfix create-coverage-report

.PHONY: test-ipfix
test-ipfix:
	go test -cover -coverprofile=ipfix-c.out ./encoding/ipfix/

.PHONY: create-coverage-report
create-coverage-report:
	go tool cover -html=ipfix-c.out -o coverage.html
	rm *-c.out

.PHONY: build
build: build-collector

.PHONY: build-collector
build-collector:
	go build -o collector cmd/collector/main.go

.SILENT: make-tmp-dir
make-tmp-dir:
	if [ ! -d ${TMPDIR} ] ; then mkdir ${TMPDIR} ; fi

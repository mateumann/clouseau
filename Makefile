.PHONY: all
all: lint test build

.PHONY: clean
clean: clean-bin

.PHONY: clean-bin
clean-bin:
	rm -f collector

.PHONY: fmt
#fmt: fmt-cmd fmt-ipfix
fmt: fmt-ipfix

.PHONY: fmt-cmd
fmt-cmd:
	gofumpt -s -w cmd/collector/*.go
	gofumports -w cmd/collector/*.go

.PHONY: fmt-ipfix
fmt-ipfix:
	gofumpt -s -w encoding/ipfix/*.go
	gofumports -w encoding/ipfix/*.go

.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test: export TMPDIR=${XDG_RUNTIME_DIR}/tmp
test: make-tmp-dir test-ipfix create-coverage-report

.PHONY: test-ipfix
test-ipfix:
	echo -n '' > ./encoding/ipfix/coverage.txt
	go test -coverprofile=./encoding/ipfix/coverage.txt -covermode=atomic ./encoding/ipfix/

.PHONY: create-coverage-report
create-coverage-report:
	cat encoding/ipfix/coverage.txt >> coverage.txt
	rm encoding/ipfix/coverage.txt

.PHONY: build
build: build-collector

.PHONY: build-collector
build-collector:
	go build -o collector cmd/collector/main.go

.SILENT: make-tmp-dir
make-tmp-dir:
	if [ ! -d ${TMPDIR} ] ; then mkdir ${TMPDIR} ; fi

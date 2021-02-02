.PHONY: all
all: lint test build

.PHONY: clean
clean: clean-bin

.PHONY: clean-bin
clean-bin:
	rm -f collector

.PHONY: fmt
fmt: fmt-cmd fmt-encoding fmt-net
#fmt: fmt-encoding

.PHONY: fmt-cmd
fmt-cmd:
	gofumpt -s -w cmd/collector/*.go
	gofumports -w cmd/collector/*.go

.PHONY: fmt-encoding
fmt-encoding:
	gofumpt -s -w encoding/ipfix/*.go
	gofumports -w encoding/ipfix/*.go

.PHONY: fmt-net
fmt-net:
	gofumpt -s -w net/collector/*.go
	gofumports -w net/collector/*.go

.PHONY: lint
lint:
	golangci-lint run

.PHONY: analyse
analyse:
	gocyclo -over 9 -ignore "_test.go$$" .
	gocognit -over 7 \
		cmd/collector/main.go \
		encoding/ipfix/decode.go \
		encoding/ipfix/element.go \
		encoding/ipfix/helpers.go \
		encoding/ipfix/message.go \
		net/collector/collector.go

.PHONY: test
test: export TMPDIR=${XDG_RUNTIME_DIR}/tmp
test: make-tmp-dir test-encoding test-net create-coverage-report

.PHONY: test-encoding
test-encoding:
	echo -n '' > ./encoding/ipfix/coverage.txt
	go test -coverprofile=./encoding/ipfix/coverage.txt -covermode=atomic ./encoding/ipfix/

.PHONY: test-net
test-net:
	echo -n '' > ./net/collector/coverage.txt
	go test -coverprofile=./net/collector/coverage.txt -covermode=atomic ./net/collector/

.PHONY: create-coverage-report
create-coverage-report:
	cat encoding/ipfix/coverage.txt >> coverage.txt
	cat net/collector/coverage.txt >> coverage.txt
	rm encoding/ipfix/coverage.txt
	rm net/collector/coverage.txt

.PHONY: build
build: build-collector

.PHONY: build-collector
build-collector:
	go build -o collector cmd/collector/main.go

.SILENT: make-tmp-dir
make-tmp-dir:
	if [ ! -d ${TMPDIR} ] ; then mkdir ${TMPDIR} ; fi

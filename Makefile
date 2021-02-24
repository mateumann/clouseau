.PHONY: all
all: lint test build

.PHONY: clean
clean: clean-bin

.PHONY: clean-bin
clean-bin:
	rm -f collector

.PHONY: fmt
fmt: fmt-cmd
	gofumpt -s -w ipfix/*.go

.PHONY: fmt-cmd
fmt-cmd:
	gofumpt -s -w cmd/collector/*.go

#.PHONY: fmt-encoding
#fmt-encoding:
#	gofumpt -s -w encoding/ipfix/*.go

#.PHONY: fmt-net
#fmt-net:
#	gofumpt -s -w net/ipfix/*.go

.PHONY: lint
lint:
	golangci-lint run

.PHONY: analyse
analyse:
	gocyclo -over 9 -ignore "_test.go$$" .
	gocognit -over 7 \
		cmd/collector/main.go \
		ipfix/decode.go \
		ipfix/element.go \
		ipfix/helpers.go \
		ipfix/message.go \
		ipfix/collector.go
	go list -json -m all | nancy sleuth

.PHONY: test
test: export TMPDIR=${XDG_RUNTIME_DIR}/tmp
#test: make-tmp-dir test-encoding test-net create-coverage-report
test: make-tmp-dir test-ipfix create-coverage-report

#.PHONY: test-encoding
#test-encoding:
#	echo -n '' > ./encoding/ipfix/coverage.txt
#	go test -coverprofile=./encoding/ipfix/coverage.txt -covermode=atomic ./encoding/ipfix/

#.PHONY: test-net
#test-net:
#	echo -n '' > ./net/ipfix/coverage.txt
#	go test -coverprofile=./net/ipfix/coverage.txt -covermode=atomic ./net/ipfix/

.PHONY: test-ipfix
test-ipfix:
	echo -n '' > ./ipfix/coverage.txt
	go test -coverprofile=./ipfix/coverage.txt -covermode=atomic ./ipfix/

.PHONY: create-coverage-report
create-coverage-report:
#	cat encoding/ipfix/coverage.txt >> coverage.txt
#	cat net/ipfix/coverage.txt >> coverage.txt
	cat ipfix/coverage.txt >> coverage.txt
#	rm encoding/ipfix/coverage.txt
#	rm net/ipfix/coverage.txt
	rm ipfix/coverage.txt

.PHONY: build
build: build-collector

.PHONY: build-collector
build-collector:
	go build -o collector cmd/collector/main.go

.SILENT: make-tmp-dir
make-tmp-dir:
	if [ ! -d ${TMPDIR} ] ; then mkdir ${TMPDIR} ; fi

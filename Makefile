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
		pkg/collector/collector_udp.go \
		pkg/encoding/decode.go \
		pkg/encoding/element.go \
		pkg/encoding/helpers.go \
		pkg/encoding/message.go
	go list -json -m all | nancy sleuth

.PHONY: test
test: export TMPDIR=${XDG_RUNTIME_DIR}/tmp
test: make-tmp-dir test-collector test-encoding test-storage test-integration create-coverage-report
test: make-tmp-dir create-coverage-report

.PHONY: test-collector
test-collector:
	echo -n '' > ./pkg/collector/coverage.txt
	go test -coverprofile=./pkg/collector/coverage.txt -covermode=atomic ./pkg/collector/

.PHONY: test-encoding
test-encoding:
	echo -n '' > ./pkg/encoding/coverage.txt
	go test -coverprofile=./pkg/encoding/coverage.txt -covermode=atomic ./pkg/encoding/

.PHONY: test-storage
test-storage:
	echo -n '' > ./pkg/storage/coverage.txt
	go test -coverprofile=./pkg/storage/coverage.txt -covermode=atomic ./pkg/storage/

.PHONY: test-integration
test-integration:
	go test ./test/


.PHONY: create-coverage-report
create-coverage-report:
	cat pkg/collector/coverage.txt >> coverage.txt
	cat pkg/encoding/coverage.txt >> coverage.txt
	cat pkg/storage/coverage.txt >> coverage.txt
	rm pkg/collector/coverage.txt
	rm pkg/encoding/coverage.txt
	rm pkg/storage/coverage.txt

.PHONY: build
build: build-collector

.PHONY: build-collector
build-collector:
	go build -o collector cmd/collector/main.go

.SILENT: make-tmp-dir
make-tmp-dir:
	if [ ! -d ${TMPDIR} ] ; then mkdir ${TMPDIR} ; fi

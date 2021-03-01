.PHONY: all
all: lint test build

.PHONY: clean
clean: clean-bin

.PHONY: clean-bin
clean-bin:
	rm -f collector

.PHONY: fmt
fmt: fmt-cmd fmt-pkg

.PHONY: fmt-cmd
fmt-cmd:
	gofumpt -s -w cmd/collector/*.go

.PHONY: fmt-pkg
fmt-pkg:
	gofumpt -s -w pkg/collecting/*.go
	gofumpt -s -w pkg/encoding/*.go
#	gofumpt -s -w pkg/storage/*.go

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
		pkg/collecting/collector.go
		pkg/encoding/decode.go \
		pkg/encoding/element.go \
		pkg/encoding/helpers.go \
		pkg/encoding/message.go \
	go list -json -m all | nancy sleuth

.PHONY: test
test: export TMPDIR=${XDG_RUNTIME_DIR}/tmp
test: make-tmp-dir test-collecting test-encoding test-storage create-coverage-report
test: make-tmp-dir create-coverage-report

.PHONY: test-collecting
test-collecting:
	echo -n '' > ./pkg/collecting/coverage.txt
	go test -coverprofile=./pkg/collecting/coverage.txt -covermode=atomic ./pkg/collecting/

.PHONY: test-encoding
test-encoding:
	echo -n '' > ./pkg/encoding/coverage.txt
	go test -coverprofile=./pkg/encoding/coverage.txt -covermode=atomic ./pkg/encoding/

.PHONY: test-storage
test-storage:
	echo -n '' > ./pkg/storage/coverage.txt
	go test -coverprofile=./pkg/storage/coverage.txt -covermode=atomic ./pkg/storage/


.PHONY: create-coverage-report
create-coverage-report:
	cat pkg/collecting/coverage.txt >> coverage.txt
	cat pkg/encoding/coverage.txt >> coverage.txt
	cat pkg/storage/coverage.txt >> coverage.txt
#	cat ipfix/coverage.txt >> coverage.txt
	rm pkg/collecting/coverage.txt
	rm pkg/encoding/coverage.txt
	rm pkg/storage/coverage.txt
#	rm ipfix/coverage.txt

.PHONY: build
build: build-collector

.PHONY: build-collector
build-collector:
	go build -o collector cmd/collector/main.go

.SILENT: make-tmp-dir
make-tmp-dir:
	if [ ! -d ${TMPDIR} ] ; then mkdir ${TMPDIR} ; fi

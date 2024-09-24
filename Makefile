.PHONY: lint test


lint:
	find . -name "*.go"  | grep -v mocks | xargs goimports -w
	find . -name "*.go"  | grep -v mocks | xargs gofmt -w
	go vet . ./conf/... ./internal/... ./jkerr/... ./jklog/... ./jkmetric/... ./pkg/...

test:
	go test -gcflags=-l -v -count=1 .
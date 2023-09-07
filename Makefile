all: test

.PHONY: test
test:
	# Don't use parallel because of config tests share the same os environment
	GOARCH=amd64 ginkgo -gcflags=all=-l -r --output-dir=test/report/ -covermode count --coverpkg=./config

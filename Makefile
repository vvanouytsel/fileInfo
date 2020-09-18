.PHONY: run
run:
	docker build . -t fileinfo:latest && docker run --rm fileinfo:latest -v /tmp

.PHONY: test
test:
	docker run -v $(shell pwd):/go/src golang:1.14 /bin/bash -c "cd src/ && go test -cover"

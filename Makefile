PROJECTNAME=$(shell basename "$(PWD)")

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

VERSION="2.0.0-"
COMMIT=`git rev-parse HEAD | cut -c 1-8`
BUILD=`date -u +%Y%m%d.%H%M%S`

compile:
	@-$(MAKE) -s go-compile

release:
	@-$(MAKE) -s go-compile-release

test:
	@-$(MAKE) -s go-test

coverage:
	@-$(MAKE) -s go-test-coverage

run:
	@-$(MAKE) -s go-compile go-run

clean:
	@-$(MAKE) go-clean


go-compile: go-clean go-build

go-compile-release: go-clean go-build-release

go-run:
	@echo "  >  Running application ..."
	./mydms.api

go-test:
	@echo "  >  Go test ..."
	go test -race ./...

go-test-coverage:
	@echo "  >  Go test coverage ..."
	go test -race -coverprofile="coverage.txt" -covermode atomic ./...

go-build:
	@echo "  >  Building binary ..."
	go build -o mydms.api

go-build-release:
	@echo "  >  Building binary..."
	GOOS=linux GOARCH=amd64 go build -ldflags="-w -s -X main.Version=${VERSION}${COMMIT} -X main.Build=${BUILD}" -tags prod -o mydms.api

go-clean:
	@echo "  >  Cleaning build cache"
	go clean ./...
	rm -f ./mydms.api

.PHONY: compile release test run clean coverage
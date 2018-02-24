NAME := lab
VERSION := v0.0.1
REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS := -ldflags="-s -w -X \"main.Version=$(VERSION)\" -X \"main.Revision=$(REVISION)\" -extldflags \"-static\""
DIST_DIRS := find * -type d -exec

.PHONY: dep
dep:
ifeq ($(shell command -v dep 2> /dev/null),)
	go get github.com/golang/dep/...
endif

.PHONY: ensure
ensure: dep
	$(GOPATH)/bin/dep ensure

.PHONY: test
test:
	go test github.com/lighttiger2505/lab/...

.PHONY: build
build:
	go build

.PHONY: install
install:
	go install

.PHONY: coverage
coverage:
	go get -u github.com/haya14busa/goverage
	goverage -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

.PHONY: cross-build
cross-build: ensure
	for os in darwin linux windows; do \
		for arch in amd64 386; do \
			GOOS=$$os GOARCH=$$arch CGO_ENABLED=0 go build -a -tags netgo -installsuffix netgo $(LDFLAGS) -o dist/$$os-$$arch/$(NAME); \
		done; \
	done

.PHONY: dist
dist:
	cd dist && \
	$(DIST_DIRS) cp ../LICENSE {} \; && \
	$(DIST_DIRS) cp ../README.md {} \; && \
	$(DIST_DIRS) tar -zcf $(NAME)-$(VERSION)-{}.tar.gz {} \; && \
	$(DIST_DIRS) zip -r $(NAME)-$(VERSION)-{}.zip {} \; && \
	cd ..

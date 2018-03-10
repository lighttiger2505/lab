NAME := lab
VERSION := 0.1.0
REVISION := $(shell git rev-parse --short HEAD)
GOVERSION := $(go version)

SRCS := $(shell find . -type f -name '*.go')
LDFLAGS := -ldflags="-s -w -X \"main.version=$(VERSION)\" -X \"main.revision=$(REVISION)\" -X \"main.goversion=$(GOVERSION)\" "
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
build: $(SRCS)
	go build $(LDFLAGS) -o bin/$(NAME)

.PHONY: install
install: $(SRCS)
	go install $(LDFLAGS)

.PHONY: coverage
coverage:
	go get -u github.com/haya14busa/goverage
	goverage -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

.PHONY: cross-build
cross-build: ensure
	for os in darwin linux windows; do \
		for arch in amd64 386; do \
			GOOS=$$os GOARCH=$$arch CGO_ENABLED=0 go build $(LDFLAGS) -o dist/$$os-$$arch/$(NAME); \
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

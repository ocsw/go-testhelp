# Copyright 2021-2022 Danielle Zephyr Malament
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Generic Go Makefile

# For each subdirectory of './cmd/' that contains a main.go file, builds a
# binary with the same name as the subdirectory

# GO_TEST allows modifying the test command and options from the command line,
# e.g.:
#
# * go test (the standard)
# * gotest (https://github.com/rakyll/gotest) colorizes standard test output
# * gotestsum (https://github.com/gotestyourself/gotestsum) has
#   color output and condensed formats; in particular, it summarizes failures
#   at the end of the output so you don't have to hunt for them (use
#   'gotestsum GOTESTSUM_ARGS -- GO_TEST_ARGS')
# * add '-count 1' to prevent result caching
# * add -v to go test / gotest, or -f standard-verbose to gotestsum, to print
#   all output, even from successful tests (among other kinds of verbosity)
#
# Sample command-line usage:
#     GO_TEST="gotestsum -f dots -- -count 3" make test
# (omit the quotes in this file)
#
# For CI systems, you probably want just 'go test', e.g.:
#     GO_TEST="go test" make test
GO_TEST ?= gotestsum -f standard-verbose --

# GO_TEST_FILES allows modifying the test scope (files and/or tests) from the
# command line, e.g.:
# GO_TEST_FILES="./cmd/foo -run TestSomethingWorks" make test
GO_TEST_FILES ?= ./...

all: lint test build

# Format imports and remove unused ones; to install:
# go install golang.org/x/tools/cmd/goimports@latest
format:
	@echo ">> Formatting..."
	goimports -w cmd pkg

# See https://github.com/golangci/golangci-lint; to install on macOS:
# brew install golangci/tap/golangci-lint
lint:
	@echo ">> Linting..."
	golangci-lint run \
		-E exportloopref,goimports,lll,revive,stylecheck,whitespace \
		--max-issues-per-linter 0 --max-same-issues 0 $(LINT_ARGS)

test:
	@echo ">> Testing..."
	$(GO_TEST) -v -cover -timeout 20s $(GO_TEST_FILES)

# Testing target that skips the cache (more convenient than setting GO_TEST on
# the command line)
test-nc:
	@echo ">> Testing..."
	$(GO_TEST) -count 1 -v -cover -timeout 20s $(GO_TEST_FILES)

# Testing target that uses gotestsum (more convenient than setting GO_TEST on
# the command line)
# Note the variables for arguments to gotestsum and go test
testsum:
	@echo ">> Testing..."
	gotestsum $(GOTESTSUM_ARGS) -- $(GO_TEST_ARGS) -cover -timeout 20s \
		$(GO_TEST_FILES)

# Testing target that uses gotestsum and skips the cache (more convenient than
# setting GO_TEST on the command line)
# Note the variables for arguments to gotestsum and go test
testsum-nc:
	@echo ">> Testing..."
	gotestsum $(GOTESTSUM_ARGS) -- $(GO_TEST_ARGS) -count 1 -cover \
		-timeout 20s $(GO_TEST_FILES)

bin:
	@echo ">> Making bin directory..."
	mkdir -p bin

output:
	@echo ">> Making output directory..."
	mkdir -p output

build: bin output
	@echo ">> Building all binaries..."
	for cmddir in cmd/*; do \
		if [ -f "$$cmddir/main.go" ]; then \
			"$(MAKE)" "build-$${cmddir#cmd/}"; \
		fi; \
	done

build-%: bin output FORCE
	bin="$(@:build-%=%)" && \
		GOOS=$$(uname | tr 'A-Z' 'a-z') \
			go build -v -o "bin/$$bin" "cmd/$$bin/main.go"

# This builds one of the available binaries, then runs it.  If the ARGS
# environment variable is set, its contents will be added to the command line
# when running the binary.  For example, 'ARGS="foo bar" make run-bin1' will
# run 'bin1 foo bar' after building bin1.
# This target was added to make it more convenient to iterate on code and run
# it, without forgetting to rebuild each time.
run-%: output build-% FORCE
	"bin/$(@:run-%=%)" $$ARGS

clean:
	@echo ">> Removing all binaries..."
	for cmddir in cmd/*; do \
		if [ -f "$$cmddir/main.go" ]; then \
			"$(MAKE)" "clean-$${cmddir#cmd/}"; \
		fi; \
	done

clean-%: FORCE
	rm -rf "bin/$(@:clean-%=%)"

fullclean:
	@echo ">> Removing all binaries and output files, including directories..."
	rm -rf bin output

# Including build-% doesn't make build-foo phony, so use the FORCE trick instead
.PHONY: all format lint lint-full test test-nc testsum testsum-nc
.PHONY: build clean fullclean
FORCE: ;

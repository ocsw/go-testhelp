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

# For each subdirectory of cmd/ that contains a main.go file, builds a
# binary with the same name as the subdirectory


#################
# include files #
#################

# Note: included files are appended to MAKEFILE_LIST and processed in that
# order; the top-level file will always be first

# Pull in repo-specific additions; ignore missing file
-include Makefile.local.mk


#############
# variables #
#############

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
#
# For CI systems, you probably want just 'go test', e.g.:
#     GO_TEST="go test" make test
GO_TEST ?= gotestsum -f standard-verbose --

# GO_TEST_FILES allows modifying the test scope (files and/or tests) from the
# command line, e.g.:
# GO_TEST_FILES="./cmd/foo -run TestSomethingWorks" make test
GO_TEST_FILES ?= ./...

# Default test timeout; can be changed as with the other GO_TEST* variables
GO_TEST_TIMEOUT ?= 20s


##################
# default target #
##################

##  :

.PHONY: all

## all: lint, test, and build all binaries (default)
all: lint test build


######################
# Makefile logistics #
######################

##  :

.PHONY: help confirm

## help: print this help message
# Note: to add a blank line to the help message, use '##  :' (with two spaces)
help:
	@echo 'Usage:'
	@sed -n 's/^## //p' $(MAKEFILE_LIST) | column -t -s ':' |  sed -e 's/^/  /'

confirm:
	@printf "%s" "Are you sure? [y/N] " && read ans && [ "$$ans" = y ]

# For wildcard targets, instead of .PHONY; see below
# (no-op)
FORCE: ;


######################################
# testing, formatting, linting, etc. #
######################################

##  :

.PHONY: format tidy lint test test-nc testsum testsum-nc

## format: format all code, including adding/removing/reordering imports
# to install:
# go install golang.org/x/tools/cmd/goimports@latest
format:
	@echo ">> Formatting..."
	for source_dir in cmd pkg internal src; do \
		if [ -d "$$source_dir" ]; then \
			goimports -w "$$source_dir"; \
		fi; \
	done

## tidy: tidy and verify module dependencies (go.mod)
tidy:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify

## lint: run static checks on the code
# See https://github.com/golangci/golangci-lint; to install on macOS:
# brew install golangci/tap/golangci-lint
lint:
	@echo ">> Linting..."
	golangci-lint run \
		-E exportloopref,goimports,lll,revive,stylecheck,whitespace \
		--max-issues-per-linter 0 --max-same-issues 0 $(LINT_ARGS)

## test: run all tests
test:
	@echo ">> Testing..."
	$(GO_TEST) -v -race -cover -timeout "$(GO_TEST_TIMEOUT)" $(GO_TEST_FILES)

## test-nc: run all tests (skipping the cache)
# Testing target that skips the cache (more convenient than setting GO_TEST on
# the command line)
test-nc:
	@echo ">> Testing..."
	$(GO_TEST) -count 1 -v -race -cover -timeout "$(GO_TEST_TIMEOUT)" \
		$(GO_TEST_FILES)

## testsum: run all tests using the gotestsum utility
# Testing target that uses gotestsum (more convenient than setting GO_TEST on
# the command line)
# Note the variables for arguments to gotestsum and go test
testsum:
	@echo ">> Testing..."
	gotestsum $(GOTESTSUM_ARGS) -- $(GO_TEST_ARGS) -race -cover \
		-timeout "$(GO_TEST_TIMEOUT)" $(GO_TEST_FILES)

## testsum-nc: run all tests using the gotestsum utility (skipping the cache)
# Testing target that uses gotestsum and skips the cache (more convenient than
# setting GO_TEST on the command line)
# Note the variables for arguments to gotestsum and go test
testsum-nc:
	@echo ">> Testing..."
	gotestsum $(GOTESTSUM_ARGS) -- $(GO_TEST_ARGS) -count 1 -race -cover \
		-timeout "$(GO_TEST_TIMEOUT)" $(GO_TEST_FILES)


########################
# building and running #
########################

##  :

# Note: including e.g. build-% doesn't make build-foo phony, so we use the
# FORCE trick instead
.PHONY: build

bin:
	@echo ">> Making bin directory..."
	mkdir -p bin

output:
	@echo ">> Making output directory..."
	mkdir -p output

## build: build all binaries available under the cmd/ directory in bin/
build: bin output
	@echo ">> Building all binaries..."
	for cmddir in cmd/*; do \
		if [ -f "$$cmddir/main.go" ]; then \
			"$(MAKE)" "build-$${cmddir#cmd/}"; \
		fi; \
	done

## build-SUBDIR: build a binary from the cmd/SUBDIR directory in bin/
build-%: bin output FORCE
	@printf "%s\n" ">> Building binary '$(@:build-%=%)'..."
	bin="$(@:build-%=%)" && \
		GOOS=$$(uname | tr 'A-Z' 'a-z') \
			go build -v -o "bin/$$bin" "./cmd/$$bin/"

## run-SUBDIR: build and run a binary from cmd/SUBDIR (allows ARGS="args")
# This builds one of the available binaries, then runs it.  If the ARGS
# environment variable is set, its contents will be added to the command line
# when running the binary.  For example, 'ARGS="foo bar" make run-bin1' will
# run 'bin1 foo bar' after building bin1.
# ('ARGS="args"' can also be passed to the make command itself, but that might
# not work on all versions of make.)
run-%: output build-% FORCE
	"bin/$(@:run-%=%)" $$ARGS


###############
# cleaning up #
###############

##  :

# Note: including clean-% doesn't make clean-foo phony, so we use the FORCE
# trick instead
.PHONY: clean fullclean

## clean: remove all built binaries from bin/
clean:
	@echo ">> Removing all binaries..."
	for cmddir in cmd/*; do \
		if [ -f "$$cmddir/main.go" ]; then \
			"$(MAKE)" "clean-$${cmddir#cmd/}"; \
		fi; \
	done

## clean-NAME: remove the binary called NAME from bin/
clean-%: FORCE
	@printf "%s\n" ">> Removing binary '$(@:clean-%=%)'..."
	rm -rf "bin/$(@:clean-%=%)"

## fullclean: completely remove all binaries and output files (bin/ and output/)
fullclean:
	@echo ">> Removing all binaries and output files, including directories..."
	rm -rf bin output

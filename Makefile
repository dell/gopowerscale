.PHONY: all
all: build

ifneq (on,$(GO111MODULE))
export GO111MODULE := on
endif

#.PHONY: dep
dep:
	go mod download && go mod verify

.PHONY: build
build:
	git config core.hooksPath hooks
	go build .
#
# Tests-related tasks
#
.PHONY: unit-test
unit-test: build
	go test -json ./... -run ^Test

.PHONY: coverage
coverage: build
	go test -json -covermode=atomic -coverpkg=./... -coverprofile goisilon_coverprofile.out ./... -run ^Test

.PHONY: all
all: go-build

ifneq (on,$(GO111MODULE))
export GO111MODULE := on
endif

#.PHONY: go-dep
go-dep:
	go mod download && go mod verify

.PHONY: go-build
go-build:
	git config core.hooksPath hooks
	go build .
#
# Tests-related tasks
#
.PHONY: go-unittest
go-unittest: go-build
	go test -json ./... -run ^Test

.PHONY: go-coverage
go-coverage: go-build
	go test -json -covermode=atomic -coverpkg=./... -coverprofile goisilon_coverprofile.out ./... -run ^Test


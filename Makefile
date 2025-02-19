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

.PHONY: actions action-help
actions: ## Run all GitHub Action checks that run on a pull request creation
	@echo "Running all GitHub Action checks for pull request events..."
	@act -l | grep -v ^Stage | grep pull_request | grep -v image_security_scan | awk '{print $$2}' | while read WF; do \
		echo "Running workflow: $${WF}"; \
		act pull_request --no-cache-server --platform ubuntu-latest=ghcr.io/catthehacker/ubuntu:act-latest --job "$${WF}"; \
	done

action-help: ## Echo instructions to run one specific workflow locally
	@echo "GitHub Workflows can be run locally with the following command:"
	@echo "act pull_request --no-cache-server --platform ubuntu-latest=ghcr.io/catthehacker/ubuntu:act-latest --job <jobid>"
	@echo ""
	@echo "Where '<jobid>' is a Job ID returned by the command:"
	@echo "act -l"
	@echo ""
	@echo "NOTE: if act is not installed, it can be downloaded from https://github.com/nektos/act"

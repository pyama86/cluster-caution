VERSION:=$(shell (grep -e 'const Version' cli.go | sed 's/const Version string = //g'))
TEST ?= $(shell go list ./... | grep -v -e vendor -e keys -e tmp)
INFO_COLOR=\033[1;34m
RESET=\033[0m
BOLD=\033[1m
GO ?= GO111MODULE=on go

release:
	git tag -a $(VERSION) -m "bump to $(VERSION)" || true
	goreleaser --rm-dist
test: ## Run test
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Testing$(RESET)"
	$(GO) test -v $(TEST) -timeout=30s -parallel=4
	$(GO) test -race $(TEST)

run:
	$(GO) run main.go cli.go $(OPT)

.PHONY: release

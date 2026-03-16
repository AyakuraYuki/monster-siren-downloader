name    = msr-downloader
dst_dir = build
target  = $(GOOS)-$(GOARCH)$(GOARM)
build   = $(dst_dir)/$(target)

VERSION ?= dev

ifeq ($(GOOS),windows)
    ext       = .exe
    chmod_cmd =
    pack_cmd  = zip -9 -r $(name)-$(target)-$(VERSION).zip $(target)
else
    ext       =
    chmod_cmd = chmod +x $(dst_dir)/$(name)
    pack_cmd  = tar czpvf $(name)-$(target)-$(VERSION).tar.gz $(target)
endif

define check_env
	@ if [ "$(GOOS)" = "" ]; then echo " <- Env variable GOOS not set"; exit 1; fi
	@ if [ "$(GOARCH)" = "" ]; then echo " <- Env variable GOARCH not set"; exit 1; fi
endef

.DEFAULT_GOAL := help
.PHONY: help
help:
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-16s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)


.PHONY: build
build: clean test ## Compile the package targeted to current platform; the package will be cleaned and tested before compilation.
	@CGO_ENABLED=0 go build -ldflags "-s -w -X main.version=$(VERSION)" -o $(build)/$(name)$(ext)
	@echo " -> $(build)/$(name)$(ext)"

.PHONY: dev
dev: clean ## Compile the development package targeted to current platform.
	@CGO_ENABLED=0 go build -o $(build)/$(name)$(ext)
	@echo " -> $(build)/$(name)$(ext)"

.PHONY: echo
echo: ## Print the parameters in Makefile.
	@echo "     name: $(name)"
	@echo "  dst_dir: $(dst_dir)"
	@echo "   target: $(target)"
	@echo "    build: $(build)"
	@echo "  VERSION: $(VERSION)"
	@echo "      ext: $(ext)"
	@echo "chmod_cmd: $(chmod_cmd)"
	@echo " pack_cmd: $(pack_cmd)"

.PHONY: test
test: ## Test the go package if it has the test cases.
	@go test -v -bench=. ./...

.PHONY: clean
clean: ## Remove build caches, temp files, and the previous build outputs.
	@go clean
	@rm -vrf $(dst_dir)
	@echo " <- done"

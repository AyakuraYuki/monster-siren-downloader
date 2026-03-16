name    = msr-downloader
dst_dir = build
target  = $(GOOS)-$(GOARCH)$(GOARM)
build   = $(dst_dir)/$(target)

VERSION ?= dev

ifeq ($(GOOS),windows)
    ext       = .exe
    chmod_cmd =
    pack_cmd  = zip -9 -r $(name)-$(target)-$(VERSION).zip "$(target)"
else
    ext       =
    chmod_cmd = chmod +x $(build)/$(name)
    pack_cmd  = tar czpvf $(name)-$(target)-$(VERSION).tar.gz "$(target)"
endif

define check_env
	@ if [ "$(GOOS)" = "" ]; then echo " <- Env variable GOOS not set"; exit 1; fi
	@ if [ "$(GOARCH)" = "" ]; then echo " <- Env variable GOARCH not set"; exit 1; fi
endef

.PHONY: help build dev release release-arm test echo clean

.DEFAULT_GOAL := help
help:
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-16s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)


build: clean test ## Compile the package targeted to current platform; the package will be cleaned and tested before compilation.
	@CGO_ENABLED=0 go build -ldflags "-s -w -X main.version=$(VERSION)" -o $(build)/$(name)$(ext)
	@echo " -> $(build)/$(name)$(ext)"

dev: clean ## Compile the development package targeted to current platform.
	@CGO_ENABLED=0 go build -o $(build)/$(name)$(ext)
	@echo " -> $(build)/$(name)$(ext)"

release: ## Build the release package.
	@$(call check_env)
	@mkdir -p $(build)
	@cp LICENSE $(build)/
	@cp README.md $(build)/
	@CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags "-s -w -X main.version=$(VERSION)" -o $(build)/$(name)$(ext)
	@$(chmod_cmd)
	@cd $(dst_dir) ; $(pack_cmd)

release-arm: ## Build the release package with GOARM.
	@$(call check_env)
	@mkdir -p $(build)
	@cp LICENSE $(build)/
	@cp README.md $(build)/
	@CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) GOARM=$(GOARM) go build -ldflags "-s -w -X main.version=$(VERSION)" -o $(build)/$(name)$(ext)
	@$(chmod_cmd)
	@cd $(dst_dir) ; $(pack_cmd)

test: ## Test the go package if it has the test cases.
	@go test -v -bench=. ./...

echo: ## Print the parameters in Makefile.
	@echo "     name: $(name)"
	@echo "  dst_dir: $(dst_dir)"
	@echo "   target: $(target)"
	@echo "    build: $(build)"
	@echo "  VERSION: $(VERSION)"
	@echo "      ext: $(ext)"
	@echo "chmod_cmd: $(chmod_cmd)"
	@echo " pack_cmd: $(pack_cmd)"

clean: ## Remove build caches, temp files, and the previous build outputs.
	@go clean
	@rm -vrf $(dst_dir)
	@echo " <- done"

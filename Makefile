NAME=msr-downloader
OUT_DIR=build  # root directory for build outputs
TARGET=$(GOOS)-$(GOARCH)$(GOARM)  # target platform identifier
BIN_DIR=$(OUT_DIR)/$(TARGET)  # Platform-specific binary directory
VERSION?=dev

ifeq ($(GOOS),windows)
  EXT=.exe
  PACK_CMD=zip -9 -r $(NAME)-$(TARGET)-$(VERSION).zip $(TARGET)
else
  EXT=
  PACK_CMD=tar czpvf $(NAME)-$(TARGET)-$(VERSION).tar.gz $(TARGET)
endif

define check_env
	@ if [ "$(GOOS)" = "" ]; then echo " <- Env variable GOOS not set"; exit 1; fi
	@ if [ "$(GOARCH)" = "" ]; then echo " <- Env variable GOARCH not set"; exit 1; fi
endef

# Self-Documented Makefile see https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.DEFAULT_GOAL := help
.PHONY: help
help:
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-16s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)


.PHONY: build
build: clean test ## Compile the package targeted to current platform; the package will be cleaned and tested before compilation.
	@go build -o $(BIN_DIR)/$(NAME)$(EXT)

.PHONY: release
release: ## release builds releasable artifacts. You need to specify two environment variables, GOOS and GOARCH, to set the target platform for the binaries.
	@$(call check_env)
	@mkdir -p $(BIN_DIR)
	@cp LICENSE $(BIN_DIR)/
	@cp README.md $(BIN_DIR)/
	@CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags "-s -w -X main.version=$(VERSION)" -o $(BIN_DIR)/$(NAME)$(EXT)
	@cd $(OUT_DIR) ; $(PACK_CMD)

.PHONY: test
test: ## Test the go package if it has the test cases.
	@go test -race -v -bench=. ./...

.PHONY: clean
clean: ## Remove build caches, temp files, and the previous build outputs.
	@go clean
	@rm -vrf $(OUT_DIR)
	@echo " <- done"

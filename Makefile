APP_NAME := jcscli
PROCS := 1
VERSION ?= dev
DIST_DIR := dist
PLATFORMS := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64

all: clean build package changelog

.PHONY: test benchmark

test:
	GOMAXPROCS=$(PROCS) go test -v .

benchmark:
	GOMAXPROCS=$(PROCS) go test -gcflags '$(FLAGS)' -benchmem -cpuprofile cpu.out -memprofile mem.out $(ARGS) .

clean:
	rm -rf $(DIST_DIR)
	mkdir -p $(DIST_DIR)

build:
	@for platform in $(PLATFORMS); do \
		GOOS=$${platform%/*}; \
		GOARCH=$${platform#*/}; \
		OUTPUT=$(DIST_DIR)/$(APP_NAME)-$$GOOS-$$GOARCH; \
		if [ $$GOOS = windows ]; then OUTPUT=$$OUTPUT.exe; fi; \
		echo "→ Building $$GOOS/$$GOARCH"; \
		GOOS=$$GOOS GOARCH=$$GOARCH go build -ldflags "-X main.version=$(VERSION)" -o $$OUTPUT ./cmd/jcscli; \
	done

package:
	@for platform in $(PLATFORMS); do \
		GOOS=$${platform%/*}; \
      GOARCH=$${platform#*/}; \
      OUTPUT=$(DIST_DIR)/$(APP_NAME)-$$GOOS-$$GOARCH; \
      if [ $$GOOS = windows ]; then \
          zip -j $${OUTPUT%.exe}.zip $$OUTPUT.exe; \
          rm $$OUTPUT.exe; \
      else \
          tar -czf $$OUTPUT.tar.gz -C $(DIST_DIR) $$(basename $$OUTPUT); \
          rm $$OUTPUT; \
      fi; \
	done

changelog:
	@DATE=$$(date +"%Y-%m-%d"); \
	echo "## [$(VERSION)] - $$DATE\n### Added\n- Release artifacts built for Linux, macOS, Windows.\n- Canonical JSON CLI tool updates.\n" >> CHANGELOG.md

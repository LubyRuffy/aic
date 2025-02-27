.PHONY: all build test coverage clean release

GO=go
PKG=./...
BINARY_NAME=aic
COVER_PROFILE=coverage.out

# 默认目标：编译
all: build

# 编译项目
build:
	$(GO) build -o $(BINARY_NAME)

# 运行测试并生成覆盖率报告
test:
	$(GO) test $(PKG) -v -race

coverage:
	$(GO) test $(PKG) -v -race -coverprofile=$(COVER_PROFILE)
	$(GO) tool cover -html=$(COVER_PROFILE)

# 清理编译和测试产物
clean:
	rm -f $(BINARY_NAME)
	rm -f $(COVER_PROFILE)
	rm -rf dist/

# 发布：运行测试、清理、使用goreleaser发布
release: test clean
	@if [ -z "$$(git tag)" ]; then \
		read -p "No git tags found. Do you want to create a new tag? (y/n): " answer; \
		if [ "$$answer" = "y" ]; then \
			read -p "Enter new tag (e.g. v1.0.0): " tag; \
			git tag "$$tag"; \
			goreleaser release --clean; \
		else \
			echo "Using snapshot mode..."; \
			goreleaser release --snapshot --clean; \
		fi \
	else \
		goreleaser release --clean; \
	fi
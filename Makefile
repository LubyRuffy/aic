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
	@if [ -n "$$(git status --porcelain)" ]; then \
		echo "警告：检测到以下未提交的更改："; \
		git status --short; \
		read -p "是否要提交这些更改？(y/n): " commit_changes; \
		if [ "$$commit_changes" = "y" ]; then \
			read -p "请输入提交信息: " commit_msg; \
			git add .; \
			git commit -m "$$commit_msg"; \
			git push; \
		else \
			echo "发布已取消"; \
			exit 1; \
		fi; \
	fi
	@if [ -z "$$(git tag)" ]; then \
		read -p "未找到git标签。是否创建新标签？(y/n): " answer; \
		if [ "$$answer" = "y" ]; then \
			read -p "请输入新标签 (例如 v1.0.0): " tag; \
			git tag "$$tag"; \
			git push origin "$$tag"; \
			goreleaser release --clean; \
		else \
			echo "使用快照模式..."; \
			goreleaser release --snapshot --clean; \
		fi \
	else \
		goreleaser release --clean; \
	fi
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
	@if [ -z "$$GITHUB_TOKEN" ]; then \
		echo "错误：未设置GITHUB_TOKEN环境变量"; \
		echo "请设置GITHUB_TOKEN环境变量以访问GitHub API"; \
		echo "示例：export GITHUB_TOKEN='your-token'"; \
		exit 1; \
	fi
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
		if git describe --tags --abbrev=0 >/dev/null 2>&1; then \
			latest_tag=$$(git describe --tags --abbrev=0); \
			tag_commit=$$(git rev-list -n 1 $$latest_tag); \
			head_commit=$$(git rev-parse HEAD); \
			if [ "$$tag_commit" != "$$head_commit" ]; then \
				echo "警告：最新标签 $$latest_tag 不是在当前commit上创建的"; \
				read -p "是否创建新标签？(y/n/r - y:创建新标签, n:取消, r:重建当前标签): " tag_action; \
				if [ "$$tag_action" = "y" ]; then \
					read -p "请输入新标签 (例如 v1.0.0): " new_tag; \
					git tag "$$new_tag"; \
					git push origin "$$new_tag"; \
					echo "已创建新标签 $$new_tag"; \
					goreleaser release --clean; \
				elif [ "$$tag_action" = "r" ]; then \
					echo "正在删除本地和远程标签..."; \
					git tag -d $$latest_tag; \
					git push origin :refs/tags/$$latest_tag; \
					echo "正在删除相关的draft releases..."; \
					curl -s -H "Authorization: token $$GITHUB_TOKEN" \
						"https://api.github.com/repos/LubyRuffy/aic/releases" \
						| jq -r ".[] | select(.tag_name == \"$$latest_tag\") | .id" \
						| while read -r release_id; do \
							curl -s -X DELETE -H "Authorization: token $$GITHUB_TOKEN" \
								"https://api.github.com/repos/LubyRuffy/aic/releases/$$release_id"; \
						done;
					echo "正在重新创建标签..."; \
					git tag $$latest_tag; \
					git push origin $$latest_tag; \
					echo "已重新创建标签 $$latest_tag"; \
					goreleaser release --clean; \
				else \
					echo "发布已取消"; \
					exit 1; \
				fi \
			else \
				goreleaser release --clean; \
			fi \
		else \
			echo "当前仓库没有任何标签"; \
			read -p "请输入新标签 (例如 v1.0.0): " new_tag; \
			git tag "$$new_tag"; \
			git push origin "$$new_tag"; \
			echo "已创建新标签 $$new_tag"; \
			goreleaser release --clean; \
		fi \
	fi
# =============================================================================
# mDNS Reflector Go - Makefile
# =============================================================================
#
# 🚀 快速开始:
#   make help          # 显示所有可用命令
#   make build         # 构建当前平台二进制
#   make dev           # 开发模式构建 (带调试信息)
#   make test          # 运行测试
#   make release       # 创建发布包
#

# =============================================================================
# 配置变量
# =============================================================================

# 项目信息
BINARY_NAME    := mdns-reflector-go
VERSION        ?= dev
BUILD_TIME     := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT     := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DIR      := build
DIST_DIR       := dist

# Go 构建配置
GOOS           := $(shell go env GOOS)
GOARCH         := $(shell go env GOARCH)
CGO_ENABLED    := 0

# 构建标志 - 优化二进制大小和性能
LDFLAGS        := -ldflags "\
	-s -w \
	-X main.version=$(VERSION) \
	-X main.commit=$(GIT_COMMIT) \
	-X main.buildTime=$(BUILD_TIME) \
	-extldflags '-static'"

# 导出环境变量
export CGO_ENABLED

# =============================================================================
# 构建目标平台配置
# =============================================================================

# macOS (Intel + Apple Silicon)
PLATFORMS_MACOS := darwin-amd64 darwin-arm64

# Linux (x86_64 + ARM64)
PLATFORMS_LINUX := linux-amd64 linux-arm64

# Windows (x86_64)
PLATFORMS_WINDOWS := windows-amd64

# 所有平台
ALL_PLATFORMS := $(PLATFORMS_MACOS) $(PLATFORMS_LINUX) $(PLATFORMS_WINDOWS)

# =============================================================================
# 主目标
# =============================================================================

.PHONY: all help
all: build

# 显示帮助信息
help:
	@echo "🚀 mDNS Reflector Go - 构建工具"
	@echo ""
	@echo "📦 主要命令:"
	@echo "  build        构建当前平台二进制 (生产优化)"
	@echo "  dev          开发模式构建 (带调试信息)"
	@echo "  build-all    构建所有平台 (macOS/Linux/Windows)"
	@echo "  release      创建发布包和校验和"
	@echo ""
	@echo "🧪 测试和质量:"
	@echo "  test         运行单元测试"
	@echo "  coverage     生成测试覆盖率报告"
	@echo "  bench        运行性能基准测试"
	@echo "  lint         代码质量检查"
	@echo ""
	@echo "🛠️  开发工具:"
	@echo "  clean        清理构建产物"
	@echo "  deps         更新依赖"
	@echo "  install      安装到系统 (/usr/local/bin)"
	@echo "  docker       构建 Docker 镜像"
	@echo ""
	@echo "🔒 安全检查:"
	@echo "  security     运行安全扫描"
	@echo "  audit        依赖审计"
	@echo ""
	@echo "📊 当前配置:"
	@echo "  版本: $(VERSION)"
	@echo "  平台: $(GOOS)/$(GOARCH)"
	@echo "  提交: $(GIT_COMMIT)"

# =============================================================================
# 构建目标
# =============================================================================

.PHONY: build dev build-all release

# 生产构建 - 优化大小和性能
build: clean
	@echo "🏗️  构建生产版本 ($(GOOS)/$(GOARCH))..."
	@mkdir -p $(BUILD_DIR)
	@echo "   📦 二进制: $(BINARY_NAME)"
	@echo "   🏷️  版本: $(VERSION)"
	@echo "   🔗 提交: $(GIT_COMMIT)"
	@go build $(LDFLAGS) -trimpath -o $(BUILD_DIR)/$(BINARY_NAME) .
	@du -h $(BUILD_DIR)/$(BINARY_NAME)
	@echo "✅ 构建完成: $(BUILD_DIR)/$(BINARY_NAME)"

# 开发构建 - 保留调试信息
dev: clean
	@echo "🔧 开发模式构建 ($(GOOS)/$(GOARCH))..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME)-dev .
	@echo "✅ 开发版本: $(BUILD_DIR)/$(BINARY_NAME)-dev"

# 构建所有平台
build-all: clean
	@echo "🌍 构建所有平台..."
	@mkdir -p $(BUILD_DIR)
	@$(MAKE) build-platforms PLATFORMS="$(PLATFORMS_MACOS)" OS_NAME="macOS"
	@$(MAKE) build-platforms PLATFORMS="$(PLATFORMS_LINUX)" OS_NAME="Linux"
	@$(MAKE) build-platforms PLATFORMS="$(PLATFORMS_WINDOWS)" OS_NAME="Windows"
	@echo "📊 构建统计:"
	@ls -lh $(BUILD_DIR)/* | grep -E '\.(exe|bin|)$$' | wc -l | xargs echo "   总文件数: "
	@du -sh $(BUILD_DIR) | cut -f1 | xargs echo "   总大小: "

# 平台构建辅助函数
build-platforms:
	@echo "   🖥️  $(OS_NAME): $(PLATFORMS)"
	@for platform in $(PLATFORMS); do \
		IFS='-' read -r os arch <<< "$$platform"; \
		echo "     🔨 $$os/$$arch..."; \
		ext=""; \
		if [ "$$os" = "windows" ]; then ext=".exe"; fi; \
		GOOS=$$os GOARCH=$$arch go build $(LDFLAGS) -trimpath -o $(BUILD_DIR)/$(BINARY_NAME)-$$platform$$ext .; \
	done

# 创建发布包
release: build-all
	@echo "📦 创建发布包..."
	@mkdir -p $(DIST_DIR)
	@$(MAKE) create-archives
	@$(MAKE) generate-checksums
	@echo "📋 发布文件:"
	@ls -lh $(DIST_DIR)/*
	@echo "📊 校验和:"
	@cat $(DIST_DIR)/checksums.txt

# 创建压缩包
create-archives:
	@echo "   📦 打包文件..."
	@cd $(BUILD_DIR) && \
	for file in $(BINARY_NAME)-*; do \
		base_name=$$(basename $$file | sed 's/\.[^.]*$$//'); \
		if [[ $$file == *.exe ]]; then \
			zip -q ../$(DIST_DIR)/$${base_name}-$(VERSION).zip $$file; \
			echo "     📚 $${base_name}-$(VERSION).zip"; \
		else \
			tar -czf ../$(DIST_DIR)/$${base_name}-$(VERSION).tar.gz $$file; \
			echo "     📦 $${base_name}-$(VERSION).tar.gz"; \
		fi; \
	done

# 生成校验和
generate-checksums:
	@echo "   🔐 生成校验和..."
	@cd $(DIST_DIR) && \
	find . -type f \( -name "*.tar.gz" -o -name "*.zip" \) -exec shasum -a 256 {} \; > checksums.txt && \
	echo "   ✅ checksums.txt 已生成"

# =============================================================================
# 测试和质量保证
# =============================================================================

.PHONY: test coverage bench lint security audit

test:
	@echo "🧪 运行测试..."
	@go test -v -race -timeout 30s ./...

coverage:
	@echo "📊 生成覆盖率报告..."
	@go test -race -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "📈 覆盖率报告: coverage.html"

bench:
	@echo "⚡ 运行性能基准测试..."
	@go test -bench=. -benchmem ./...

lint: fmt vet
	@echo "🔍 代码质量检查..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "   ⚠️  golangci-lint 未安装，跳过高级检查"; \
	fi

security:
	@echo "🔒 安全扫描..."
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "   ⚠️  gosec 未安装，跳过安全扫描"; \
	fi

audit:
	@echo "📋 依赖审计..."
	@if command -v govulncheck >/dev/null 2>&1; then \
		govulncheck ./...; \
	else \
		echo "   ⚠️  govulncheck 未安装，使用 go mod verify"; \
		go mod verify; \
	fi

# =============================================================================
# 开发工具
# =============================================================================

.PHONY: clean deps install docker fmt vet

clean:
	@echo "🧹 清理构建产物..."
	@rm -rf $(BUILD_DIR) $(DIST_DIR)
	@rm -f coverage.out coverage.html
	@echo "✅ 清理完成"

deps:
	@echo "📦 更新依赖..."
	@go mod download
	@go mod tidy
	@echo "✅ 依赖已更新"

install: build
	@echo "📥 安装到系统 (/usr/local/bin)..."
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@echo "✅ 已安装: $(BINARY_NAME)"

docker:
	@echo "🐳 构建 Docker 镜像..."
	@docker build -t $(BINARY_NAME):$(VERSION) .
	@echo "✅ Docker 镜像: $(BINARY_NAME):$(VERSION)"

fmt:
	@echo "🎨 格式化代码..."
	@go fmt ./...
	@echo "✅ 代码已格式化"

vet:
	@echo "🔬 静态分析..."
	@go vet ./...
	@echo "✅ 静态分析完成"

# =============================================================================
# 特殊目标
# =============================================================================

# 开发模式别名
.PHONY: run watch

run: dev
	@echo "🎮 运行开发版本..."
	@./$(BUILD_DIR)/$(BINARY_NAME)-dev --help

watch:
	@echo "👀 监听文件变化..."
	@if command -v air >/dev/null 2>&1; then \
		air; \
	else \
		echo "   ⚠️  air 未安装，请运行: go install github.com/cosmtrek/air@latest"; \
	fi

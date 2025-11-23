#!/bin/bash

# =============================================================================
# mDNS Reflector Go - Release 文件上传脚本
# =============================================================================
#
# 🚀 功能:
#   - 自动上传所有发布包到 GitHub Release
#   - 验证文件完整性
#   - 生成上传报告
#
# 📋 使用方法:
#   1. 确保已创建 GitHub Release v1.1.0
#   2. 设置 GitHub Token (可选): export GITHUB_TOKEN=your_token
#   3. 运行: ./upload-release.sh
#

set -e

# 配置
REPO="FangTianwd/mdns_reflector_go"
TAG="v1.1.0"
DIST_DIR="dist"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

log_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

log_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

log_error() {
    echo -e "${RED}❌ $1${NC}"
}

# 检查依赖
check_dependencies() {
    log_info "检查依赖..."

    if ! command -v curl >/dev/null 2>&1; then
        log_error "需要安装 curl"
        exit 1
    fi

    if ! command -v jq >/dev/null 2>&1; then
        log_warning "未安装 jq，将使用基本模式"
        USE_JQ=false
    else
        USE_JQ=true
    fi

    log_success "依赖检查完成"
}

# 验证文件
verify_files() {
    log_info "验证发布文件..."

    if [ ! -d "$DIST_DIR" ]; then
        log_error "dist 目录不存在，请先运行: make release VERSION=v1.0.2"
        exit 1
    fi

    local expected_files=(
        "checksums.txt"
        "mdns-reflector-go-darwin-amd64-v1.0.2.tar.gz"
        "mdns-reflector-go-darwin-arm64-v1.0.2.tar.gz"
        "mdns-reflector-go-linux-amd64-v1.0.2.tar.gz"
        "mdns-reflector-go-linux-arm64-v1.0.2.tar.gz"
        "mdns-reflector-go-windows-amd64-v1.0.2.zip"
    )

    for file in "${expected_files[@]}"; do
        if [ ! -f "$DIST_DIR/$file" ]; then
            log_error "缺少文件: $file"
            exit 1
        fi
    done

    log_success "所有文件验证通过"
}

# 验证校验和
verify_checksums() {
    log_info "验证文件校验和..."

    cd "$DIST_DIR"

    if [ -f "checksums.txt" ]; then
        if shasum -a 256 --check checksums.txt >/dev/null 2>&1; then
            log_success "校验和验证通过"
        else
            log_error "校验和验证失败"
            exit 1
        fi
    else
        log_warning "未找到 checksums.txt 文件"
    fi

    cd ..
}

# 检查 Release 存在
check_release() {
    log_info "检查 GitHub Release 状态..."

    local response
    response=$(curl -s -w "%{http_code}" -o /tmp/release_check.json \
        "https://api.github.com/repos/$REPO/releases/tags/$TAG")

    local status_code=${response: -3}

    if [ "$status_code" = "200" ]; then
        log_success "Release v1.0.2 已存在"
        return 0
    elif [ "$status_code" = "404" ]; then
        log_error "Release v1.0.2 不存在，请先在 GitHub 上创建"
        log_info "访问: https://github.com/$REPO/releases/new"
        exit 1
    else
        log_error "检查 Release 失败 (HTTP $status_code)"
        exit 1
    fi
}

# 获取 Release ID
get_release_id() {
    log_info "获取 Release ID..."

    if [ "$USE_JQ" = true ]; then
        RELEASE_ID=$(jq -r '.id' /tmp/release_check.json)
    else
        # 基本模式：提取 id
        RELEASE_ID=$(grep '"id"' /tmp/release_check.json | head -1 | sed 's/.*: *//' | tr -d ',')
    fi

    if [ -z "$RELEASE_ID" ] || [ "$RELEASE_ID" = "null" ]; then
        log_error "无法获取 Release ID"
        exit 1
    fi

    log_success "Release ID: $RELEASE_ID"
}

# 上传文件 (使用 GitHub CLI 如果可用)
upload_files() {
    log_info "开始上传文件..."

    # 检查是否安装了 GitHub CLI
    if command -v gh >/dev/null 2>&1; then
        log_info "使用 GitHub CLI 上传..."
        upload_with_gh
    else
        log_warning "未安装 GitHub CLI，请手动上传文件"
        manual_upload_instructions
    fi
}

# 使用 GitHub CLI 上传
upload_with_gh() {
    cd "$DIST_DIR"

    local files=(
        "checksums.txt"
        "mdns-reflector-go-darwin-amd64-v1.0.2.tar.gz"
        "mdns-reflector-go-darwin-arm64-v1.0.2.tar.gz"
        "mdns-reflector-go-linux-amd64-v1.0.2.tar.gz"
        "mdns-reflector-go-linux-arm64-v1.0.2.tar.gz"
        "mdns-reflector-go-windows-amd64-v1.0.2.zip"
    )

    for file in "${files[@]}"; do
        log_info "上传: $file"
        if gh release upload "$TAG" "$file" --clobber; then
            log_success "上传成功: $file"
        else
            log_error "上传失败: $file"
            return 1
        fi
    done

    cd ..
    log_success "所有文件上传完成"
}

# 手动上传指导
manual_upload_instructions() {
    log_warning "请手动上传以下文件到 GitHub Release:"

    echo ""
    echo "📋 步骤:"
    echo "1. 访问: https://github.com/$REPO/releases/tag/$TAG"
    echo "2. 点击 'Edit' 按钮"
    echo "3. 在文件上传区域，拖拽或选择以下文件:"
    echo ""

    ls -1 "$DIST_DIR"/ | while read -r file; do
        echo "   📎 $file"
    done

    echo ""
    echo "4. 点击 'Update release'"
    echo ""
    echo "🔐 文件校验和 (验证完整性):"
    echo "----------------------------------------"
    cat "$DIST_DIR/checksums.txt"
    echo "----------------------------------------"
}

# 主函数
main() {
    echo "🚀 mDNS Reflector Go - Release 文件上传工具"
    echo "=========================================="
    echo ""

    check_dependencies
    verify_files
    verify_checksums
    check_release
    get_release_id
    upload_files

    echo ""
    log_success "🎉 上传流程完成！"
    echo ""
    echo "📋 下一步:"
    echo "1. 验证 GitHub Release 页面已显示所有文件"
    echo "2. 测试安装: brew install fangtianwd/tap/mdns-reflector-go"
    echo ""
    echo "🔗 Release 页面: https://github.com/$REPO/releases/tag/$TAG"
}

# 运行主函数
main "$@"


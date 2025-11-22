#!/bin/bash

# 完整的 Homebrew Tap 设置脚本

set -e

echo "🎯 完整的 Homebrew Tap 设置流程"
echo "================================="
echo ""

# 配置
USERNAME="FangTianwd"
PROJECT_NAME="mdns-reflector-go"
VERSION="v1.0.0"
TAP_REPO="https://github.com/${USERNAME}/homebrew-tap.git"

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

# 步骤 1: 检查前置条件
check_prerequisites() {
    log_info "步骤 1: 检查前置条件"

    # 检查必要文件
    local required_files=("mdns-reflector-go.rb" "Makefile" "scripts/update-formula.sh")
    for file in "${required_files[@]}"; do
        if [ ! -f "$file" ]; then
            log_error "缺少文件: $file"
            exit 1
        fi
    done

    # 检查 git
    if ! command -v git &> /dev/null; then
        log_error "需要安装 git"
        exit 1
    fi

    # 检查是否在 git 仓库中
    if ! git rev-parse --git-dir > /dev/null 2>&1; then
        log_error "当前目录不是 git 仓库"
        exit 1
    fi

    log_success "前置条件检查完成"
    echo ""
}

# 步骤 2: 创建版本标签
create_version_tag() {
    log_info "步骤 2: 创建版本标签 $VERSION"

    # 检查标签是否已存在
    if git tag -l | grep -q "^${VERSION}$"; then
        log_warning "标签 $VERSION 已存在，跳过创建"
    else
        git tag "$VERSION"
        git push origin "$VERSION"
        log_success "创建并推送标签 $VERSION"
    fi

    echo ""
}

# 步骤 3: 克隆或更新 tap 仓库
setup_tap_repo() {
    log_info "步骤 3: 设置 Tap 仓库"

    if [ -d "homebrew-tap" ]; then
        log_info "更新现有的 homebrew-tap 目录"
        cd homebrew-tap
        git pull origin main
        cd ..
    else
        log_info "克隆 homebrew-tap 仓库"
        git clone "$TAP_REPO" homebrew-tap
    fi

    log_success "Tap 仓库准备完成"
    echo ""
}

# 步骤 4: 复制并配置 formula
setup_formula() {
    log_info "步骤 4: 配置 Formula"

    cd homebrew-tap

    # 复制 formula 文件
    cp ../mdns-reflector-go.rb Formula/

    # 更新 formula 信息
    local formula_file="Formula/mdns-reflector-go.rb"

    # 更新 homepage
    sed -i.bak "s|homepage \".*\"|homepage \"https://github.com/${USERNAME}/${PROJECT_NAME}\"|" "$formula_file"

    # 更新 url 中的用户名
    sed -i.bak "s|url \"https://github.com/[^/]*/|url \"https://github.com/${USERNAME}/|" "$formula_file"

    # 清理备份文件
    rm -f "${formula_file}.bak"

    log_success "Formula 配置完成"
    echo ""
}

# 步骤 5: 测试 formula
test_formula() {
    log_info "步骤 5: 测试 Formula"

    if command -v brew &> /dev/null; then
        local formula_file="Formula/mdns-reflector-go.rb"

        log_info "运行 brew audit..."
        if brew audit --strict "$formula_file"; then
            log_success "Formula 语法检查通过"
        else
            log_warning "Formula 语法检查失败，请手动修复"
        fi

        log_info "运行 brew style..."
        if brew style "$formula_file"; then
            log_success "Formula 样式检查通过"
        else
            log_warning "Formula 样式检查失败，请手动修复"
        fi
    else
        log_warning "未找到 brew 命令，跳过测试"
    fi

    echo ""
}

# 步骤 6: 提交更改
commit_changes() {
    log_info "步骤 6: 提交更改到 Tap 仓库"

    # 检查是否有更改
    if git diff --quiet && git diff --staged --quiet; then
        log_info "没有需要提交的更改"
    else
        git add .
        git commit -m "Add ${PROJECT_NAME} ${VERSION}"
        git push origin main
        log_success "更改已提交并推送"
    fi

    cd ..
    echo ""
}

# 步骤 7: 等待 GitHub Actions 并更新 SHA256
wait_and_update_sha256() {
    log_info "步骤 7: 等待发布并更新 SHA256"

    echo "请完成以下手动步骤:"
    echo ""
    echo "1. 🌐 访问 https://github.com/${USERNAME}/${PROJECT_NAME}/releases"
    echo "2. 📦 等待 GitHub Actions 完成构建 (大约 2-3 分钟)"
    echo "3. 🔍 确认 $VERSION release 已创建"
    echo ""
    echo "完成上述步骤后，运行:"
    echo "./scripts/update-formula.sh $VERSION"
    echo ""

    read -p "按回车键继续..."
    echo ""
}

# 步骤 8: 最终测试
final_test() {
    log_info "步骤 8: 最终测试"

    echo "运行以下命令测试安装:"
    echo ""
    echo "# 添加你的 tap"
    echo "brew tap ${USERNAME}/homebrew-tap"
    echo ""
    echo "# 安装软件"
    echo "brew install ${PROJECT_NAME}"
    echo ""
    echo "# 启动服务"
    echo "brew services start ${PROJECT_NAME}"
    echo ""
    echo "# 验证"
    echo "brew services list | grep ${PROJECT_NAME}"
    echo ""

    log_success "设置完成！🎉"
}

# 主函数
main() {
    echo "开始设置你的 Homebrew Tap..."
    echo ""

    check_prerequisites
    create_version_tag
    setup_tap_repo
    setup_formula
    test_formula
    commit_changes
    wait_and_update_sha256
    final_test

    echo ""
    echo "📚 相关文档:"
    echo "- 快速开始: QUICKSTART.md"
    echo "- 故障排除: TROUBLESHOOTING.md"
    echo "- 详细指南: HOMEBREW_TAP_SETUP.md"
}

# 运行主函数
main "$@"

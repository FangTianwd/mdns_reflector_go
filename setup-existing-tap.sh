#!/bin/bash

# 设置现有的 Homebrew Tap 仓库的脚本

set -e

echo "🔄 设置现有的 Homebrew Tap 仓库"
echo "=================================="
echo ""

# 配置变量
TAP_REPO="https://github.com/FangTianwd/homebrew-tap.git"
TAP_NAME="FangTianwd/homebrew-tap"
PROJECT_DIR="$(pwd)"
FORMULA_FILE="mdns-reflector-go.rb"

# 步骤 1: 克隆 tap 仓库
clone_tap_repo() {
    echo "📥 步骤 1: 克隆 tap 仓库..."

    if [ -d "homebrew-tap" ]; then
        echo "⚠️  本地已存在 homebrew-tap 目录，正在更新..."
        cd homebrew-tap
        git pull origin main
        cd ..
    else
        git clone "$TAP_REPO" homebrew-tap
    fi

    echo "✅ Tap 仓库克隆完成!"
    echo ""
}

# 步骤 2: 检查仓库结构
check_repo_structure() {
    echo "🔍 步骤 2: 检查仓库结构..."

    cd homebrew-tap

    echo "📂 当前目录结构:"
    find . -type f -name "*.rb" | head -10

    # 检查是否有 Formula 目录
    if [ -d "Formula" ]; then
        echo "✅ 发现 Formula 目录"
        FORMULA_DIR="Formula"
    else
        echo "ℹ️  没有 Formula 目录，将直接在根目录放置 formula"
        FORMULA_DIR="."
    fi

    cd ..
    echo ""
}

# 步骤 3: 复制并更新 formula
setup_formula() {
    echo "📋 步骤 3: 设置 formula..."

    cd homebrew-tap

    # 复制 formula 文件
    cp "../$FORMULA_FILE" "$FORMULA_DIR/"

    # 更新 formula 中的信息
    FORMULA_PATH="$FORMULA_DIR/$FORMULA_FILE"
    sed -i.bak "s|homepage \".*\"|homepage \"https://github.com/FangTianwd/mdns-reflector-go\"|" "$FORMULA_PATH"
    sed -i.bak "s|url \".*\"|url \"https://github.com/FangTianwd/mdns-reflector-go/archive/refs/tags/v#{version}.tar.gz\"|" "$FORMULA_PATH"
    sed -i.bak "s|sha256 \".*\"|sha256 \"CHANGE_THIS_WITH_ACTUAL_SHA256\"|" "$FORMULA_PATH"

    # 删除备份文件
    rm -f "$FORMULA_DIR/$FORMULA_FILE.bak"

    echo "✅ Formula 已复制并更新"
    echo "📄 Formula 位置: $FORMULA_PATH"
    echo ""

    cd ..
}

# 步骤 4: 设置工作流
setup_workflows() {
    echo "🔧 步骤 4: 检查工作流..."

    cd homebrew-tap

    if [ ! -d ".github/workflows" ]; then
        mkdir -p .github/workflows
        echo "📁 创建了 .github/workflows 目录"
    fi

    # 检查是否有现有的 workflow
    if [ -f ".github/workflows/test.yml" ]; then
        echo "✅ 发现现有的 workflow 文件"
    else
        echo "ℹ️  没有找到测试 workflow，建议添加一个"
    fi

    cd ..
    echo ""
}

# 步骤 5: 显示后续步骤
show_next_steps() {
    echo "🎯 后续步骤:"
    echo ""
    echo "1. 📝 更新 formula 中的 SHA256 值:"
    echo "   ./scripts/update-formula.sh v1.0.0"
    echo ""
    echo "2. 🧪 测试 formula:"
    echo "   cd homebrew-tap"
    echo "   brew audit --strict $FORMULA_DIR/$FORMULA_FILE"
    echo "   brew style $FORMULA_DIR/$FORMULA_FILE"
    echo ""
    echo "3. 📤 提交更改:"
    echo "   cd homebrew-tap"
    echo "   git add ."
    echo "   git commit -m \"Add mdns-reflector-go formula\""
    echo "   git push origin main"
    echo ""
    echo "4. 🏷️  为你的项目创建 GitHub Release (v1.0.0)"
    echo ""
    echo "5. 🧪 测试安装:"
    echo "   brew tap $TAP_NAME"
    echo "   brew install mdns-reflector-go"
    echo ""
}

# 主函数
main() {
    echo "开始设置你的现有 Homebrew Tap..."
    echo ""

    clone_tap_repo
    check_repo_structure
    setup_formula
    setup_workflows

    echo "🎉 Tap 仓库设置完成!"
    echo ""
    echo "📂 本地 tap 目录: $(pwd)/homebrew-tap"
    echo "🔗 远程仓库: $TAP_REPO"
    echo ""

    show_next_steps
}

# 运行主函数
main "$@"

#!/bin/bash

# 创建 Homebrew Tap 仓库的辅助脚本

set -e

echo "🚀 Homebrew Tap 创建助手"
echo "=========================="
echo ""

# 检查是否安装了 Git 和 GitHub CLI
check_dependencies() {
    echo "📋 检查依赖..."

    if ! command -v git &> /dev/null; then
        echo "❌ 请先安装 Git"
        echo "   brew install git"
        exit 1
    fi

    if ! command -v gh &> /dev/null; then
        echo "⚠️  建议安装 GitHub CLI 以便自动化创建仓库"
        echo "   brew install gh"
        echo "   gh auth login"
        echo ""
        echo "如果没有安装，请手动在浏览器中创建仓库"
    fi

    echo "✅ 依赖检查完成"
    echo ""
}

# 获取用户输入
get_user_info() {
    echo "📝 请输入你的信息："

    # 获取 GitHub 用户名
    if command -v gh &> /dev/null; then
        GITHUB_USERNAME=$(gh api user --jq '.login')
        echo "检测到 GitHub 用户名: $GITHUB_USERNAME"
        read -p "使用这个用户名? (y/n): " use_detected
        if [[ $use_detected != "y" && $use_detected != "Y" ]]; then
            read -p "请输入你的 GitHub 用户名: " GITHUB_USERNAME
        fi
    else
        read -p "请输入你的 GitHub 用户名: " GITHUB_USERNAME
    fi

    echo "将创建仓库: $GITHUB_USERNAME/homebrew-tap"
    echo ""
}

# 使用 GitHub CLI 创建仓库
create_repo_with_gh() {
    echo "🔧 使用 GitHub CLI 创建仓库..."

    # 检查是否已登录
    if ! gh auth status &> /dev/null; then
        echo "请先登录 GitHub:"
        echo "gh auth login"
        exit 1
    fi

    # 创建仓库
    gh repo create "$GITHUB_USERNAME/homebrew-tap" \
        --description "Homebrew formulae for my projects" \
        --public \
        --add-readme \
        --disable-wiki \
        --disable-issues \
        --disable-projects

    echo "✅ 仓库创建成功!"
}

# 手动创建指南
manual_creation_guide() {
    echo "🌐 手动创建 GitHub 仓库指南:"
    echo ""
    echo "1. 打开浏览器，访问: https://github.com/new"
    echo ""
    echo "2. 填写仓库信息:"
    echo "   📦 仓库名称: homebrew-tap"
    echo "   📝 描述: Homebrew formulae for my projects"
    echo "   🌍 可见性: Public (公开)"
    echo ""
    echo "3. 取消勾选以下选项:"
    echo "   ❌ Add a README file"
    echo "   ❌ Add .gitignore"
    echo "   ❌ Choose a license"
    echo ""
    echo "4. 点击 'Create repository' 按钮"
    echo ""
    echo "创建完成后，按回车键继续..."
    read -p ""
}

# 克隆仓库
clone_repo() {
    echo "📥 克隆仓库到本地..."

    if [ -d "homebrew-tap" ]; then
        echo "⚠️  本地已存在 homebrew-tap 目录，正在备份..."
        mv homebrew-tap homebrew-tap.backup.$(date +%Y%m%d_%H%M%S)
    fi

    git clone "https://github.com/$GITHUB_USERNAME/homebrew-tap.git"

    echo "✅ 仓库克隆完成!"
    echo ""
}

# 设置仓库
setup_repo() {
    echo "⚙️  设置仓库..."

    cd homebrew-tap

    # 创建 .github/workflows 目录
    mkdir -p .github/workflows

    # 创建基本的 workflow 文件 (可选)
    cat > .github/workflows/test.yml << 'EOF'
name: Test Formulae

on:
  push:
    paths:
      - '**.rb'
  pull_request:
    paths:
      - '**.rb'

jobs:
  test:
    runs-on: macos-latest

    steps:
      - name: Checkout
        uses: actions/checkout@v4

      - name: Set up Homebrew
        uses: Homebrew/actions/setup-homebrew@master

      - name: Test formulae
        run: |
          for formula in *.rb; do
            echo "Testing $formula..."
            brew audit --strict "$formula"
            brew style "$formula"
          done
EOF

    # 提交初始文件
    git add .
    git commit -m "Initial commit: Add basic tap structure"
    git push origin main

    cd ..
    echo "✅ 仓库设置完成!"
    echo ""
}

# 主函数
main() {
    check_dependencies
    get_user_info

    # 询问用户是否使用 GitHub CLI
    if command -v gh &> /dev/null; then
        echo "你想要:"
        echo "1) 使用 GitHub CLI 自动创建 (推荐)"
        echo "2) 手动在浏览器中创建"
        read -p "请选择 (1/2): " choice

        case $choice in
            1)
                create_repo_with_gh
                ;;
            2)
                manual_creation_guide
                ;;
            *)
                echo "❌ 无效选择"
                exit 1
                ;;
        esac
    else
        manual_creation_guide
    fi

    clone_repo
    setup_repo

    echo "🎉 Homebrew Tap 创建完成!"
    echo ""
    echo "📂 本地目录: $(pwd)/homebrew-tap"
    echo "🔗 远程仓库: https://github.com/$GITHUB_USERNAME/homebrew-tap"
    echo ""
    echo "下一步:"
    echo "1. 将项目的 mdns-reflector-go.rb 文件复制到 homebrew-tap 目录"
    echo "2. 更新 formula 中的用户名"
    echo "3. 提交并推送更改"
}

# 运行主函数
main "$@"

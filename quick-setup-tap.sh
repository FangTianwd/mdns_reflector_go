#!/bin/bash

# 快速设置现有 Homebrew Tap 的脚本

set -e

echo "🚀 快速设置你的 Homebrew Tap"
echo "=============================="
echo ""

# 检查是否在项目根目录
if [ ! -f "mdns-reflector-go.rb" ]; then
    echo "❌ 请在项目根目录运行此脚本"
    echo "   项目根目录应该包含 mdns-reflector-go.rb 文件"
    exit 1
fi

# 1. 克隆仓库（如果不存在）
if [ ! -d "homebrew-tap" ]; then
    echo "📥 克隆 homebrew-tap 仓库..."
    git clone https://github.com/FangTianwd/homebrew-tap.git
fi

cd homebrew-tap

# 2. 复制 formula
echo "📋 复制 formula 文件..."
cp ../mdns-reflector-go.rb Formula/

# 3. 更新 formula 信息
echo "🔄 更新 formula 信息..."
FORMULA_FILE="Formula/mdns-reflector-go.rb"

# 更新 homepage
sed -i.bak 's|homepage ".*"|homepage "https://github.com/FangTianwd/mdns-reflector-go"|' "$FORMULA_FILE"

# 更新 url 中的用户名
sed -i.bak 's|url "https://github.com/[^/]*/|url "https://github.com/FangTianwd/|' "$FORMULA_FILE"

# 清理备份文件
rm -f "${FORMULA_FILE}.bak"

echo "✅ Formula 已更新"

# 4. 检查 formula 语法
echo "🔍 检查 formula 语法..."
if command -v brew &> /dev/null; then
    brew audit --strict "$FORMULA_FILE" || echo "⚠️  Formula 检查失败，请手动修复"
    brew style "$FORMULA_FILE" || echo "⚠️  样式检查失败，请手动修复"
else
    echo "⚠️  未找到 brew 命令，跳过语法检查"
fi

echo ""
echo "🎯 下一步操作:"
echo ""
echo "1. 🔍 检查 Formula/mdns-reflector-go.rb 文件是否正确"
echo "2. 📝 手动编辑并设置正确的 sha256 值 (稍后用脚本更新)"
echo "3. 🧪 测试构建: brew install --build-from-source Formula/mdns-reflector-go.rb"
echo "4. 📤 提交更改: git add . && git commit -m \"Add mdns-reflector-go formula\""
echo "5. 🚀 推送: git push origin main"
echo "6. 🏷️ 在 https://github.com/FangTianwd/mdns-reflector-go 创建 Release"
echo "7. 🔄 更新 SHA256: cd .. && ./scripts/update-formula.sh v1.0.0"
echo ""
echo "📂 Tap 目录: $(pwd)"
echo "🔗 仓库: https://github.com/FangTianwd/homebrew-tap"

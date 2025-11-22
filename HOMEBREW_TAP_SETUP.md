# Homebrew Tap 设置指南

本文档介绍如何为你的Go项目创建自己的Homebrew tap。

## 前置要求

- GitHub 账户
- 已发布的项目版本 (GitHub Release)
- Homebrew 已安装 (macOS/Linux)

## 步骤 1: 创建 Homebrew Tap 仓库

1. 在 GitHub 上创建一个新仓库，命名格式为 `homebrew-tap`
   ```
   名称: homebrew-tap
   描述: Homebrew formulae for my projects
   ```

2. 克隆到本地:
   ```bash
   git clone https://github.com/your-username/homebrew-tap.git
   cd homebrew-tap
   ```

## 步骤 2: 创建 Formula

1. 将项目中的 `mdns-reflector-go.rb` 文件复制到 tap 仓库:
   ```bash
   cp /path/to/your/project/mdns-reflector-go.rb .
   ```

2. 更新 formula 中的信息:
   - 将 `your-username` 替换为你的 GitHub 用户名
   - 运行更新脚本获取正确的 SHA256:
     ```bash
     # 在项目根目录运行
     ./scripts/update-formula.sh v1.0.0  # 替换为实际版本号
     ```

3. 提交并推送:
   ```bash
   git add mdns-reflector-go.rb
   git commit -m "Add mdns-reflector-go formula"
   git push origin main
   ```

## 步骤 3: 测试 Formula

1. 在本地测试 formula:
   ```bash
   # 安装前先卸载 (如果已安装)
   brew uninstall mdns-reflector-go

   # 从本地tap安装
   brew install --build-from-source mdns-reflector-go

   # 验证安装
   mdns-reflector-go --help

   # 测试服务
   brew services start mdns-reflector-go
   brew services list | grep mdns-reflector-go
   ```

2. 修复任何问题后，重新提交。

## 步骤 4: 发布新版本

当你发布新版本时:

1. 创建 GitHub Release (会触发自动构建)
2. 更新 formula:
   ```bash
   ./scripts/update-formula.sh v1.1.0  # 新版本号
   ```
3. 提交更新:
   ```bash
   git add mdns-reflector-go.rb
   git commit -m "Update mdns-reflector-go to v1.1.0"
   git push
   ```

## 步骤 5: 让用户安装

用户现在可以通过以下方式安装:

```bash
# 添加你的tap
brew tap your-username/homebrew-tap

# 安装你的工具
brew install mdns-reflector-go

# 启动服务
brew services start mdns-reflector-go
```

## 高级配置

### Formula 中的服务配置

当前的 formula 已经包含了服务配置:

```ruby
service do
  run [opt_bin/"mdns-reflector-go"]
  keep_alive true
  log_path var/"log/mdns-reflector-go.log"
  error_log_path var/"log/mdns-reflector-go-error.log"
end
```

### 添加依赖

如果你的项目需要其他 Homebrew 包，在 formula 中添加:

```ruby
depends_on "some-package"
```

### 瓶装 (Bottles)

对于更好的用户体验，可以设置自动瓶装。但这需要 CI/CD 配置，通常用于受欢迎的项目。

## 故障排除

### 常见问题

1. **SHA256 不匹配**: 重新运行 `update-formula.sh` 脚本
2. **构建失败**: 检查 Go 版本和依赖
3. **服务启动失败**: 检查日志文件 `~/Library/Logs/Homebrew/mdns-reflector-go/`

### 调试命令

```bash
# 查看 formula
brew cat mdns-reflector-go

# 检查语法
brew audit --strict mdns-reflector-go

# 详细安装日志
brew install -v mdns-reflector-go
```

## 贡献给官方 Homebrew

如果你希望你的工具被包含在官方 Homebrew 中:

1. 确保代码质量和测试覆盖
2. 提交 PR 到 [homebrew-core](https://github.com/Homebrew/homebrew-core)
3. 遵循他们的 [formula 规范](https://docs.brew.sh/Formula-Cookbook)

## 参考资料

- [Homebrew Formula Cookbook](https://docs.brew.sh/Formula-Cookbook)
- [Homebrew Services](https://docs.brew.sh/How-to-Service-Your-Formula)
- [Creating a Tap](https://docs.brew.sh/How-to-Create-and-Maintain-a-Tap)

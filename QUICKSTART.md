# 快速开始指南

## 🚀 3 分钟内让 mDNS Reflector 运行起来

### 前置要求
- macOS (支持 Intel 和 Apple Silicon)
- Go 1.19+ (用于构建)
- Homebrew (用于安装)

### 步骤 1: 构建项目

```bash
# 克隆项目
git clone https://github.com/your-username/mdns-reflector-go.git
cd mdns-reflector-go

# 构建二进制文件
make build

# 验证构建成功
./build/mdns-reflector-go --help
```

### 步骤 2: 配置网络接口

```bash
# 查看可用网络接口
ifconfig | grep -E "^\w+:" | awk -F: '{print $1}'

# 配置反射接口 (根据你的网络环境调整)
./build/mdns-reflector-go --config-ifaces en1,bridge100
```

**常见接口配置**:
- **Wi-Fi + Docker**: `en1,bridge100`
- **有线网络 + Docker**: `en0,bridge100`
- **仅 Wi-Fi**: `en1`
- **多网络**: `en0,en1,bridge100`

### 步骤 3: 运行服务

```bash
# 前台运行测试
./build/mdns-reflector-go

# 如果正常工作，按 Ctrl+C 停止
```

### 步骤 4: 安装到系统 (可选)

```bash
# 安装到 /usr/local/bin
make install

# 验证安装
mdns-reflector-go --help
```

### 步骤 5: 设置为系统服务

```bash
# 启动服务
brew services start mdns-reflector-go

# 检查状态
brew services list | grep mdns-reflector-go

# 查看日志
tail -f ~/Library/Logs/Homebrew/mdns-reflector-go/*.log
```

### 步骤 6: 验证工作

```bash
# 检查进程
ps aux | grep mdns-reflector-go

# 测试 mDNS 解析
# 安装 dns-sd (如果没有)
brew install mdns-sd

# 监听 mDNS 流量
dns-sd -B _services._dns-sd._udp

# 在 Docker 容器中测试
docker run --rm alpine ping host.docker.internal
```

---

## 🎯 Homebrew Tap 快速设置

### 为你的项目创建 Tap

1. **创建 Tap 仓库**:
   ```bash
   # 在 GitHub 上创建名为 homebrew-tap 的公开仓库
   ```

2. **克隆并配置**:
   ```bash
   git clone https://github.com/your-username/homebrew-tap.git
   cd homebrew-tap

   # 复制 formula 文件
   cp ../mdns-reflector-go/mdns-reflector-go.rb .

   # 更新 formula 中的用户名
   sed -i 's/your-username/YOUR_USERNAME/g' mdns-reflector-go.rb
   ```

3. **发布版本**:
   ```bash
   # 创建 Git tag
   git tag v1.0.0
   git push origin v1.0.0

   # GitHub Actions 会自动构建发布
   ```

4. **更新 Formula**:
   ```bash
   # 在项目目录运行
   ./scripts/update-formula.sh v1.0.0

   # 提交到 tap 仓库
   cd ../homebrew-tap
   git add mdns-reflector-go.rb
   git commit -m "Update mdns-reflector-go to v1.0.0"
   git push
   ```

### 使用你的 Tap

```bash
# 添加你的 tap
brew tap your-username/homebrew-tap

# 安装
brew install mdns-reflector-go

# 启动服务
brew services start mdns-reflector-go
```

---

## 🔧 故障排除

### 服务启动失败
```bash
# 检查权限 (macOS 会弹出权限对话框)
# 如果没有弹出: 系统偏好设置 -> 安全性与隐私 -> 本地网络

# 查看详细日志
brew services list
tail -f ~/Library/Logs/Homebrew/mdns-reflector-go/*.log
```

### 找不到网络接口
```bash
# 查看所有接口
ifconfig

# Docker Desktop 通常使用 bridge100
# Wi-Fi 通常使用 en0 或 en1
```

### 构建失败
```bash
# 如果遇到 Go 版本问题
brew install go@1.21
export PATH="/opt/homebrew/opt/go@1.21/bin:$PATH"

# 清理并重试
go clean -modcache
go clean -cache
make build
```

---

## 📝 配置示例

### 基本配置
```bash
# 保存配置
mdns-reflector-go --config-ifaces en1,bridge100

# 使用命令行参数
mdns-reflector-go --ifaces en1,bridge100
```

### 配置文件位置
```
~/Library/Application Support/jiangshengcheng.mdns-reflector-go/config.yml
```

### 示例配置文件
```yaml
ifaces:
  - en1      # Wi-Fi 接口
  - bridge100 # Docker 网络接口
```

---

## 🎉 成功标志

当一切正常时，你应该看到:
- ✅ 服务状态为 `started`
- ✅ 日志中显示 "mDNS reflector started"
- ✅ Docker 容器可以解析 `host.docker.internal`
- ✅ 本地网络设备可以被 Docker 容器发现

---

## 📚 更多资源

- [完整文档](README.md)
- [故障排除](TROUBLESHOOTING.md)
- [Homebrew Tap 指南](HOMEBREW_TAP_SETUP.md)
- [GitHub Issues](https://github.com/your-username/mdns-reflector-go/issues)

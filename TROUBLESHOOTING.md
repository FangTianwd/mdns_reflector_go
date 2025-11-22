# 故障排除指南

## 构建问题

### Go 版本兼容性问题

**问题**: 遇到 `compile: version "go1.23.12" does not match go tool version "go1.24.3"` 错误

**原因**: Go 1.24 与某些系统库存在兼容性问题

**解决方案**:

1. **降级 Go 版本** (推荐):
   ```bash
   # 使用 Go 1.21 或 1.22
   brew uninstall go
   brew install go@1.21
   export PATH="/opt/homebrew/opt/go@1.21/bin:$PATH"
   ```

2. **使用 CGO_ENABLED=0**:
   ```bash
   CGO_ENABLED=0 go build -o mdns-reflector-go .
   ```

3. **清理缓存**:
   ```bash
   go clean -modcache
   go clean -cache
   go mod download
   ```

### 网络接口权限问题 (macOS)

**问题**: 服务启动失败，提示网络权限

**解决方案**:
1. 启动服务时会弹出权限对话框，需要点击"允许"
2. 如果没有弹出，可以手动授予:
   ```bash
   # 系统偏好设置 -> 安全性与隐私 -> 隐私 -> 本地网络
   # 找到 mdns-reflector-go 并勾选
   ```

### 配置文件问题

**问题**: 无法加载配置文件

**解决方案**:
1. 检查配置文件路径: `~/Library/Application Support/jiangshengcheng.mdns-reflector-go/config.yml`
2. 确保目录存在:
   ```bash
   mkdir -p "~/Library/Application Support/jiangshengcheng.mdns-reflector-go"
   ```
3. 使用 `--config-ifaces` 参数设置接口:
   ```bash
   mdns-reflector-go --config-ifaces en1,bridge100
   ```

## Homebrew Tap 问题

### Formula SHA256 不匹配

**问题**: `SHA256 mismatch` 错误

**解决方案**:
```bash
# 重新生成 formula
./scripts/update-formula.sh v1.0.0

# 或者手动更新
shasum -a 256 /path/to/archive.tar.gz
```

### Tap 添加失败

**问题**: `brew tap` 命令失败

**解决方案**:
```bash
# 检查仓库是否存在
curl -s https://api.github.com/repos/your-username/homebrew-tap

# 确保仓库是公开的
# 在 GitHub 仓库设置中确认 Visibility 为 Public
```

### 服务启动失败

**问题**: `brew services start mdns-reflector-go` 失败

**解决方案**:
1. 检查日志:
   ```bash
   brew services list
   tail -f ~/Library/Logs/Homebrew/mdns-reflector-go/*.log
   ```

2. 手动测试:
   ```bash
   mdns-reflector-go --ifaces en1,bridge100
   ```

3. 检查权限:
   ```bash
   # 确保二进制文件有执行权限
   chmod +x /opt/homebrew/bin/mdns-reflector-go
   ```

## 网络配置问题

### 如何找到正确的网络接口

**macOS**:
```bash
# 查看所有接口
ifconfig

# 或者使用 networksetup
networksetup -listallhardwareports
```

**查找 mDNS 流量**:
```bash
# 安装 dns-sd (如果没有)
brew install mdns-sd

# 监听 mDNS 流量
dns-sd -B _services._dns-sd._udp
```

### Docker Desktop 集成

对于 Docker Desktop 用户，通常需要桥接接口:
- `bridge100` (Docker Desktop 网络)
- `en0` 或 `en1` (Wi-Fi 接口)

## 性能问题

### CPU 使用率过高

**原因**: 可能是网络接口配置错误或大量 mDNS 流量

**解决方案**:
1. 检查是否配置了正确的接口
2. 减少不必要的接口
3. 使用 `--help` 查看所有选项

### 内存泄漏

**原因**: goroutine 泄漏或缓冲区问题

**解决方案**:
1. 重启服务
2. 检查是否有网络循环
3. 监控资源使用: `ps aux | grep mdns-reflector-go`

## 调试命令

```bash
# 查看服务状态
brew services list | grep mdns

# 查看进程
ps aux | grep mdns-reflector-go

# 查看网络连接
lsof -i :5353

# 测试 mDNS 解析
dig @224.0.0.251 -p 5353 PTR _services._dns-sd._udp local

# 查看系统日志
log show --predicate 'process == "mdns-reflector-go"' --last 1h
```

## 获取帮助

如果问题仍然存在:

1. 检查 GitHub Issues: https://github.com/your-username/mdns-reflector-go/issues
2. 提供以下信息:
   - Go 版本: `go version`
   - macOS 版本: `sw_vers`
   - 错误日志
   - 配置文件内容
   - 网络接口列表: `ifconfig`

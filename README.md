# mDNS Reflector Go

[![Go Version](https://img.shields.io/badge/Go-1.19+-blue.svg)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Release](https://img.shields.io/github/v/release/FangTianwd/mdns_reflector_go)](https://github.com/FangTianwd/mdns_reflector_go/releases)

mDNS 报文反射工具，用于在不同网络接口之间转发多播 DNS 报文。特别适用于 Docker 环境下的服务发现。

## ✨ 功能特性

- 🚀 **高效性能** - 原生 Go 实现，资源占用低
- 🔄 **多接口支持** - 同时监听多个网络接口
- 🐳 **Docker 友好** - 完美支持容器化环境
- 🍺 **开箱即用** - Homebrew 一键安装
- 🔧 **灵活配置** - 支持命令行和配置文件
- 📊 **运行监控** - 内置日志和状态监控

## 📦 安装方式

### Homebrew (推荐)

```bash
# 添加个人 tap
brew tap fangtianwd/homebrew-tap

# 安装最新版本
brew install mdns-reflector-go

# 配置网络接口
mdns-reflector-go --config-ifaces en1,bridge100

# 启动服务
brew services start mdns-reflector-go

# 验证运行
brew services list | grep mdns-reflector-go
```

### 从源码构建

**前置要求：** Go 1.19+

```bash
# 克隆项目
git clone https://github.com/FangTianwd/mdns_reflector_go.git
cd mdns_reflector_go

# 构建
make build

# 安装到系统 (可选)
sudo make install
```

## 🚀 快速开始

### 1. 配置网络接口

首先需要确定要反射的网络接口名称：

```bash
# macOS: 查看网络接口
ifconfig | grep -E "^\w+:" | awk -F: '{print $1}'

# 常见配置：
# - WiFi + Docker: en1,bridge100
# - 有线网 + Docker: en0,bridge100
```

### 2. 运行服务

```bash
# 直接运行
mdns-reflector-go --ifaces en1,bridge100

# 或保存配置后启动服务
mdns-reflector-go --config-ifaces en1,bridge100
brew services start mdns-reflector-go
```

### 3. 验证工作

```bash
# 检查服务状态
ps aux | grep mdns-reflector-go

# 测试 mDNS 解析 (需要安装 dns-sd)
dns-sd -B _services._dns-sd._udp

# Docker 容器内测试
docker run --rm alpine ping host.docker.internal
```

## 📖 使用方法

### 命令行参数

```bash
./mdns-reflector-go [选项]

选项：
  -config-ifaces string
        持久化保存需要反射的网络接口，使用逗号分隔
        示例: -config-ifaces=eth0,wlan0

  -ifaces string
        临时指定需要反射的网络接口，使用逗号分隔
        示例: -ifaces=eth0,wlan0
```

### 配置文件

程序会在用户配置目录创建配置文件：

- **macOS**: `~/Library/Application Support/jiangshengcheng.mdns-reflector-go/config.yml`
- **Linux**: `~/.config/jiangshengcheng.mdns-reflector-go/config.yml`

```yaml
ifaces:
  - eth0
  - wlan0
  - docker0
```

## 🔧 网络接口配置示例

### macOS + Docker Desktop
```bash
mdns-reflector-go --config-ifaces en1,bridge100
```

### Ubuntu + Docker
```bash
mdns-reflector-go --config-ifaces eth0,docker0
```

### 多网络环境
```bash
mdns-reflector-go --config-ifaces eth0,wlan0,docker0
```

## 🛠️ 开发

### 环境要求
- Go 1.19+
- make

### 常用命令

```bash
# 构建
make build          # 构建当前平台
make build-all      # 构建多平台

# 测试
make test

# 清理
make clean

# 安装
make install        # 安装到 /usr/local/bin
```

### 项目结构

```
.
├── main.go              # 主程序
├── Makefile             # 构建脚本
├── go.mod               # Go 模块
├── .github/workflows/   # CI/CD 配置
├── homebrew-tap/        # Homebrew 配置
│   └── Formula/         # Formula 文件
└── README.md           # 文档
```

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

### 开发流程

1. Fork 本项目
2. 创建特性分支: `git checkout -b feature/amazing-feature`
3. 提交更改: `git commit -m 'Add amazing feature'`
4. 推送分支: `git push origin feature/amazing-feature`
5. 提交 Pull Request

## �� 问题反馈

如果遇到问题，请：

1. 查看[故障排除文档](https://github.com/FangTianwd/homebrew-tap)
2. 在 [Issues](https://github.com/FangTianwd/mdns_reflector_go/issues) 中搜索
3. 提交新 Issue，包含：
   - 操作系统版本
   - Docker 版本（如适用）
   - 网络接口配置
   - 错误日志

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🙏 致谢

- [mDNS 协议规范](https://tools.ietf.org/html/rfc6762)
- [Go 语言](https://golang.org/)
- [Homebrew](https://brew.sh/)

---

**注意**: 本工具需要网络接口访问权限。在 macOS 上首次运行时会弹出权限请求，请允许访问。

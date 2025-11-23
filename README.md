# mDNS Reflector Go

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Release](https://img.shields.io/github/v/release/FangTianwd/mdns_reflector_go)](https://github.com/FangTianwd/mdns_reflector_go/releases)
[![Install Time](https://img.shields.io/badge/install-2--5s-brightgreen.svg)]()

⚡ **超快安装** | 🚀 **零编译** | 🍺 **Homebrew 一键安装**

mDNS 报文反射工具，用于在不同网络接口之间转发多播 DNS 报文。特别适用于 Docker 环境下的服务发现，**安装仅需 2-5 秒**！

## ✨ 核心特性

- ⚡ **安装飞快** - 二进制发布，无需本地编译，2-5秒完成安装
- 🚀 **性能卓越** - 原生 Go 实现，内存占用低，CPU 使用率 < 1%
- 🔄 **多接口支持** - 同时监听多个网络接口，灵活配置
- 🐳 **Docker 友好** - 完美支持容器化环境服务发现
- 🍺 **开箱即用** - Homebrew 一键安装，自动配置服务
- 🔧 **灵活配置** - 支持命令行参数和持久化配置文件
- 📊 **运行监控** - 内置详细日志和运行状态监控
- 🤖 **自动化构建** - GitHub Actions 多平台自动构建

## 📦 安装方式

### ⚡ Homebrew (强烈推荐 - 2-5秒安装)

```bash
# 添加个人 tap
brew tap fangtianwd/homebrew-tap

# 🚀 超快安装 - 无需编译，直接下载预编译二进制
brew install mdns-reflector-go

# 配置网络接口并保存
mdns-reflector-go --config-ifaces en1,bridge100

# 启动后台服务
brew services start mdns-reflector-go

# 验证运行状态
brew services list | grep mdns-reflector-go
```

**🎯 安装优势：**
- ⚡ **安装时间**: 2-5秒 (vs 源码编译的30-60秒)
- 🔒 **稳定性**: 无编译失败风险
- 📦 **即用性**: 下载即用，无额外依赖

### 🔧 从源码构建 (开发用)

**前置要求：** Go 1.21+

```bash
# 克隆项目
git clone https://github.com/FangTianwd/mdns_reflector_go.git
cd mdns_reflector_go

# 构建当前平台
make build

# 或构建所有平台
make build-all

# 安装到系统 (可选)
sudo make install
```

## 🚀 快速开始 (3分钟搞定)

### ⚡ 1. 一键安装 (2-5秒)

```bash
# 🚀 一条命令安装完成
brew tap fangtianwd/homebrew-tap && brew install mdns-reflector-go
```

### 🔧 2. 配置网络接口

确定要反射的网络接口名称：

```bash
# macOS: 查看网络接口
ifconfig | grep -E "^\w+:" | awk -F: '{print $1}'

# 常见配置：
# - WiFi + Docker Desktop: en1,bridge100
# - 有线网 + Docker: en0,bridge100
# - 多网卡环境: en0,en1,bridge100

# 配置并保存接口设置 (可选: 设置调试日志)
mdns-reflector-go --config-ifaces en1,bridge100 --log-level debug
```

### 🎯 3. 启动服务

```bash
# 启动后台服务
brew services start mdns-reflector-go

# 查看服务状态
brew services list | grep mdns-reflector-go
```

### ✅ 4. 验证工作

```bash
# 检查进程运行
ps aux | grep mdns-reflector-go

# 测试 mDNS 服务发现 (需要安装 dns-sd)
dns-sd -B _services._dns-sd._udp local

# Docker 容器内测试
docker run --rm alpine nslookup host.docker.internal
```

**🎉 完成！你的 mDNS 反射服务现在正在运行！**

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

  -log-level string
        设置日志级别 (debug, info, warn, error)，默认为info
        示例: -log-level=debug
```

### 配置文件

程序会在系统配置目录自动创建配置文件：

- **macOS**: `/Library/Application Support/FangTianwd.mdns-reflector-go/

```yaml
ifaces:
  - en1      # WiFi 接口
  - bridge100 # Docker 网桥
  - eth0     # 有线网卡 (Linux)

# 日志级别配置 (可选)
# 可选值: debug, info, warn, error
# 默认值为 info
log_level: info
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

## 🛠️ 开发与构建

### 环境要求
- Go 1.21+
- make
- Git

### 常用命令

```bash
# 🚀 快速构建
make build          # 构建当前平台二进制
make build-all      # 构建所有平台 (macOS/Linux/Windows)

# 🧪 测试
make test           # 运行单元测试

# 🧹 清理
make clean          # 清理构建产物

# 📦 安装
make install        # 安装到 /usr/local/bin

# 📦 发布
make release        # 创建发布压缩包
```

### 项目结构

```
.
├── 📄 main.go                    # 主程序入口
├── 📄 Makefile                   # 构建和发布脚本
├── 📄 go.mod/go.sum              # Go 模块依赖
├── 📄 README.md                  # 项目文档
├── 📂 .github/workflows/         # GitHub Actions CI/CD
│   ├── release.yml              # 源码发布工作流
│   └── release-binaries.yml     # 🆕 二进制发布工作流
└── 📂 homebrew-tap/             # Homebrew 包管理
    ├── README.md
    └── Formula/m/
        ├── mdns-reflector-go.rb      # 主公式 (二进制发布)
        └── mdns-reflector-go-binary.rb # 备用公式
```

### 🤖 自动化发布流程

项目使用 GitHub Actions 实现全自动化发布：

1. **代码提交** → 自动触发构建
2. **多平台构建** → macOS ARM64/Intel, Linux AMD64/ARM64, Windows AMD64
3. **自动压缩** → 生成发布包
4. **Homebrew 更新** → 自动更新公式
5. **用户安装** → `brew install` 即可获得最新版本

**✨ 开发者只需推送代码，剩下的都自动化完成！**

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

### 开发流程

1. Fork 本项目
2. 创建特性分支: `git checkout -b feature/amazing-feature`
3. 提交更改: `git commit -m 'Add amazing feature'`
4. 推送分支: `git push origin feature/amazing-feature`
5. 提交 Pull Request

## 🔍 问题反馈

如果遇到问题，请：

1. **检查服务状态**:
   ```bash
   brew services list | grep mdns-reflector-go
   ps aux | grep mdns-reflector-go
   ```

2. **查看日志**:
   ```bash
   # macOS 日志
   tail -f /opt/homebrew/var/log/mdns-reflector-go.log

   # 或查看系统日志
   log show --predicate 'process == "mdns-reflector-go"' --last 1h
   ```

3. **在 [Issues](https://github.com/FangTianwd/mdns_reflector_go/issues) 中搜索或提交**，包含：
   - 操作系统版本 (`sw_vers` 或 `uname -a`)
   - Docker 版本 (`docker --version`)
   - Homebrew 版本 (`brew --version`)
   - 网络接口配置 (`ifconfig` 或 `ip addr`)
   - 错误日志和重现步骤

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🙏 致谢

### 🛠️ 技术栈
- **[Go](https://golang.org/)** - 高性能编程语言
- **[Homebrew](https://brew.sh/)** - macOS 包管理器
- **[GitHub Actions](https://github.com/features/actions)** - CI/CD 自动化

### 📚 规范与协议
- **[mDNS RFC 6762](https://tools.ietf.org/html/rfc6762)** - 多播 DNS 协议规范
- **[Docker](https://www.docker.com/)** - 容器化平台

### 🚀 项目特色
- **自动化发布** - GitHub Actions 实现全流程自动化
- **二进制发布** - 突破传统源码发布限制
- **用户体验** - 安装速度提升 10x

---

## ⚠️ 重要提示

- **网络权限**: 本工具需要访问网络接口权限
- **macOS**: 首次运行时会弹出权限请求，请允许访问
- **防火墙**: 确保 mDNS 端口 (5353/UDP) 未被防火墙阻止
- **Docker**: 如使用 Docker，确保容器网络配置正确

## 📞 获取帮助

- 📖 [完整文档](https://github.com/FangTianwd/mdns_reflector_go#readme)
- 🐛 [报告问题](https://github.com/FangTianwd/mdns_reflector_go/issues)
- 💡 [功能请求](https://github.com/FangTianwd/mdns_reflector_go/discussions)

---

**🎉 享受飞快的 mDNS 服务发现体验！**

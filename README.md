# mDNS Reflector Go

mDNS 报文反射工具，用于在不同网络接口之间转发多播 DNS 报文。

## 功能特性

- 🚀 高效的 mDNS 报文反射
- 🔄 支持多网络接口
- 🐳 Docker 环境友好
- �� Homebrew 安装支持

## 安装方式

### Homebrew (推荐)

```bash
# 添加个人 tap
brew tap fangtianwd/homebrew-tap

# 安装
brew install mdns-reflector-go

# 配置接口
mdns-reflector-go --config-ifaces en1,bridge100

# 启动服务
brew services start mdns-reflector-go
```

### 从源码构建

```bash
# 克隆项目
git clone https://github.com/FangTianwd/mdns_reflector_go.git
cd mdns_reflector_go

# 构建
make build

# 安装
make install
```

## 使用方法

```bash
# 基本用法
./mdns-reflector-go -ifaces=en1,bridge100

# 参数说明
  -config-ifaces string
        持久化需要反射mDNS报文的网络接口，使用逗号分隔
  -ifaces string  
        指定需要反射mDNS报文的网络接口，使用逗号分隔
```

## 接口配置示例

### macOS + Docker Desktop
```bash
mdns-reflector-go --config-ifaces en1,bridge100
```

### Linux 环境
```bash
mdns-reflector-go --config-ifaces eth0,docker0
```

## 开发

```bash
# 运行测试
make test

# 构建多平台二进制
make build-all

# 创建发布包
make release
```

## 许可证

MIT License

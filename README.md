## 安装方式

### Homebrew (推荐)

#### 使用官方tap (如果已合并)
```bash
brew install mdns-reflector-go
```

#### 使用个人tap
```bash
# 添加你的个人tap (将 username 替换为你的GitHub用户名)
brew tap your-username/homebrew-tap

# 安装
brew install mdns-reflector-go

# 配置接口
mdns-reflector-go --config-ifaces en1,bridge100

# 启动服务 (会弹出"本地网络"权限窗口，需要授予)
brew services start mdns-reflector-go

# 验证启动
ps aux | grep mdns
```

### 从源码构建

#### 前置要求
- Go 1.19+

#### 构建步骤
```bash
# 克隆项目
git clone https://github.com/your-username/mdns-reflector-go.git
cd mdns-reflector-go

# 构建
make build

# 安装到系统 (可选)
make install
```

## Usage

    ./mdns-reflector-go -ifaces=en1,bridge100
    
    Usage of ./mdns-reflector-go:
      -config-ifaces string
            持久化需要反射mDNS报文的网络接口，使用逗号分隔，例如：-config-ifaces=eth0,en0
      -ifaces string
            指定需要反射mDNS报文的网络接口，使用逗号分隔，例如：-ifaces=eth0,en0

## FAQ

* 如何知道需要进行反射的ifaces name？
  
  - orbstack侧
    启动orbstack后，`ifconfig`观察输出接口哪个网段跟docker内部的网段匹配
    or
    `dns-sd -B _hap._tcp`后观察 if 列数值（代表interface index），再通过`ip link show` 最前面的数值就是if index了（你可能需要 `brew install iproute2mac` 来使用ip command）
  
  - 本地网络侧
    wifi的话，直接按住option键点击wifi图标，出现一个窗口，接口名称字段就是了

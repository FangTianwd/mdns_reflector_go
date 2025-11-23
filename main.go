package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"gopkg.in/yaml.v3"
)

// 常量定义 - 避免重复计算
const (
	MDNS_IP     = "224.0.0.251"
	MDNS_PORT   = 5353
	BUFFER_SIZE = 65535
	READ_BUFFER = 4096
	CONFIG_PERM = 0644
	DIR_PERM    = 0755
)

// 日志级别常量
const (
	LogLevelDebug = "debug"
	LogLevelInfo  = "info"
	LogLevelWarn  = "warn"
	LogLevelError = "error"
)

// 获取 HOME 环境变量
var (
	homeDir        string
	configFilePath string
)

// Config 表示配置文件结构
type Config struct {
	Ifaces   []string `yaml:"ifaces"`
	LogLevel string   `yaml:"log_level,omitempty"`
}

// Logger 简单的日志包装器
type Logger struct {
	level string
}

// NewLogger 创建新的日志器
func NewLogger(level string) *Logger {
	if level == "" {
		level = LogLevelInfo
	}
	return &Logger{level: level}
}

// shouldLog 检查是否应该记录该级别的日志
func (l *Logger) shouldLog(level string) bool {
	levels := map[string]int{
		LogLevelDebug: 1,
		LogLevelInfo:  2,
		LogLevelWarn:  3,
		LogLevelError: 4,
	}

	currentLevel, exists := levels[l.level]
	if !exists {
		currentLevel = levels[LogLevelInfo]
	}

	targetLevel, exists := levels[level]
	if !exists {
		return false
	}

	return currentLevel <= targetLevel
}

// log 内部日志记录方法
func (l *Logger) log(level, format string, args ...interface{}) {
	if !l.shouldLog(level) {
		return
	}

	prefix := fmt.Sprintf("[%s] %s ", time.Now().Format("2006-01-02 15:04:05"), strings.ToUpper(level))
	message := fmt.Sprintf(format, args...)

	if level == LogLevelError {
		log.Printf(prefix + message)
	} else {
		log.Printf(prefix + message)
	}
}

// Debug 记录调试日志
func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(LogLevelDebug, format, args...)
}

// Info 记录信息日志
func (l *Logger) Info(format string, args ...interface{}) {
	l.log(LogLevelInfo, format, args...)
}

// Warn 记录警告日志
func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(LogLevelWarn, format, args...)
}

// Error 记录错误日志
func (l *Logger) Error(format string, args ...interface{}) {
	l.log(LogLevelError, format, args...)
}

// Reflector 表示mDNS反射器
type Reflector struct {
	config   Config
	conns    sync.Map // 并发安全的连接存储
	mdnsAddr *net.UDPAddr
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
	logger   *Logger
}

// NewReflector 创建新的反射器实例
func NewReflector(logLevel string) (*Reflector, error) {
	ctx, cancel := context.WithCancel(context.Background())

	// 预解析mDNS地址
	mdnsAddr := &net.UDPAddr{
		IP:   net.ParseIP(MDNS_IP),
		Port: MDNS_PORT,
	}

	return &Reflector{
		conns:    sync.Map{},
		mdnsAddr: mdnsAddr,
		ctx:      ctx,
		cancel:   cancel,
		logger:   NewLogger(logLevel),
	}, nil
}

// initConfig 初始化配置路径
func initConfig() error {
	var err error
	homeDir, err = os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("获取用户主目录失败: %w", err)
	}

	configFilePath = filepath.Join(homeDir, "Library/Application Support/FangTianwd.mdns-reflector-go/config.yml")
	return nil
}

// setupSignalHandler 设置信号处理器，支持优雅关闭
func (r *Reflector) setupSignalHandler() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		log.Printf("接收到信号 %v，开始优雅关闭...", sig)
		r.Shutdown()
	}()
}

// loadInterfaces 加载网络接口列表
func (r *Reflector) loadInterfaces(ifaceNames, configIfaces, logLevel string, logger *Logger) ([]*net.Interface, error) {
	var ifaceNameList []string

	if configIfaces != "" {
		// 保存配置模式
		if err := r.saveConfig(configIfaces, logLevel); err != nil {
			return nil, fmt.Errorf("保存配置失败: %w", err)
		}
		logger.Info("配置已保存，程序退出")
		os.Exit(0)
		os.Exit(0)
	}

	// 加载接口列表
	if ifaceNames == "" {
		config, err := r.loadConfig()
		if err != nil {
			return nil, fmt.Errorf("加载配置文件失败: %w", err)
		}

		if len(config.Ifaces) == 0 {
			return nil, fmt.Errorf("必须指定至少一个网络接口，可以通过--config-ifaces参数或config.yml配置文件")
		}

		ifaceNameList = config.Ifaces
		logger.Info("从配置文件加载接口: %v", ifaceNameList)
	} else {
		ifaceNameList = strings.Split(ifaceNames, ",")
		logger.Info("从命令行参数加载接口")
	}

	// 验证并获取接口
	ifaces := make([]*net.Interface, 0, len(ifaceNameList))
	for _, name := range ifaceNameList {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}

		iface, err := net.InterfaceByName(name)
		if err != nil {
			return nil, fmt.Errorf("无法获取接口 %s: %w", name, err)
		}
		ifaces = append(ifaces, iface)
	}

	if len(ifaces) == 0 {
		return nil, fmt.Errorf("没有找到有效的网络接口")
	}

	return ifaces, nil
}

// setupConnections 设置UDP连接
func (r *Reflector) setupConnections(ifaces []*net.Interface) error {
	r.logger.Info("开始设置 %d 个网络接口的连接", len(ifaces))

	for _, iface := range ifaces {
		r.logger.Debug("正在设置接口 %s (%s)", iface.Name, iface.HardwareAddr)

		conn, err := net.ListenMulticastUDP("udp4", iface, r.mdnsAddr)
		if err != nil {
			r.logger.Error("无法在接口 %s 上监听mDNS报文: %v", iface.Name, err)
			return fmt.Errorf("无法在接口 %s 上监听mDNS报文: %w", iface.Name, err)
		}

		// 设置读取缓冲区大小
		if err := conn.SetReadBuffer(READ_BUFFER); err != nil {
			r.logger.Warn("无法设置接口 %s 的读取缓冲区: %v", iface.Name, err)
		} else {
			r.logger.Debug("成功设置接口 %s 的读取缓冲区为 %d 字节", iface.Name, READ_BUFFER)
		}

		r.conns.Store(iface.Name, conn)
		r.logger.Info("接口 %s 配置完成", iface.Name)

		// 启动goroutine读取报文
		r.wg.Add(1)
		go r.readPackets(iface, conn)
		r.logger.Debug("为接口 %s 启动了读取goroutine", iface.Name)
	}

	r.logger.Info("所有接口配置完成，共启动 %d 个监听goroutine", len(ifaces))
	return nil
}

// readPackets 读取并处理UDP报文
func (r *Reflector) readPackets(iface *net.Interface, conn *net.UDPConn) {
	defer r.wg.Done()

	buf := make([]byte, BUFFER_SIZE)

	for {
		select {
		case <-r.ctx.Done():
			r.logger.Info("停止监听接口 %s", iface.Name)
			return
		default:
			n, _, err := conn.ReadFromUDP(buf)
			if err != nil {
				if r.ctx.Err() != nil {
					// 上下文已取消，不需要记录错误
					return
				}
				r.logger.Error("读取接口 %s 报文错误: %v", iface.Name, err)
				continue
			}

			r.logger.Debug("从接口 %s 收到 %d 字节的mDNS报文", iface.Name, n)

			// 转发报文到其他接口
			r.forwardPacket(iface, buf[:n])
		}
	}
}

// forwardPacket 转发报文到其他接口
func (r *Reflector) forwardPacket(srcIface *net.Interface, packet []byte) {
	r.logger.Debug("开始转发来自接口 %s 的 %d 字节报文", srcIface.Name, len(packet))

	r.conns.Range(func(key, value interface{}) bool {
		name := key.(string)
		conn := value.(*net.UDPConn)

		if name == srcIface.Name {
			// 不将报文发送回其来源接口
			r.logger.Debug("跳过来源接口 %s", name)
			return true
		}

		// 发送报文到多播地址
		if _, err := conn.WriteToUDP(packet, r.mdnsAddr); err != nil {
			r.logger.Error("在接口 %s 上发送报文错误: %v", name, err)
		} else {
			r.logger.Debug("成功在接口 %s 上转发报文", name)
		}

		return true
	})
}

// saveConfig 保存配置到文件
func (r *Reflector) saveConfig(ifacesStr string, logLevel string) error {
	ifacesList := strings.Split(ifacesStr, ",")
	for i, iface := range ifacesList {
		ifacesList[i] = strings.TrimSpace(iface)
	}

	config := Config{
		Ifaces:   ifacesList,
		LogLevel: logLevel,
	}

	r.logger.Debug("准备保存配置: 接口=%v, 日志级别=%s", ifacesList, logLevel)

	data, err := yaml.Marshal(&config)
	if err != nil {
		r.logger.Error("序列化配置失败: %v", err)
		return fmt.Errorf("序列化配置失败: %w", err)
	}

	// 创建目录
	dir := filepath.Dir(configFilePath)
	if err := os.MkdirAll(dir, DIR_PERM); err != nil {
		r.logger.Error("创建配置目录失败: %v", err)
		return fmt.Errorf("创建目录失败: %w", err)
	}

	// 写入文件
	if err := os.WriteFile(configFilePath, data, CONFIG_PERM); err != nil {
		r.logger.Error("写入配置文件失败: %v", err)
		return fmt.Errorf("写入配置文件失败: %w", err)
	}

	r.logger.Info("配置已保存到 %s", configFilePath)
	return nil
}

// loadConfig 从文件加载配置
func (r *Reflector) loadConfig() (Config, error) {
	config := Config{}

	// 检查文件是否存在
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		return config, fmt.Errorf("配置文件 config.yml 不存在")
	}

	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return config, fmt.Errorf("读取配置文件失败: %w", err)
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return config, fmt.Errorf("解析配置文件失败: %w", err)
	}

	return config, nil
}

// loadConfigForLogLevel 专门用于加载日志级别配置
func loadConfigForLogLevel() (Config, error) {
	config := Config{}

	// 初始化配置路径
	if err := initConfig(); err != nil {
		return config, err
	}

	// 检查文件是否存在
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		return config, err
	}

	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return config, err
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return config, err
	}

	return config, nil
}

// Run 启动反射器
func (r *Reflector) Run(ifaceNames, configIfaces, logLevel string, logger *Logger) error {
	// 设置信号处理器
	r.setupSignalHandler()

	// 加载接口
	ifaces, err := r.loadInterfaces(ifaceNames, configIfaces, logLevel, logger)
	if err != nil {
		return err
	}

	// 设置连接
	if err := r.setupConnections(ifaces); err != nil {
		return err
	}

	r.logger.Info("mDNS reflector started, 监听 %d 个接口", len(ifaces))

	// 等待关闭信号
	<-r.ctx.Done()

	// 等待所有goroutine结束
	r.wg.Wait()

	r.logger.Info("mDNS reflector stopped")
	return nil
}

// Shutdown 优雅关闭反射器
func (r *Reflector) Shutdown() {
	r.cancel()

	// 关闭所有连接
	r.conns.Range(func(key, value interface{}) bool {
		if conn, ok := value.(*net.UDPConn); ok {
			conn.Close()
		}
		return true
	})
}

func main() {
	// 初始化配置
	if err := initConfig(); err != nil {
		log.Fatal(err)
	}

	// 解析命令行参数
	var ifaceNames string
	var configIfaces string
	var logLevel string

	flag.StringVar(&ifaceNames, "ifaces", "", "指定需要反射mDNS报文的网络接口，使用逗号分隔，例如：-ifaces=eth0,en0")
	flag.StringVar(&configIfaces, "config-ifaces", "", "持久化需要反射mDNS报文的网络接口，使用逗号分隔，例如：-config-ifaces=eth0,en0")
	flag.StringVar(&logLevel, "log-level", "", "设置日志级别 (debug, info, warn, error)，默认为info，例如：-log-level=debug")
	flag.Parse()

	// 如果没有指定日志级别，尝试从配置中读取
	if logLevel == "" {
		config, err := loadConfigForLogLevel()
		if err == nil && config.LogLevel != "" {
			logLevel = config.LogLevel
		}
	}

	// 创建并运行反射器
	reflector, err := NewReflector(logLevel)
	if err != nil {
		log.Fatal(err)
	}

	// 创建临时logger用于配置加载
	tempLogger := NewLogger(logLevel)
	if err := reflector.Run(ifaceNames, configIfaces, logLevel, tempLogger); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"os/user"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"golang.org/x/net/dns/dnsmessage"
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

// getUserConfigPath 获取用户级别的配置路径
func getUserConfigPath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(usr.HomeDir, "Library", "Application Support", "FangTianwd.mdns-reflector-go", "config.yml"), nil
}

// initConfig 初始化配置路径
func initConfig(configPath string) error {
	if configPath != "" {
		// 用户明确指定了配置路径，直接使用
		configFilePath = configPath
	} else {
		// 用户没有指定配置路径，自动查找现有配置文件
		// 1. 优先尝试用户级配置文件（默认使用用户目录）
		if userPath, err := getUserConfigPath(); err == nil {
			if _, err := os.Stat(userPath); err == nil {
				// 用户级配置文件存在，使用它
				configFilePath = userPath
				return nil
			}
		}

		// 2. 尝试系统级配置文件（保持向后兼容性）
		systemPath := "/Library/Application Support/FangTianwd.mdns-reflector-go/config.yml"
		if _, err := os.Stat(systemPath); err == nil {
			// 系统级配置文件存在，使用它
			configFilePath = systemPath
			return nil
		}

		// 3. 都没有找到，使用用户级路径作为默认（用于保存新配置）
		if userPath, err := getUserConfigPath(); err == nil {
			configFilePath = userPath
		} else {
			// 如果获取用户路径失败，返回错误
			return fmt.Errorf("无法获取用户配置路径: %w", err)
		}
	}
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
func (r *Reflector) loadInterfaces(ifaceNames, configIfaces, configLogLevel, logLevel string, logger *Logger) ([]*net.Interface, error) {
	var ifaceNameList []string

	if configIfaces != "" || configLogLevel != "" {
		// 保存配置模式 - 处理接口配置和日志级别配置
		if err := r.saveConfig(configIfaces, configLogLevel); err != nil {
			return nil, fmt.Errorf("保存配置失败: %w", err)
		}
		logger.Info("配置已保存，程序退出")
		os.Exit(0)
	}

	// 加载接口列表
	if ifaceNames == "" {
		config, err := r.loadConfig()
		if err != nil {
			return nil, fmt.Errorf("未找到配置文件: %s\n\n请使用以下方式之一指定网络接口:\n  1. 使用 -ifaces 参数: -ifaces=en0,en1\n  2. 使用 -config-ifaces 参数创建配置文件: -config-ifaces=en0,en1\n  3. 使用 -config 参数指定配置文件路径: -config=/path/to/config.yml\n\n详细用法请查看帮助: -h", configFilePath)
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

	// 创建一个全局的多播监听器，监听所有接口
	r.logger.Info("创建全局mDNS多播监听器")
	conn, err := net.ListenMulticastUDP("udp4", nil, r.mdnsAddr) // nil表示所有接口
	if err != nil {
		r.logger.Error("无法创建多播监听器: %v", err)
		return fmt.Errorf("无法创建多播监听器: %w", err)
	}
	r.logger.Debug("成功创建全局mDNS多播监听器")

	// 设置读取缓冲区
	if err := conn.SetReadBuffer(READ_BUFFER); err != nil {
		r.logger.Warn("无法设置读取缓冲区: %v", err)
	}

	// 为每个接口存储相同的连接（用于发送）
	for _, iface := range ifaces {
		r.logger.Debug("配置接口 %s (%s)", iface.Name, iface.HardwareAddr)
		r.conns.Store(iface.Name, conn)
		r.logger.Info("接口 %s 配置完成", iface.Name)
	}

	// 只启动一个读取goroutine来处理所有接口的报文
	r.wg.Add(1)
	go r.readPacketsGlobal(conn)
	r.logger.Debug("启动了全局读取goroutine")

	r.logger.Info("所有接口配置完成，共启动 1 个监听goroutine")
	return nil
}

// readPacketsGlobal 全局读取并处理UDP报文
func (r *Reflector) readPacketsGlobal(conn *net.UDPConn) {
	defer r.wg.Done()

	buf := make([]byte, BUFFER_SIZE)

	for {
		select {
		case <-r.ctx.Done():
			r.logger.Info("停止全局mDNS监听")
			return
		default:
			r.logger.Debug("等待接收mDNS报文")
			n, srcAddr, err := conn.ReadFromUDP(buf)
			if err != nil {
				if r.ctx.Err() != nil {
					// 上下文已取消，不需要记录错误
					return
				}
				r.logger.Error("读取mDNS报文错误: %v", err)
				continue
			}

			r.logger.Debug("收到 %d 字节的mDNS报文 from %s", n, srcAddr.String())

			// 查找来源接口
			var srcIface *net.Interface
			r.conns.Range(func(key, value interface{}) bool {
				name := key.(string)
				if ifaces, err := net.Interfaces(); err == nil {
					for _, iface := range ifaces {
						if iface.Name == name {
							// 检查报文是否来自这个接口
							if addrs, err := iface.Addrs(); err == nil {
								for _, addr := range addrs {
									if ipnet, ok := addr.(*net.IPNet); ok {
										if ipnet.Contains(srcAddr.IP) {
											srcIface = &iface
											return false
										}
									}
								}
							}
						}
					}
				}
				return true
			})

			if srcIface != nil {
				r.logger.Debug("报文来自接口 %s", srcIface.Name)
				r.forwardPacket(srcIface, buf[:n])
			} else {
				r.logger.Debug("无法确定报文来源接口，转发到所有接口")
				// 如果无法确定来源，转发到所有接口
				r.conns.Range(func(key, value interface{}) bool {
					name := key.(string)
					conn := value.(*net.UDPConn)
					if _, err := conn.WriteToUDP(buf[:n], r.mdnsAddr); err != nil {
						r.logger.Error("转发到接口 %s 错误: %v", name, err)
					} else {
						r.logger.Debug("转发到接口 %s 成功", name)
					}
					return true
				})
			}
		}
	}
}

// readPackets 保留旧的接口特定方法（已不使用）
func (r *Reflector) readPackets(iface *net.Interface, conn *net.UDPConn) {
	// 这个方法现在不使用了，由readPacketsGlobal替代
}

// forwardPacket 转发报文到其他接口
func (r *Reflector) forwardPacket(srcIface *net.Interface, packet []byte) {
	// 解析mDNS报文内容
	queryName, queryType := parseMDNSMessage(packet)

	var successCount, errorCount int
	var successInterfaces []string
	var errorDetails []string

	r.conns.Range(func(key, value interface{}) bool {
		name := key.(string)
		conn := value.(*net.UDPConn)

		if name == srcIface.Name {
			// 不将报文发送回其来源接口
			return true
		}

		// 发送报文到多播地址
		if _, err := conn.WriteToUDP(packet, r.mdnsAddr); err != nil {
			errorCount++
			if r.logger.shouldLog(LogLevelDebug) {
				errorDetails = append(errorDetails, fmt.Sprintf("%s: %v", name, err))
			}
		} else {
			successCount++
			successInterfaces = append(successInterfaces, name)
		}

		return true
	})

	// 详细的转发日志：显示源接口、目标接口和报文内容
	if successCount > 0 && errorCount == 0 {
		r.logger.Info("mDNS转发: %s -> [%s] (%d字节, %s %s)",
			srcIface.Name, strings.Join(successInterfaces, ","), len(packet), queryType, queryName)
	} else if successCount > 0 && errorCount > 0 {
		r.logger.Info("mDNS转发: %s -> [%s] (%d字节, %s %s), 部分失败:%d个接口",
			srcIface.Name, strings.Join(successInterfaces, ","), len(packet), queryType, queryName, errorCount)
		if r.logger.shouldLog(LogLevelDebug) && len(errorDetails) > 0 {
			r.logger.Debug("转发失败详情: %s", strings.Join(errorDetails, "; "))
		}
	} else if errorCount > 0 {
		r.logger.Warn("mDNS转发失败: %s -> 所有目标接口 (%d字节, %s %s)",
			srcIface.Name, len(packet), queryType, queryName)
		if r.logger.shouldLog(LogLevelDebug) && len(errorDetails) > 0 {
			r.logger.Debug("转发失败详情: %s", strings.Join(errorDetails, "; "))
		}
	}
}

// saveConfig 保存配置到文件
func (r *Reflector) saveConfig(ifacesStr string, logLevel string) error {
	// 读取现有配置（如果存在）
	var config Config
	if _, err := os.Stat(configFilePath); err == nil {
		// 配置文件存在，读取现有配置
		if data, err := os.ReadFile(configFilePath); err == nil {
			if err := yaml.Unmarshal(data, &config); err != nil {
				r.logger.Warn("读取现有配置失败，将创建新配置: %v", err)
				config = Config{}
			}
		}
	}

	// 更新接口配置（如果提供了）
	if ifacesStr != "" {
		ifacesList := strings.Split(ifacesStr, ",")
		for i, iface := range ifacesList {
			ifacesList[i] = strings.TrimSpace(iface)
		}
		config.Ifaces = ifacesList
	}

	// 更新日志级别配置（如果提供了）
	if logLevel != "" {
		config.LogLevel = logLevel
	}

	r.logger.Debug("准备保存配置: 接口=%v, 日志级别=%s", config.Ifaces, config.LogLevel)

	data, err := yaml.Marshal(&config)
	if err != nil {
		r.logger.Error("序列化配置失败: %v", err)
		return fmt.Errorf("序列化配置失败: %w", err)
	}

	// 尝试保存到指定的路径
	actualPath, err := r.saveConfigToPath(configFilePath, data)
	if err != nil {
		// 如果是权限错误且使用的是系统级路径，尝试用户级路径
		if isPermissionError(err) && configFilePath == "/Library/Application Support/FangTianwd.mdns-reflector-go/config.yml" {
			r.logger.Warn("系统级配置目录权限不足，尝试使用用户级配置目录")
			userPath, userErr := getUserConfigPath()
			if userErr != nil {
				r.logger.Error("获取用户配置路径失败: %v", userErr)
				return fmt.Errorf("写入配置文件失败: %w", err)
			}

			actualPath, err = r.saveConfigToPath(userPath, data)
			if err != nil {
				r.logger.Error("写入用户级配置文件失败: %v", err)
				return fmt.Errorf("写入配置文件失败: %w", err)
			}

			r.logger.Info("配置已保存到用户级目录: %s", actualPath)
		} else {
			r.logger.Error("写入配置文件失败: %v", err)
			return fmt.Errorf("写入配置文件失败: %w", err)
		}
	} else {
		r.logger.Info("配置已保存到 %s", actualPath)
	}

	return nil
}

// saveConfigToPath 尝试将配置保存到指定路径
func (r *Reflector) saveConfigToPath(configPath string, data []byte) (string, error) {
	// 创建目录
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, DIR_PERM); err != nil {
		return "", err
	}

	// 写入文件
	if err := os.WriteFile(configPath, data, CONFIG_PERM); err != nil {
		return "", err
	}

	return configPath, nil
}

// isPermissionError 检查是否是权限错误
func isPermissionError(err error) bool {
	if err == nil {
		return false
	}
	// 检查是否是权限拒绝错误
	return strings.Contains(err.Error(), "permission denied") ||
		   strings.Contains(err.Error(), "access denied")
}

// parseMDNSMessage 解析mDNS报文内容，返回查询名称和类型
func parseMDNSMessage(packet []byte) (string, string) {
	var msg dnsmessage.Message
	if err := msg.Unpack(packet); err != nil {
		return "未知", "未知"
	}

	// 如果是查询报文
	if len(msg.Questions) > 0 {
		question := msg.Questions[0]
		name := question.Name.String()
		// 去掉末尾的点
		if len(name) > 0 && name[len(name)-1] == '.' {
			name = name[:len(name)-1]
		}
		qtype := question.Type.String()
		return name, fmt.Sprintf("查询:%s", qtype)
	}

	// 如果是响应报文
	if len(msg.Answers) > 0 {
		answer := msg.Answers[0]
		name := answer.Header.Name.String()
		// 去掉末尾的点
		if len(name) > 0 && name[len(name)-1] == '.' {
			name = name[:len(name)-1]
		}
		qtype := answer.Header.Type.String()
		return name, fmt.Sprintf("响应:%s", qtype)
	}

	return "未知", "未知"
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
func loadConfigForLogLevel(configPath string) (Config, error) {
	config := Config{}

	// 初始化配置路径
	if err := initConfig(configPath); err != nil {
		return config, err
	}

	// 检查文件是否存在，如果不存在就返回空配置（日志级别配置是可选的）
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		return config, nil // 不报错，返回空配置
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
func (r *Reflector) Run(ifaceNames, configIfaces, configLogLevel, logLevel string, logger *Logger) error {
	// 设置信号处理器
	r.setupSignalHandler()

	// 加载接口
	ifaces, err := r.loadInterfaces(ifaceNames, configIfaces, configLogLevel, logLevel, logger)
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
	// 解析命令行参数
	var ifaceNames string
	var configIfaces string
	var configLogLevel string
	var logLevel string
	var configFile string

	flag.StringVar(&ifaceNames, "ifaces", "", "指定需要反射mDNS报文的网络接口，使用逗号分隔，例如：-ifaces=eth0,en0")
	flag.StringVar(&configIfaces, "config-ifaces", "", "持久化需要反射mDNS报文的网络接口，使用逗号分隔，例如：-config-ifaces=eth0,en0")
	flag.StringVar(&configLogLevel, "config-log-level", "", "持久化日志级别配置 (debug, info, warn, error)，例如：-config-log-level=debug")
	flag.StringVar(&logLevel, "log-level", "", "设置日志级别 (debug, info, warn, error)，默认为info，例如：-log-level=debug")
	flag.StringVar(&configFile, "config", "", "指定配置文件路径，默认为用户级配置目录")
	flag.Parse()

	// 初始化配置
	if err := initConfig(configFile); err != nil {
		log.Fatal(err)
	}

	// 如果没有指定日志级别，尝试从配置中读取
	if logLevel == "" {
		config, err := loadConfigForLogLevel(configFile)
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
	if err := reflector.Run(ifaceNames, configIfaces, configLogLevel, logLevel, tempLogger); err != nil {
		log.Fatal(err)
	}
}

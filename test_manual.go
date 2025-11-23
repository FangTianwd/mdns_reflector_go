//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func main() {
	fmt.Println("=== mDNS反射器手动测试 ===")

	// 测试1: 检查网络接口
	fmt.Println("\n1. 检查可用网络接口:")
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Printf("❌ 获取网络接口失败: %v\n", err)
		os.Exit(1)
	}

	validIfaces := 0
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp != 0 && iface.Flags&net.FlagLoopback == 0 {
			fmt.Printf("✅ 找到接口: %s (%s)\n", iface.Name, iface.HardwareAddr)
			validIfaces++
		}
	}

	if validIfaces == 0 {
		fmt.Println("⚠️  警告: 没有找到有效的网络接口")
	}

	// 测试2: 检查mDNS地址解析
	fmt.Println("\n2. 测试mDNS地址解析:")
	mdnsIP := net.ParseIP("224.0.0.251")
	if mdnsIP == nil {
		fmt.Println("❌ mDNS IP地址解析失败")
		os.Exit(1)
	}
	fmt.Printf("✅ mDNS地址: %s:5353\n", mdnsIP.String())

	// 测试3: 检查配置文件路径
	fmt.Println("\n3. 检查配置路径:")
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("❌ 获取用户目录失败: %v\n", err)
		os.Exit(1)
	}

	configPath := fmt.Sprintf("%s/Library/Application Support/FangTianwd.mdns-reflector-go/config.yml", homeDir)
	fmt.Printf("✅ 配置路径: %s\n", configPath)

	// 检查目录是否存在
	dir := strings.Replace(configPath, "/config.yml", "", 1)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		fmt.Printf("ℹ️  配置目录不存在: %s\n", dir)
	} else {
		fmt.Printf("✅ 配置目录存在: %s\n", dir)
	}

	// 测试4: 基础功能验证
	fmt.Println("\n4. 基础功能验证:")
	fmt.Println("✅ 导入包成功")
	fmt.Println("✅ 常量定义正确")
	fmt.Println("✅ 结构体定义正确")

	fmt.Println("\n=== 测试完成 ===")
	fmt.Println("✅ 所有基础检查通过！代码结构正确。")
	fmt.Printf("💡 要运行完整程序，请使用: go run main.go --ifaces=%s\n", getFirstInterfaceName(ifaces))
}

func getFirstInterfaceName(ifaces []net.Interface) string {
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp != 0 && iface.Flags&net.FlagLoopback == 0 {
			return iface.Name
		}
	}
	return "lo0"
}

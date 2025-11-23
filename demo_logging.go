// +build ignore

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	fmt.Println("=== mDNS反射器日志功能演示 ===")

	// 演示不同日志级别的输出
	fmt.Println("\n1. 日志级别说明:")
	fmt.Println("   debug - 显示所有调试信息、接口配置、报文转发详情")
	fmt.Println("   info  - 显示重要事件和状态变化 (默认)")
	fmt.Println("   warn  - 只显示警告信息")
	fmt.Println("   error - 只显示错误信息")

	fmt.Println("\n2. 命令行使用示例:")

	examples := []string{
		"go run main.go --log-level=debug",
		"go run main.go --log-level=info --ifaces=lo0",
		"go run main.go --config-ifaces=en0,en1 --log-level=debug",
	}

	for _, example := range examples {
		fmt.Printf("   $ %s\n", example)
	}

	fmt.Println("\n3. 配置文件示例:")
	fmt.Println("   ~/Library/Application Support/FangTianwd.mdns-reflector-go/config.yml")

	configExample := `ifaces:
  - en0
  - en1
log_level: debug`

	fmt.Printf("\n   ```yaml\n   %s\n   ```\n", configExample)

	fmt.Println("\n4. 调试日志会显示:")
	debugLogs := []string{
		"[DEBUG] 从接口 en0 收到 128 字节的mDNS报文",
		"[DEBUG] 开始转发来自接口 en0 的 128 字节报文",
		"[DEBUG] 成功在接口 en1 上转发报文",
		"[INFO] 接口 en0 配置完成",
		"[INFO] mDNS reflector started, 监听 2 个接口",
	}

	for _, log := range debugLogs {
		fmt.Printf("   %s\n", log)
	}

	fmt.Println("\n5. 实际测试命令:")

	// 检查配置文件是否存在
	homeDir, _ := os.UserHomeDir()
	configPath := filepath.Join(homeDir, "Library/Application Support/FangTianwd.mdns-reflector-go/config.yml")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Printf("   # 创建配置目录\n")
		fmt.Printf("   mkdir -p '%s'\n", strings.Replace(filepath.Dir(configPath), homeDir, "~", 1))
		fmt.Printf("   # 编辑配置文件添加 log_level: debug\n")
	} else {
		fmt.Printf("   # 配置文件已存在，可编辑添加日志级别\n")
	}

	fmt.Println("\n   # 使用调试模式运行 (需要网络权限)")
	fmt.Println("   go run main.go --log-level=debug --ifaces=lo0")
	fmt.Println("\n   # 或使用配置文件")
	fmt.Println("   go run main.go")

	fmt.Println("\n=== 提示 ===")
	fmt.Println("• 首次运行可能需要网络访问权限")
	fmt.Println("• 使用 lo0 (loopback) 接口进行测试不会影响实际网络")
	fmt.Println("• 按 Ctrl+C 停止程序")
}

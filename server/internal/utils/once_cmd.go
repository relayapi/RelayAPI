package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"relayapi/server/internal/config"
)

// printUsage 打印使用说明
func printUsage() {
	fmt.Println("\n生成客户端配置文件的用法:")
	fmt.Println("1. 使用默认配置 (localhost:8840):")
	fmt.Println("   relayapi-server --gen \"\"")
	fmt.Println("   或")
	fmt.Println("   relayapi-server --gen=")
	fmt.Println()
	fmt.Println("2. 指定主机和端口:")
	fmt.Println("   relayapi-server --gen example.com:8080")
	fmt.Println()
	fmt.Println("3. 只指定主机 (使用默认端口 8840):")
	fmt.Println("   relayapi-server --gen example.com")
	fmt.Println()
	fmt.Println("4. 查看此帮助信息:")
	fmt.Println("   relayapi-server --gen help")
	fmt.Println()
	fmt.Println("提示: 使用重定向保存配置到文件:")
	fmt.Println("   relayapi-server --gen > config.rai")
	fmt.Println()
}

// OnceCMDGenerateClientConfig 生成客户端配置并直接退出程序
func OnceCMDGenerateClientConfig(genArg string) {
	// 处理帮助命令
	if genArg == "help" || genArg == "-h" || genArg == "--help" {
		printUsage()
		os.Exit(0)
	}

	var host string
	var port int

	// 解析参数
	if genArg != "" {
		parts := strings.Split(genArg, ":")
		if parts[0] != "" {
			host = parts[0]
		}
		if len(parts) > 1 {
			if _, err := fmt.Sscanf(parts[1], "%d", &port); err != nil {
				fmt.Fprintf(os.Stderr, "错误: 无效的端口号: %v\n", err)
				fmt.Println("\n使用 --gen help 查看使用说明")
				os.Exit(1)
			}
		}
	}

	// 生成配置
	cfg, err := config.DefaultClientConfig(host, port)
	if err != nil {
		fmt.Fprintf(os.Stderr, "错误: 生成配置失败: %v\n", err)
		os.Exit(1)
	}

	// 转换为 JSON
	jsonData, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "错误: JSON 转换失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(jsonData))
	os.Exit(0)
}

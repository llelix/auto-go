package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/urfave/cli/v2"
)

// CreateMockServerCommand 创建启动mock服务器的命令
func CreateMockServerCommand() *cli.Command {
	return &cli.Command{
		Name:  "mock-server",
		Usage: "启动测试用的mock服务器",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "port",
				Aliases: []string{"p"},
				Usage:   "指定服务器端口",
				Value:   "8080",
			},
			&cli.BoolFlag{
				Name:    "install",
				Aliases: []string{"i"},
				Usage:   "安装依赖并初始化mock服务器",
				Value:   false,
			},
		},
		Action: executeMockServerCommand,
	}
}

// executeMockServerCommand 执行mock服务器命令
func executeMockServerCommand(c *cli.Context) error {
	install := c.Bool("install")
	port := c.String("port")

	// 获取当前可执行文件的目录
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("获取可执行文件路径失败: %w", err)
	}
	exeDir := filepath.Dir(exePath)
	mockServerDir := filepath.Join(exeDir, "mock_server")

	// 检查mock_server目录是否存在
	if _, err := os.Stat(mockServerDir); os.IsNotExist(err) {
		return fmt.Errorf("mock_server目录不存在，请确保mock_server文件夹位于auto-go根目录下")
	}

	// 如果需要安装依赖
	if install {
		fmt.Println("🔄 正在安装mock_server依赖...")
		installCmd := exec.Command("go", "mod", "download")
		installCmd.Dir = mockServerDir
		installCmd.Stdout = os.Stdout
		installCmd.Stderr = os.Stderr
		
		if err := installCmd.Run(); err != nil {
			return fmt.Errorf("安装依赖失败: %w", err)
		}
		fmt.Println("✅ 依赖安装完成")
	}

	// 获取mock_server的main.go路径
	mainGoPath := filepath.Join(mockServerDir, "main.go")
	
	// 检查main.go是否存在
	if _, err := os.Stat(mainGoPath); os.IsNotExist(err) {
		return fmt.Errorf("mock_server/main.go文件不存在")
	}

	// 切换到mock_server目录
	if err := os.Chdir(mockServerDir); err != nil {
		return fmt.Errorf("切换到mock_server目录失败: %w", err)
	}

	fmt.Printf("🚀 启动Mock服务器，端口: %s\n", port)
	fmt.Println("📝 访问地址: http://localhost:" + port)
	fmt.Println("📄 简单表单: http://localhost:" + port + "/simple-form")
	fmt.Println("📄 复杂表单: http://localhost:" + port + "/complex-form")
	fmt.Println("🎯 拖拽测试: http://localhost:" + port + "/drag-and-drop")
	fmt.Println("ℹ️  信息获取: http://localhost:" + port + "/info-page")
	fmt.Println("🎮 交互测试: http://localhost:" + port + "/interactive-test")
	fmt.Println("按 Ctrl+C 停止服务器")

	// 设置端口环境变量
	os.Setenv("PORT", port)

	// 使用不同的运行命令基于操作系统
	var runCmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		runCmd = exec.Command("go", "run", "main.go")
	default:
		runCmd = exec.Command("go", "run", "main.go")
	}

	// 连接标准输入输出
	runCmd.Stdout = os.Stdout
	runCmd.Stderr = os.Stderr
	runCmd.Stdin = os.Stdin

	// 启动mock服务器
	if err := runCmd.Run(); err != nil {
		return fmt.Errorf("启动mock服务器失败: %w", err)
	}

	return nil
}
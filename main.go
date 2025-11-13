package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mike/auto-go/config"
	"github.com/mike/auto-go/operater"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    "auto-go",
		Usage:   "基于Playwright的浏览器自动化工具",
		Version: "1.0.0",
		Commands: []*cli.Command{
			{
				Name:  "run",
				Usage: "执行自动化任务",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "config",
						Aliases: []string{"c"},
						Usage:   "配置文件路径",
						Value:   "config.json",
					},
					&cli.StringFlag{
						Name:    "tasks",
						Aliases: []string{"t"},
						Usage:   "任务配置文件路径",
						Value:   "tasks.json",
					},
					&cli.BoolFlag{
						Name:    "headless",
						Usage:   "无头模式运行",
						Value:   true,
					},
					&cli.BoolFlag{
						Name:  "interactive",
						Usage: "交互式模式（显示浏览器）",
					},
					&cli.StringFlag{
						Name:    "chrome-path",
						Aliases: []string{"p"},
						Usage:   "指定Chrome浏览器可执行文件路径",
					},
				},
				Action: func(c *cli.Context) error {
					// 加载配置
					appConfig, err := config.LoadConfig(c.String("config"))
					if err != nil {
						return fmt.Errorf("加载配置失败: %w", err)
					}

					// 加载任务
					tasks, err := operater.LoadTasksFromFile(c.String("tasks"))
					if err != nil {
						return fmt.Errorf("加载任务失败: %w", err)
					}

					// 设置浏览器模式
					headless := appConfig.Browser.Headless
					if c.Bool("interactive") {
						headless = false
					}

					return runAutomationWithChrome(appConfig, tasks, headless, c.String("chrome-path"))
				},
			},
			{
				Name:  "init",
				Usage: "初始化配置文件",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "config",
						Aliases: []string{"c"},
						Usage:   "配置文件路径",
						Value:   "config.json",
					},
					&cli.StringFlag{
						Name:    "tasks",
						Aliases: []string{"t"},
						Usage:   "任务配置文件路径",
						Value:   "tasks.json",
					},
				},
				Action: func(c *cli.Context) error {
					// 创建默认配置文件
					if err := config.CreateDefaultConfigFile(c.String("config")); err != nil {
						return fmt.Errorf("创建配置文件失败: %w", err)
					}

					// 创建示例任务文件
					if err := createSampleTasksFile(c.String("tasks")); err != nil {
						return fmt.Errorf("创建任务文件失败: %w", err)
					}

					fmt.Println("✅ 初始化完成！")
					fmt.Printf("📁 配置文件: %s\n", c.String("config"))
					fmt.Printf("📁 任务文件: %s\n", c.String("tasks"))
					return nil
				},
			},
			{
				Name:  "test",
				Usage: "测试浏览器连接",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "headless",
						Usage: "无头模式",
						Value: true,
					},
					&cli.StringFlag{
						Name:    "chrome-path",
						Aliases: []string{"p"},
						Usage:   "指定Chrome浏览器可执行文件路径",
					},
				},
				Action: func(c *cli.Context) error {
					return testBrowserWithChrome(c.Bool("headless"), c.String("chrome-path"))
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// runAutomationWithChrome 执行自动化任务（支持指定Chrome路径）
func runAutomationWithChrome(appConfig *config.Config, tasks []operater.Task, headless bool, chromePath string) error {
	fmt.Println("🚀 开始自动化任务...")

	// 创建浏览器管理器
	bm := operater.NewBrowserManager()
	defer func() {
		if err := bm.Close(); err != nil {
			fmt.Printf("⚠️  关闭浏览器失败: %v\n", err)
		}
	}()

	// 启动浏览器
	fmt.Printf("🌐 启动浏览器 (headless: %v, chromePath: %s)...\n", headless, chromePath)
	if err := bm.LaunchWithExecutable(headless, chromePath); err != nil {
		return fmt.Errorf("启动浏览器失败: %w", err)
	}
	fmt.Println("✅ 浏览器启动成功")

	// 创建任务管理器
	tm := operater.NewTaskManager(bm)

	// 执行任务
	fmt.Printf("📋 准备执行 %d 个任务\n", len(tasks))
	results := tm.ExecuteTasks(tasks)

	// 保存任务结果
	resultFile := fmt.Sprintf("results_%s.json", time.Now().Format("20060102_150405"))
	if err := tm.SaveTaskResults(results, resultFile); err != nil {
		fmt.Printf("⚠️  保存任务结果失败: %v\n", err)
	} else {
		fmt.Printf("📊 任务结果已保存: %s\n", resultFile)
	}

	// 统计结果
	successCount := 0
	for _, result := range results {
		if result.Success {
			successCount++
		}
	}

	fmt.Printf("🎯 任务完成统计: %d/%d 成功\n", successCount, len(results))
	return nil
}

// runAutomation 执行自动化任务（兼容旧版本）
func runAutomation(appConfig *config.Config, tasks []operater.Task, headless bool) error {
	return runAutomationWithChrome(appConfig, tasks, headless, "")
}

// testBrowserWithChrome 测试浏览器连接（支持指定Chrome路径）
func testBrowserWithChrome(headless bool, chromePath string) error {
	fmt.Println("🧪 测试浏览器连接...")

	bm := operater.NewBrowserManager()
	defer func() {
		if err := bm.Close(); err != nil {
			fmt.Printf("⚠️  关闭浏览器失败: %v\n", err)
		}
	}()

	// 启动浏览器
	if err := bm.LaunchWithExecutable(headless, chromePath); err != nil {
		return fmt.Errorf("浏览器启动失败: %w", err)
	}

	fmt.Println("✅ 浏览器连接测试成功")

	// 导航到测试页面
	if err := bm.Navigate("https://www.google.com"); err != nil {
		return fmt.Errorf("导航测试失败: %w", err)
	}

	fmt.Println("✅ 页面导航测试成功")

	// 截取屏幕截图
	if err := bm.Screenshot("test_screenshot.png"); err != nil {
		return fmt.Errorf("截图测试失败: %w", err)
	}

	fmt.Println("✅ 截图功能测试成功")
	fmt.Println("🎉 所有测试通过！")
	return nil
}

// testBrowser 测试浏览器连接（兼容旧版本）
func testBrowser(headless bool) error {
	return testBrowserWithChrome(headless, "")
}

// createSampleTasksFile 创建示例任务文件
func createSampleTasksFile(filename string) error {
	tasks := []operater.Task{
		{
			Name: "示例表单填写",
			URL:  "https://example.com/form",
			FormFields: map[string]string{
				"#name":     "张三",
				"#email":    "zhangsan@example.com",
				"#phone":    "13800138000",
				"#message":  "这是一个测试消息",
			},
			ClickBefore: []string{
				"#agree-terms",
			},
			ClickAfter: []string{
				"#submit-btn",
			},
			WaitTime:   5,
			Screenshot: true,
		},
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("创建任务文件失败: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	
	if err := encoder.Encode(tasks); err != nil {
		return fmt.Errorf("编码任务文件失败: %w", err)
	}

	return nil
}
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

// 应用程序信息
const (
	AppName  = "auto-go"
	AppUsage = "基于Playwright的浏览器自动化工具"
	AppVer   = "1.0.0"
)

func main() {
	app := createApp()
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

// createApp 创建CLI应用程序
func createApp() *cli.App {
	return &cli.App{
		Name:    AppName,
		Usage:   AppUsage,
		Version: AppVer,
		Commands: []*cli.Command{
			createRunCommand(),
			createInitCommand(),
		},
	}
}

// createRunCommand 创建运行命令
func createRunCommand() *cli.Command {
	return &cli.Command{
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
		},
		Action: executeRunCommand,
	}
}

// createInitCommand 创建初始化命令
func createInitCommand() *cli.Command {
	return &cli.Command{
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
		Action: executeInitCommand,
	}
}

// executeRunCommand 执行运行命令
func executeRunCommand(c *cli.Context) error {
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

	// 获取Chrome路径
	chromePath := getChromePath(c, appConfig)

	return executeAutomation(appConfig, tasks, headless, chromePath)
}

// executeInitCommand 执行初始化命令
func executeInitCommand(c *cli.Context) error {
	// 创建默认配置文件
	if err := config.CreateDefaultConfigFile(c.String("config")); err != nil {
		return fmt.Errorf("创建配置文件失败: %w", err)
	}

	// 创建示例任务文件
	if err := createSampleTasksFile(c.String("tasks")); err != nil {
		return fmt.Errorf("创建任务文件失败: %w", err)
	}

	printInitSuccess(c.String("config"), c.String("tasks"))
	return nil
}

// getChromePath 获取Chrome浏览器路径
func getChromePath(c *cli.Context, appConfig *config.Config) string {
	// 优先使用命令行参数
	if chromePath := c.String("chrome-path"); chromePath != "" {
		return chromePath
	}

	// 否则使用配置文件中的路径
	return appConfig.Browser.ExecutablePath
}

// executeAutomation 执行自动化任务
func executeAutomation(appConfig *config.Config, tasks []operater.Task, headless bool, chromePath string) error {
	fmt.Printf("🚀 开始执行%d个任务\n", len(tasks))

	// 创建和管理浏览器
	bm, err := setupBrowser(headless, chromePath)
	if err != nil {
		return fmt.Errorf("浏览器设置失败: %w", err)
	}
	defer cleanupBrowser(bm)

	// 执行任务
	results := executeTasks(bm, tasks)

	// 保存结果
	saveTaskResults(results)

	// 打印统计信息
	printTaskStatistics(results, len(tasks))

	return nil
}

// setupBrowser 设置浏览器
func setupBrowser(headless bool, chromePath string) (*operater.BrowserManager, error) {
	printBrowserStart(headless, chromePath)

	bm := operater.NewBrowserManager()
	if err := bm.LaunchWithExecutable(headless, chromePath); err != nil {
		return nil, fmt.Errorf("启动浏览器失败: %w", err)
	}

	printBrowserSuccess()
	return bm, nil
}

// cleanupBrowser 清理浏览器资源
func cleanupBrowser(bm *operater.BrowserManager) {
	if err := bm.Close(); err != nil {
		fmt.Printf("⚠️  关闭浏览器失败: %v\n", err)
	}
}

// executeTasks 执行所有任务
func executeTasks(bm *operater.BrowserManager, tasks []operater.Task) []operater.TaskResult {
	tm := operater.NewTaskManager(bm)
	return tm.ExecuteTasks(tasks)
}

// saveTaskResults 保存任务结果
func saveTaskResults(results []operater.TaskResult) {
	resultFile := fmt.Sprintf("results_%s.json", time.Now().Format("20060102_150405"))

	tm := &operater.TaskManager{}
	if err := tm.SaveTaskResults(results, resultFile); err != nil {
		fmt.Printf("⚠️  保存任务结果失败: %v\n", err)
	} else {
		fmt.Printf("📊 任务结果已保存: %s\n", resultFile)
	}
}

// createSampleTasksFile 创建示例任务文件
func createSampleTasksFile(filename string) error {
	tasks := []operater.Task{
		{
			Name: "示例表单填写",
			URL:  "https://example.org/form",
			FormFields: map[string]string{
				"#name":    "张三",
				"#email":   "zhangsan@example.com",
				"#phone":   "13800138000",
				"#message": "这是一个测试消息",
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

func printBrowserStart(headless bool, chromePath string) {
	fmt.Printf("🌐 启动浏览器 (headless: %v, chromePath: %s)...\n", headless, chromePath)
}

func printBrowserSuccess() {
	fmt.Println("✅ 浏览器启动成功")
}

func printTaskStatistics(results []operater.TaskResult, totalTasks int) {
	successCount := 0
	for _, result := range results {
		if result.Success {
			successCount++
		}
	}
	fmt.Printf("🎯 任务完成统计: %d/%d 成功\n", successCount, totalTasks)
}

func printInitSuccess(configPath, tasksPath string) {
	fmt.Println("✅ 初始化完成！")
	fmt.Printf("📁 配置文件: %s\n", configPath)
	fmt.Printf("📁 任务文件: %s\n", tasksPath)
}

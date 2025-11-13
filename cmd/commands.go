package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/mike/auto-go/config"
	"github.com/mike/auto-go/internal/automation"
	"github.com/mike/auto-go/internal/print"
	"github.com/mike/auto-go/operater"
	"github.com/urfave/cli/v2"
)

// CreateRunCommand 创建运行命令
func CreateRunCommand() *cli.Command {
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
			&cli.BoolFlag{
				Name:    "interactive",
				Aliases: []string{"i"},
				Usage:   "交互模式，显示浏览器界面",
				Value:   false,
			},
			&cli.StringFlag{
				Name:    "chrome-path",
				Aliases: []string{"cp"},
				Usage:   "指定Chrome可执行文件路径",
				Value:   "",
			},
		},
		Action: executeRunCommand,
	}
}

// CreateInitCommand 创建初始化命令
func CreateInitCommand() *cli.Command {
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

	// 创建自动化执行器
	auto := automation.New(appConfig, tasks, headless, chromePath)
	
	// 执行自动化任务
	return auto.Execute()
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

	print.InitSuccess(c.String("config"), c.String("tasks"))
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

// createSampleTasksFile 创建示例任务文件
func createSampleTasksFile(tasksPath string) error {
	sampleTasks := []operater.Task{
		{
			Name: "示例任务 - 页面导航",
			URL:  "https://example.com",
			WaitTime: 3,
			Screenshot: true,
			Actions: []operater.Action{
				{
					Type: operater.ActionWaitAppear,
					Selector: "h1",
					Timeout: 5,
					ErrorMessage: "等待页面标题出现失败",
				},
				{
					Type: operater.ActionGetText,
					Selector: "h1",
					OutputKey: "pageTitle",
					ErrorMessage: "获取页面标题失败",
				},
			},
		},
		{
			Name: "示例表单填写",
			URL:  "https://example.org/form",
			WaitTime: 5,
			Screenshot: true,
			Actions: []operater.Action{
				{
					Type: operater.ActionWaitAppear,
					Selector: "#agree-terms",
					Timeout: 5,
					ErrorMessage: "等待同意条款复选框出现失败",
				},
				{
					Type: operater.ActionClick,
					Selector: "#agree-terms",
					ErrorMessage: "点击同意条款复选框失败",
				},
				{
					Type: operater.ActionFill,
					Selector: "#name",
					Value: "张三",
					ErrorMessage: "填写姓名失败",
				},
				{
					Type: operater.ActionFill,
					Selector: "#email",
					Value: "zhangsan@example.com",
					ErrorMessage: "填写邮箱失败",
				},
				{
					Type: operater.ActionFill,
					Selector: "#phone",
					Value: "13800138000",
					ErrorMessage: "填写电话失败",
				},
				{
					Type: operater.ActionFill,
					Selector: "#message",
					Value: "这是一个测试消息",
					ErrorMessage: "填写消息失败",
				},
				{
					Type: operater.ActionClick,
					Selector: "#submit-btn",
					ErrorMessage: "点击提交按钮失败",
				},
				{
					Type: operater.ActionWaitAppear,
					Selector: "#success-message",
					Timeout: 10,
					ErrorMessage: "等待成功消息出现失败",
				},
			},
		},
	}
	
	file, err := os.Create(tasksPath)
	if err != nil {
		return fmt.Errorf("创建任务文件失败: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(sampleTasks); err != nil {
		return fmt.Errorf("编码任务文件失败: %w", err)
	}

	return nil
}
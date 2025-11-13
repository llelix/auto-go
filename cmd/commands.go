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
			Name: "示例任务",
			URL:  "https://example.com",
			WaitTime: 3,
			Screenshot: true,
		},
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
package cmd

import (
	_ "embed"
	"fmt"

	"github.com/mike/auto-go/config"
	"github.com/mike/auto-go/internal/automation"
	"github.com/mike/auto-go/internal/logger"
	"github.com/mike/auto-go/internal/operator"
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
	tasks, err := operator.LoadTasksFromFile(c.String("tasks"))
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

	logger.InitSuccess(c.String("config"), c.String("tasks"))
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

var tasksTemplate string

// createSampleTasksFile 创建示例任务文件
func createSampleTasksFile(tasksPath string) error {

	return nil
}

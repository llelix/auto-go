package cmd

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/mike/auto-go/config"
	"github.com/mike/auto-go/internal/automation"
	"github.com/mike/auto-go/internal/logger"
	"github.com/mike/auto-go/internal/operator"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
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
				Usage:   "任务配置文件路径(.yaml)",
				Value:   "tasks.yaml",
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
				Usage:   "任务配置文件路径(.yaml)",
				Value:   "tasks.yaml",
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

	// 智能检测任务文件（优先使用YAML格式）
	tasksFile := c.String("tasks")
	if tasksFile == "tasks.yaml" {
		// 如果用户使用默认值，尝试智能检测文件
		tasksFile = detectTasksFile()
	}

	// 加载任务
	tasks, err := operator.LoadTasksFromFile(tasksFile)
	if err != nil {
		return fmt.Errorf("加载任务失败: %w (文件: %s)", err, tasksFile)
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

	// 智能检测任务文件格式（根据用户指定的文件名决定格式）
	tasksFile := c.String("tasks")
	if err := createSampleTasksFile(tasksFile); err != nil {
		return fmt.Errorf("创建任务文件失败: %w", err)
	}

	logger.InitSuccess(c.String("config"), tasksFile)
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

// createSampleTasksFile 创建示例任务文件（YAML格式）
func createSampleTasksFile(tasksPath string) error {
	// 示例任务配置
	sampleTasks := []operator.Task{
		{
			Name:     "示例任务",
			URL:      "http://localhost:8080",
			WaitTime: 3,
			Actions: []operator.NodeItem{
				{
					Action: &operator.Action{
						Type:     "click",
						Selector: "#button",
						Timeout:  10,
					},
				},
			},
		},
	}

	// 创建文件
	file, err := os.Create(tasksPath)
	if err != nil {
		return fmt.Errorf("创建任务文件失败: %w", err)
	}
	defer file.Close()

	// 使用YAML格式
	content, err := yaml.Marshal(sampleTasks)
	if err != nil {
		return fmt.Errorf("YAML编码失败: %w", err)
	}

	// 添加文件头注释
	header := "# 自动化任务配置文件 (YAML格式)# 使用YAML格式，支持注释，可读性更好"

	_, err = file.WriteString(header + string(content))
	if err != nil {
		return fmt.Errorf("写入YAML文件失败: %w", err)
	}

	return nil
}

// detectTasksFile 智能检测任务文件（YAML格式）
func detectTasksFile() string {
	// 优先检测的文件格式和顺序
	filePriorities := []string{
		"tasks.yaml",              // YAML格式（推荐）
		"tasks.yml",               // YAML格式（备用）
		"tasks_sequence.yaml",     // 序列任务YAML
		"tasks_with_control.yaml", // 控制任务YAML
	}

	for _, filename := range filePriorities {
		if _, err := os.Stat(filename); err == nil {
			return filename
		}
	}

	// 如果没有任何任务文件存在，返回默认的YAML文件名
	return "tasks.yaml"
}

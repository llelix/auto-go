package automation

import (
	"fmt"
	"time"

	"github.com/mike/auto-go/config"
	"github.com/mike/auto-go/internal/logger"
	"github.com/mike/auto-go/internal/operator"
)

// Automation 自动化执行器
type Automation struct {
	Config     *config.Config
	Tasks      []operator.Task
	Headless   bool
	ChromePath string
}

// New 创建新的自动化执行器
func New(cfg *config.Config, tasks []operator.Task, headless bool, chromePath string) *Automation {
	return &Automation{
		Config:     cfg,
		Tasks:      tasks,
		Headless:   headless,
		ChromePath: chromePath,
	}
}

// Execute 执行自动化任务
func (a *Automation) Execute() error {
	logger.StartExecution(len(a.Tasks))

	// 创建和管理浏览器
	bm, err := a.setupBrowser()
	if err != nil {
		return fmt.Errorf("浏览器设置失败: %w", err)
	}
	defer a.cleanupBrowser(bm)

	// 执行任务
	results := a.executeTasks(bm)

	// 保存结果
	a.saveTaskResults(results)

	// 打印统计信息
	logger.TaskStatistics(results, len(a.Tasks))

	return nil
}

// setupBrowser 设置浏览器
func (a *Automation) setupBrowser() (*operator.BrowserManager, error) {
	logger.BrowserStart(a.Headless, a.ChromePath)

	bm := operator.NewBrowserManager()
	if err := bm.LaunchWithExecutable(a.Headless, a.ChromePath); err != nil {
		return nil, fmt.Errorf("启动浏览器失败: %w", err)
	}

	logger.BrowserSuccess()
	return bm, nil
}

// cleanupBrowser 清理浏览器资源
func (a *Automation) cleanupBrowser(bm *operator.BrowserManager) {
	if err := bm.Close(); err != nil {
		fmt.Printf("⚠️  关闭浏览器失败: %v\n", err)
	}
}

// executeTasks 执行所有任务
func (a *Automation) executeTasks(bm *operator.BrowserManager) []logger.TaskResult {
	tm := operator.NewTaskManager(bm)
	return tm.ExecuteTasks(a.Tasks)
}

// saveTaskResults 保存任务结果
func (a *Automation) saveTaskResults(results []logger.TaskResult) {
	resultFile := fmt.Sprintf("results_%s.json", time.Now().Format("20060102_150405"))

	tm := &operator.TaskManager{}
	if err := tm.SaveTaskResults(results, resultFile); err != nil {
		fmt.Printf("⚠️  保存任务结果失败: %v\n", err)
	} else {
		fmt.Printf("📊 任务结果已保存: %s\n", resultFile)
	}
}

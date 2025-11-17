package operator

import (
	"fmt"
	"os"
	"time"

	"github.com/mike/auto-go/internal/logger"
	"gopkg.in/yaml.v3"
)

// ActionType 定义操作类型
type ActionType string

const (
	ActionClick         ActionType = "click"
	ActionFill          ActionType = "fill"
	ActionHover         ActionType = "hover"
	ActionSelect        ActionType = "select"
	ActionScroll        ActionType = "scroll"
	ActionRightClick    ActionType = "right_click"
	ActionDragDrop      ActionType = "drag_drop"
	ActionWaitAppear    ActionType = "wait_appear"
	ActionWaitDisappear ActionType = "wait_disappear"
	ActionGetText       ActionType = "get_text"
	ActionGetAttribute  ActionType = "get_attribute"
)

// Action 定义单个元素操作
type Action struct {
	Type         ActionType `json:"type"`
	Selector     string     `json:"selector"`
	Value        string     `json:"value,omitempty"`
	Target       string     `json:"target,omitempty"`        // 用于拖拽目标或其他需要第二个元素的场景
	Attribute    string     `json:"attribute,omitempty"`     // 用于获取属性
	Timeout      int        `json:"timeout,omitempty"`       // 超时时间(秒)，默认10秒
	OutputKey    string     `json:"output_key,omitempty"`    // 用于存储操作结果的键名
	ErrorMessage string     `json:"error_message,omitempty"` // 自定义错误信息
}

// Task 定义自动化任务
type Task struct {
	Name       string      `json:"name"`
	URL        string      `json:"url"`
	WaitTime   int         `json:"wait_time,omitempty"`
	Screenshot bool        `json:"screenshot,omitempty"`
	Actions    []NodeItem  `json:"actions"` // 灵活操作序列，支持流程控制
}

// TaskManager 管理自动化任务
type TaskManager struct {
	BrowserManager *BrowserManager
}

// NewTaskManager 创建新的任务管理器
func NewTaskManager(bm *BrowserManager) *TaskManager {
	return &TaskManager{
		BrowserManager: bm,
	}
}

// ExecuteTask 执行单个任务
func (tm *TaskManager) ExecuteTask(task Task) logger.TaskResult {
	startTime := time.Now()
	result := logger.TaskResult{
		TaskName:  task.Name,
		StartTime: startTime.Format("2006-01-02 15:04:05"),
	}

	defer func() {
		endTime := time.Now()
		result.EndTime = endTime.Format("2006-01-02 15:04:05")
		result.Duration = endTime.Sub(startTime).Seconds()
	}()

	// 检查任务是否包含操作序列
	if len(task.Actions) == 0 {
		result.Success = false
		result.Error = "任务未定义操作序列(Actions)，请至少添加一个操作"
		return result
	}

	// 导航到指定URL
	if err := tm.BrowserManager.Navigate(task.URL); err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("导航失败: %v", err)
		return result
	}

	// 等待页面加载
	time.Sleep(time.Duration(task.WaitTime) * time.Second)

	// 执行操作序列
	if err := tm.executeActions(task.Actions); err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("执行操作序列失败: %v", err)
		return result
	}

	// 截取屏幕截图
	if task.Screenshot {
		screenshotFile := fmt.Sprintf("screenshots/%s_%s.png", task.Name, time.Now().Format("20060102_150405"))
		if err := os.MkdirAll("screenshots", 0755); err != nil {
			result.Success = false
			result.Error = fmt.Sprintf("创建截图目录失败: %v", err)
			return result
		}

		if err := tm.BrowserManager.Screenshot(screenshotFile); err != nil {
			result.Success = false
			result.Error = fmt.Sprintf("截图失败: %v", err)
			return result
		}
		result.Screenshot = screenshotFile
	}

	result.Success = true
	return result
}

// executeActions 执行操作序列（支持流程控制）
func (tm *TaskManager) executeActions(items []NodeItem) error {
	// 创建控制执行器
	executor := NewControlExecutor(tm)
	
	// 执行节点项序列
	if err := executor.ExecuteNodeItems(items); err != nil {
		return err
	}
	
	return nil
}

// ExecuteTasks 批量执行任务
func (tm *TaskManager) ExecuteTasks(tasks []Task) []logger.TaskResult {
	var results []logger.TaskResult

	for _, task := range tasks {
		fmt.Printf("🚀 开始执行任务: %s\n", task.Name)
		result := tm.ExecuteTask(task)
		results = append(results, result)

		if result.Success {
			fmt.Printf("✅ 任务完成: %s (耗时: %.2fs)\n", task.Name, result.Duration)
		} else {
			fmt.Printf("❌ 任务失败: %s - %s\n", task.Name, result.Error)
		}

		// 任务间等待时间
		time.Sleep(2 * time.Second)
	}

	return results
}

// SaveTaskResults 保存任务结果到JSON文件
func (tm *TaskManager) SaveTaskResults(results []logger.TaskResult, filename string) error {
	return logger.SaveTaskResults(results, filename)
}

// LoadTasksFromFile 从YAML文件加载任务配置
func LoadTasksFromFile(filename string) ([]Task, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("读取YAML文件失败: %w", err)
	}

	// 首先尝试直接解析为Task数组
	var tasks []Task
	if err := yaml.Unmarshal(content, &tasks); err != nil {
		return nil, fmt.Errorf("YAML解码失败: %w", err)
	}

	// 将每个任务中的Action数组转换为NodeItem数组
	for i := range tasks {
		if len(tasks[i].Actions) > 0 {
			// 检查第一个元素是否是有效的NodeItem
			// 如果不是有效的NodeItem，说明是直接的Action对象
			if !tasks[i].Actions[0].IsAction() && !tasks[i].Actions[0].IsControlNode() {
				// 这种情况说明YAML解码没有正确创建NodeItem包装器
				// 我们需要手动创建NodeItem数组
				var nodeItems []NodeItem
				for _, action := range tasks[i].Actions {
					// 创建一个新的NodeItem并设置Action字段
					nodeItems = append(nodeItems, NodeItem{Action: action.Action})
				}
				tasks[i].Actions = nodeItems
			}
		}
	}

	return tasks, nil
}

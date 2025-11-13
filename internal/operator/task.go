package operator

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mike/auto-go/internal/logger"
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
	Name       string   `json:"name"`
	URL        string   `json:"url"`
	WaitTime   int      `json:"wait_time,omitempty"`
	Screenshot bool     `json:"screenshot,omitempty"`
	Actions    []Action `json:"actions"` // 灵活操作序列，必填
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

// executeActions 执行操作序列
func (tm *TaskManager) executeActions(actions []Action) error {
	for i, action := range actions {
		var err error

		switch action.Type {
		case ActionClick:
			err = tm.BrowserManager.Click(action.Selector)

		case ActionFill:
			if action.Value == "" {
				err = fmt.Errorf("fill操作需要提供value参数")
			} else {
				err = tm.BrowserManager.FillForm(map[string]string{action.Selector: action.Value})
			}

		case ActionHover:
			err = tm.BrowserManager.Hover(action.Selector)

		case ActionSelect:
			if action.Value == "" {
				err = fmt.Errorf("select操作需要提供value参数")
			} else {
				err = tm.BrowserManager.SelectOption(action.Selector, action.Value)
			}

		case ActionScroll:
			err = tm.BrowserManager.ScrollToElement(action.Selector)

		case ActionRightClick:
			err = tm.BrowserManager.RightClick(action.Selector)

		case ActionDragDrop:
			if action.Target == "" {
				err = fmt.Errorf("drag_drop操作需要提供target参数")
			} else {
				err = tm.BrowserManager.DragAndDrop(action.Selector, action.Target)
			}

		case ActionWaitAppear:
			timeout := time.Duration(10) * time.Second
			if action.Timeout > 0 {
				timeout = time.Duration(action.Timeout) * time.Second
			}
			err = tm.BrowserManager.WaitForSelector(action.Selector, timeout)

		case ActionWaitDisappear:
			timeout := time.Duration(10) * time.Second
			if action.Timeout > 0 {
				timeout = time.Duration(action.Timeout) * time.Second
			}
			err = tm.BrowserManager.WaitForElementDisappear(action.Selector, timeout)

		case ActionGetText:
			text, getTextErr := tm.BrowserManager.GetText(action.Selector)
			if getTextErr != nil {
				err = getTextErr
			} else {
				log.Printf("📝 获取元素文本: %s = '%s'", action.Selector, text)
				// 如果提供了输出键名，可以在这里存储结果
				if action.OutputKey != "" {
					// 这里可以扩展为将结果存储到某个上下文中
					log.Printf("📋 文本已存储到键: %s", action.OutputKey)
				}
			}

		case ActionGetAttribute:
			if action.Attribute == "" {
				err = fmt.Errorf("get_attribute操作需要提供attribute参数")
			} else {
				attr, getAttrErr := tm.BrowserManager.GetAttribute(action.Selector, action.Attribute)
				if getAttrErr != nil {
					err = getAttrErr
				} else {
					log.Printf("🏷️ 获取元素属性: %s.%s = '%s'", action.Selector, action.Attribute, attr)
					// 如果提供了输出键名，可以在这里存储结果
					if action.OutputKey != "" {
						// 这里可以扩展为将结果存储到某个上下文中
						log.Printf("📋 属性值已存储到键: %s", action.OutputKey)
					}
				}
			}

		default:
			err = fmt.Errorf("不支持的操作类型: %s", action.Type)
		}

		if err != nil {
			if action.ErrorMessage != "" {
				return fmt.Errorf("操作失败 [%d]: %s", i+1, action.ErrorMessage)
			}
			return fmt.Errorf("操作失败 [%d]: %s - %v", i+1, action.Type, err)
		}

		// 操作间添加短暂延迟，提高执行稳定性
		time.Sleep(500 * time.Millisecond)
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

// LoadTasksFromFile 从JSON文件加载任务配置
func LoadTasksFromFile(filename string) ([]Task, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("打开任务文件失败: %w", err)
	}
	defer file.Close()

	var tasks []Task

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&tasks); err != nil {
		return nil, fmt.Errorf("解码任务文件失败: %w", err)
	}

	return tasks, nil
}

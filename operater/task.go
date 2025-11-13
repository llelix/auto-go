package operater

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Task 定义自动化任务
type Task struct {
	Name        string            `json:"name"`
	URL         string            `json:"url"`
	FormFields  map[string]string `json:"form_fields"`
	ClickBefore []string          `json:"click_before,omitempty"`
	ClickAfter  []string          `json:"click_after,omitempty"`
	WaitTime    int               `json:"wait_time,omitempty"`
	Screenshot  bool              `json:"screenshot,omitempty"`
}

// TaskResult 任务执行结果
type TaskResult struct {
	TaskName   string    `json:"task_name"`
	Success    bool      `json:"success"`
	Error      string    `json:"error,omitempty"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	Duration   float64   `json:"duration"`
	Screenshot string    `json:"screenshot,omitempty"`
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
func (tm *TaskManager) ExecuteTask(task Task) TaskResult {
	result := TaskResult{
		TaskName:  task.Name,
		StartTime: time.Now(),
	}

	defer func() {
		result.EndTime = time.Now()
		result.Duration = result.EndTime.Sub(result.StartTime).Seconds()
	}()

	// 导航到指定URL
	if err := tm.BrowserManager.Navigate(task.URL); err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("导航失败: %v", err)
		return result
	}

	// 等待页面加载
	time.Sleep(time.Duration(task.WaitTime) * time.Second)

	// 执行前置点击操作
	for i, selector := range task.ClickBefore {
		if err := tm.BrowserManager.Click(selector); err != nil {
			result.Success = false
			result.Error = fmt.Sprintf("前置点击失败 [%d]: %s - %v", i+1, selector, err)
			return result
		}
	}

	// 填写表单
	if len(task.FormFields) > 0 {
		if err := tm.BrowserManager.FillForm(task.FormFields); err != nil {
			result.Success = false
			result.Error = fmt.Sprintf("表单填写失败: %v", err)
			return result
		}
	}

	// 执行后置点击操作
	for i, selector := range task.ClickAfter {
		if err := tm.BrowserManager.Click(selector); err != nil {
			result.Success = false
			result.Error = fmt.Sprintf("后置点击失败 [%d]: %s - %v", i+1, selector, err)
			return result
		}
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

// ExecuteTasks 批量执行任务
func (tm *TaskManager) ExecuteTasks(tasks []Task) []TaskResult {
	var results []TaskResult
	
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
func (tm *TaskManager) SaveTaskResults(results []TaskResult, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("创建结果文件失败: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	
	if err := encoder.Encode(results); err != nil {
		return fmt.Errorf("编码结果失败: %w", err)
	}

	return nil
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
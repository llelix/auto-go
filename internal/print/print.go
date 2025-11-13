package print

import (
	"fmt"

	"github.com/mike/auto-go/internal/logger"
)

// StartExecution 打印任务开始信息
func StartExecution(taskCount int) {
	fmt.Printf("🚀 开始执行%d个任务\n", taskCount)
}

// BrowserStart 打印浏览器启动信息
func BrowserStart(headless bool, chromePath string) {
	if chromePath != "" {
		fmt.Printf("🌐 使用系统Chrome: %s (无头模式: %v)\n", chromePath, headless)
	} else {
		fmt.Printf("🌐 使用Playwright内置浏览器 (无头模式: %v)\n", headless)
	}
}

// BrowserSuccess 打印浏览器启动成功信息
func BrowserSuccess() {
	fmt.Println("✅ 浏览器启动成功")
}

// TaskStatistics 打印任务统计信息
func TaskStatistics(results []logger.TaskResult, totalTasks int) {
	var successCount, failureCount int

	for _, result := range results {
		if result.Success {
			successCount++
		} else {
			failureCount++
		}
	}

	fmt.Println("\n📊 任务执行统计:")
	fmt.Printf("   总任务数: %d\n", totalTasks)
	fmt.Printf("   成功: %d\n", successCount)
	fmt.Printf("   失败: %d\n", failureCount)
	fmt.Printf("   成功率: %.2f%%\n", float64(successCount)/float64(totalTasks)*100)
}

// InitSuccess 打印初始化成功信息
func InitSuccess(configPath, tasksPath string) {
	fmt.Println("✅ 初始化完成!")
	fmt.Printf("   配置文件: %s\n", configPath)
	fmt.Printf("   任务文件: %s\n", tasksPath)
	fmt.Println("\n请编辑配置文件后使用 'auto-go run' 命令执行任务")
}

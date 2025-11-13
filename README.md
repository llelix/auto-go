# Auto-Go: 基于Playwright的浏览器自动化工具

一个使用Go语言和Playwright实现的强大浏览器自动化工具，专门用于表单填写、网页操作和自动化测试。

## 功能特性

- ✅ **Chrome浏览器自动化** - 基于Playwright的无头浏览器操作
- ✅ **智能表单填写** - 支持CSS选择器定位表单元素
- ✅ **任务配置管理** - JSON格式的任务配置文件
- ✅ **截图功能** - 自动截取操作过程截图
- ✅ **批量任务执行** - 支持批量执行多个自动化任务
- ✅ **命令行界面** - 友好的命令行操作界面
- ✅ **跨平台支持** - Windows/Linux/macOS全平台支持

## 快速开始

### 1. 安装依赖

```bash
go mod download
```

### 2. 初始化配置文件

```bash
go run main.go init
```

这将创建默认的配置文件 `config.json` 和示例任务文件 `tasks.json`。

### 3. 测试浏览器连接

```bash
go run main.go test --interactive
```

### 4. 执行自动化任务

```bash
# 无头模式执行
go run main.go run

# 交互式模式（显示浏览器窗口）
go run main.go run --interactive

# 指定配置文件
go run main.go run --config custom-config.json --tasks my-tasks.json
```

## 配置文件说明

### 应用配置 (config.json)

```json
{
  "browser": {
    "headless": true,
    "user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
    "timeout": 30
  },
  "tasks": {
    "default_wait_time": 5,
    "auto_screenshot": true
  },
  "logging": {
    "level": "info",
    "file": "logs/auto-go.log",
    "console": true
  }
}
```

### 任务配置 (tasks.json)

```json
[
  {
    "name": "用户注册表单",
    "url": "https://example.com/register",
    "form_fields": {
      "#username": "testuser",
      "#email": "test@example.com",
      "#password": "securepassword123"
    },
    "click_before": ["#agree-terms"],
    "click_after": ["#submit-btn"],
    "wait_time": 5,
    "screenshot": true
  }
]
```

## 任务配置字段说明

- **name**: 任务名称（用于标识和日志输出）
- **url**: 目标网页URL
- **form_fields**: 表单字段映射（CSS选择器 -> 填充值）
- **click_before**: 填写表单前需要点击的元素
- **click_after**: 填写表单后需要点击的元素  
- **wait_time**: 页面加载等待时间（秒）
- **screenshot**: 是否截取屏幕截图

## 支持的CSS选择器示例

```json
{
  "#username": "用户名",           // ID选择器
  ".form-input": "表单值",         // 类选择器
  "input[name='email']": "邮箱",   // 属性选择器
  "form > input:first-child": "值" // 复杂选择器
}
```

## 开发指南

### 项目结构

```
auto-go/
├── main.go                 # 主程序入口
├── config/                 # 配置管理模块
│   └── config.go
├── operater/               # 浏览器操作模块
│   ├── browser.go
│   └── task.go
├── internal/               # 内部工具模块
│   └── time/
│       └── time.go
├── config.json             # 应用配置文件
├── tasks.json              # 任务配置文件
└── README.md
```

### 扩展自定义操作

可以在 `operater/browser.go` 中添加新的浏览器操作方法：

```go
// 示例：选择下拉框选项
func (bm *BrowserManager) SelectOption(selector, value string) error {
    return bm.Page.SelectOption(selector, playwright.SelectOptionValues{
        Labels: playwright.StringSlice(value),
    })
}
```

## 常见问题

### Q: 浏览器启动失败怎么办？
A: 确保系统已安装Chrome浏览器，或运行 `go run main.go test` 进行诊断。

### Q: 如何调试选择器问题？
A: 使用交互模式运行，观察浏览器实际页面元素结构。

### Q: 任务执行速度太慢？
A: 调整 `wait_time` 参数，或在配置中减少截图等非必要操作。

## 许可证

MIT License

## 贡献

欢迎提交Issue和Pull Request来改进这个项目！
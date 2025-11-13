# Auto-Go: 基于Playwright的浏览器自动化工具

一个使用Go语言和Playwright实现的强大浏览器自动化工具，支持灵活的操作序列，专门用于表单填写、网页操作和自动化测试。

## 功能特性

- ✅ **Chrome浏览器自动化** - 基于Playwright的无头浏览器操作
- ✅ **灵活操作序列** - 支持点击、填写、滚动、拖拽、等待等多种元素操作
- ✅ **智能元素交互** - 支持CSS选择器定位、文本获取、属性获取等高级功能
- ✅ **任务配置管理** - JSON格式的任务配置文件，支持自定义操作序列
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
    "wait_time": 5,
    "screenshot": true,
    "actions": [
      {
        "type": "wait_appear",
        "selector": "#username",
        "timeout": 5,
        "error_message": "等待用户名输入框出现失败"
      },
      {
        "type": "fill",
        "selector": "#username",
        "value": "testuser",
        "error_message": "填写用户名失败"
      },
      {
        "type": "fill",
        "selector": "#email",
        "value": "test@example.com",
        "error_message": "填写邮箱失败"
      },
      {
        "type": "fill",
        "selector": "#password",
        "value": "securepassword123",
        "error_message": "填写密码失败"
      },
      {
        "type": "click",
        "selector": "#agree-terms",
        "error_message": "点击同意条款复选框失败"
      },
      {
        "type": "click",
        "selector": "#submit-btn",
        "error_message": "点击提交按钮失败"
      },
      {
        "type": "wait_appear",
        "selector": "#success-message",
        "timeout": 10,
        "error_message": "等待成功消息出现失败"
      }
    ]
  }
]
```

## 任务配置字段说明

- **name**: 任务名称（用于标识和日志输出）
- **url**: 目标网页URL
- **actions**: 操作序列数组，定义要执行的元素操作（必填）
  - **type**: 操作类型（如click、fill、select、wait_appear等）
  - **selector**: CSS选择器，用于定位元素
  - **value**: 操作值（如填写的内容或选择的选项）
  - **target**: 目标元素（用于拖拽操作）
  - **attribute**: 属性名（用于获取属性操作）
  - **timeout**: 超时时间（秒），默认为10秒
  - **output_key**: 输出键名，用于存储操作结果
  - **error_message**: 自定义错误信息
- **wait_time**: 页面加载等待时间（秒）
- **screenshot**: 是否截取屏幕截图

## 支持的操作类型

- **click**: 点击元素
- **fill**: 填写表单字段
- **hover**: 鼠标悬停在元素上
- **select**: 从下拉菜单中选择选项
- **scroll**: 滚动到元素可见区域
- **right_click**: 右键点击元素
- **drag_drop**: 拖拽元素到另一个位置
- **wait_appear**: 等待元素出现
- **wait_disappear**: 等待元素消失
- **get_text**: 获取元素的文本内容
- **get_attribute**: 获取元素的属性值

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
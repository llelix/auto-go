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
- ✅ **新增：流程控制结构** - 支持for循环、if条件分支控制
- ✅ **新增：变量管理** - 支持变量存储和引用
- ✅ **新增：表达式求值** - 支持布尔表达式和算术运算
- ✅ **新增：控制流管理** - 支持break、continue控制流

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

### 3. 启动Mock服务器（可选）

如果您需要测试auto-go的自动化功能，可以启动内置的Mock服务器：

```bash
# 使用默认设置启动Mock服务器（端口8080）
cd mockserver && go run main.go 

```

Mock服务器提供了各种测试场景，包括表单填写、拖拽操作、信息获取等。详细说明请参考 [Mock服务器文档](mock_server/README.md)。

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

### 基础任务配置 (tasks.json)

```json
[
  {
    "name": "简单表单测试",
    "url": "http://localhost:8080/simple-form",
    "wait_time": 3,
    "screenshot": true,
    "actions": [
      {
        "type": "wait_appear",
        "selector": "#name",
        "timeout": 5,
        "error_message": "等待姓名输入框出现失败"
      },
      {
        "type": "fill",
        "selector": "#name",
        "value": "测试用户",
        "error_message": "填写姓名失败"
      },
      {
        "type": "fill",
        "selector": "#email",
        "value": "test@example.com",
        "error_message": "填写邮箱失败"
      },
      {
        "type": "fill",
        "selector": "#phone",
        "value": "13800138000",
        "error_message": "填写电话失败"
      },
      {
        "type": "fill",
        "selector": "#message",
        "value": "这是一个测试消息，用于验证auto-go的表单填写功能。",
        "error_message": "填写消息失败"
      }
      ]
    }
]
```

### 流程控制任务配置 (tasks_with_control.json)

```json
[
  {
    "name": "带流程控制的任务",
    "url": "http://localhost:8080/test-page",
    "actions": [
      {
        "type": "for",
        "children": [
          {
            "type": "click",
            "selector": ".button-{{i}}",
            "error_message": "点击按钮 {{i}} 失败"
          }
        ],
        "variable": "i",
        "from": 1,
        "to": 5,
        "step": 1
      },
      {
        "type": "if",
        "condition": "pageTitle == '登录页面'",
        "children": [
          {
            "type": "fill",
            "selector": "#username",
            "value": "testuser",
            "error_message": "填写用户名失败"
          }
        ]
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

### 基础操作类型
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

### 流程控制操作类型
- **for**: for循环控制结构
  - `variable`: 循环变量名
  - `from`: 起始值
  - `to`: 结束值  
  - `step`: 步长（默认为1）
  - `children`: 循环体内的操作序列
- **if**: 条件分支控制结构
  - `condition`: 布尔表达式条件
  - `children`: 条件为真时执行的操作序列
- **break**: 跳出当前循环
- **continue**: 跳过当前循环迭代

### 表达式语法
支持变量引用和布尔表达式：
- **变量引用**: `{{变量名}}`（在selector、value中引用）
- **比较操作**: `==`, `!=`, `>`, `<`, `>=`, `<=`
- **逻辑操作**: `&&`, `||`, `!`
- **算术操作**: 支持基本的数值运算

示例表达式：
- `pageTitle == '登录页面'`
- `notificationCount > 0 && userRole == 'admin'`
- `index >= 1 && index <= 5`

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
├── operator/               # 浏览器操作模块
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

可以在 `operator/browser.go` 中添加新的浏览器操作方法：

```go
// 示例：选择下拉框选项
func (bm *BrowserManager) SelectOption(selector, value string) error {
    return bm.Page.SelectOption(selector, playwright.SelectOptionValues{
        Labels: playwright.StringSlice(value),
    })
}
```

### 流程控制示例

```go
// 使用ControlExecutor执行带流程控制的任务
func (tm *TaskManager) ExecuteTask(task Task) error {
    executor := NewControlExecutor(tm)
    return executor.ExecuteNodeItems(task.Actions)
}
```

## 常见问题

### Q: 浏览器启动失败怎么办？
A: 确保系统已安装Chrome浏览器，或运行 `go run main.go test` 进行诊断。

### Q: 如何调试选择器问题？
A: 使用交互模式运行，观察浏览器实际页面元素结构。

### Q: 任务执行速度太慢？
A: 调整 `wait_time` 参数，或在配置中减少截图等非必要操作。

### Q: 如何使用流程控制功能？
A: 参考 `tasks_with_control.json` 示例文件，使用for循环和if条件分支构建复杂自动化流程。

### Q: 如何调试表达式求值？
A: 在控制执行器中启用详细日志，查看变量赋值和表达式求值结果。

## 许可证

MIT License

## 贡献

欢迎提交Issue和Pull Request来改进这个项目！
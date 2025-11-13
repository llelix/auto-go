# 灵活操作序列使用指南

本指南介绍如何使用auto-go的新灵活操作序列功能，使元素操作更加灵活和强大。

## 功能概述

新的操作序列允许您按顺序执行各种元素操作，包括点击、填写表单、选择下拉框、滚动、拖拽等多种操作类型。相比旧版本只支持简单的点击和填写，新系统提供了更丰富的交互能力。

## 操作类型

支持以下操作类型：

1. `click` - 点击元素
2. `fill` - 填写表单字段
3. `hover` - 鼠标悬停在元素上
4. `select` - 从下拉菜单中选择选项
5. `scroll` - 滚动到元素可见区域
6. `right_click` - 右键点击元素
7. `drag_drop` - 拖拽元素到另一个位置
8. `wait_appear` - 等待元素出现
9. `wait_disappear` - 等待元素消失
10. `get_text` - 获取元素的文本内容
11. `get_attribute` - 获取元素的属性值

## 操作配置参数

每个操作可以包含以下参数：

- `type` (必填): 操作类型
- `selector` (必填): CSS选择器，用于定位元素
- `value` (可选): 操作值，如填写的内容或选择的选项
- `target` (可选): 目标元素，用于拖拽操作
- `attribute` (可选): 属性名，用于获取属性值操作
- `timeout` (可选): 超时时间(秒)，默认为10秒
- `output_key` (可选): 输出键名，用于存储操作结果
- `error_message` (可选): 自定义错误信息

## 示例任务配置

### 基本表单填写

```json
{
  "name": "基本表单填写",
  "url": "https://example.com/form",
  "wait_time": 3,
  "screenshot": true,
  "actions": [
    {
      "type": "fill",
      "selector": "#username",
      "value": "testuser",
      "error_message": "用户名填写失败"
    },
    {
      "type": "fill",
      "selector": "#email",
      "value": "test@example.com",
      "error_message": "邮箱填写失败"
    },
    {
      "type": "click",
      "selector": "#submit-btn",
      "error_message": "提交按钮点击失败"
    }
  ]
}
```

### 复杂交互操作

```json
{
  "name": "复杂交互操作",
  "url": "https://example.com/advanced-form",
  "wait_time": 2,
  "actions": [
    {
      "type": "wait_appear",
      "selector": "#form-container",
      "timeout": 5,
      "error_message": "等待表单容器出现失败"
    },
    {
      "type": "scroll",
      "selector": "#advanced-options",
      "error_message": "滚动到高级选项失败"
    },
    {
      "type": "hover",
      "selector": "#help-tooltip",
      "error_message": "悬停帮助提示失败"
    },
    {
      "type": "select",
      "selector": "#country",
      "value": "China",
      "error_message": "选择国家失败"
    },
    {
      "type": "wait_appear",
      "selector": "#state-dropdown",
      "timeout": 5,
      "error_message": "等待省份下拉菜单出现失败"
    },
    {
      "type": "select",
      "selector": "#state",
      "value": "Beijing",
      "error_message": "选择省份失败"
    }
  ]
}
```

### 拖拽操作

```json
{
  "name": "拖拽操作示例",
  "url": "https://example.com/drag-and-drop",
  "actions": [
    {
      "type": "wait_appear",
      "selector": "#draggable-element",
      "error_message": "等待可拖拽元素出现失败"
    },
    {
      "type": "drag_drop",
      "selector": "#draggable-element",
      "target": "#drop-zone",
      "error_message": "拖拽操作失败"
    }
  ]
}
```

### 信息获取操作

```json
{
  "name": "信息获取示例",
  "url": "https://example.com/product-page",
  "actions": [
    {
      "type": "get_text",
      "selector": "#product-name",
      "output_key": "productName",
      "error_message": "获取产品名称失败"
    },
    {
      "type": "get_text",
      "selector": "#product-price",
      "output_key": "productPrice",
      "error_message": "获取产品价格失败"
    },
    {
      "type": "get_attribute",
      "selector": "#product-image",
      "attribute": "src",
      "output_key": "productImageSrc",
      "error_message": "获取产品图片链接失败"
    }
  ]
}
```

## 任务配置要求

新的操作序列系统要求每个任务必须包含`actions`数组，并定义至少一个操作：

```json
{
  "name": "任务名称",
  "url": "https://example.com",
  "actions": [
    {
      "type": "操作类型",
      "selector": "CSS选择器",
      "value": "操作值（可选）",
      "error_message": "自定义错误信息"
    }
  ],
  "wait_time": 3,
  "screenshot": true
}
```

## 使用建议

1. **等待元素出现**: 在操作元素前，先使用`wait_appear`确保元素已经出现
2. **添加操作间隔**: 系统会在每个操作后自动添加500ms的间隔，您可以通过`wait_time`参数增加页面级别的等待时间
3. **使用自定义错误信息**: 为每个操作提供清晰的`error_message`，便于调试
4. **合理设置超时**: 对于加载慢的页面或元素，使用`timeout`参数调整等待时间
5. **使用截图功能**: 在关键步骤启用截图，便于排查问题

## 测试方法

1. 使用提供的示例文件`example-flexible-tasks.json`进行测试
2. 使用测试程序`test-flexible-operations.go`验证基本功能
3. 在非无头模式下测试，观察浏览器实际行为

## 进阶功能

您可以通过扩展`Action`类型和`executeActions`函数来添加更多自定义操作。例如添加键盘操作、文件上传等特殊交互。

## 故障排除

1. **元素找不到**: 检查CSS选择器是否正确，使用浏览器开发者工具验证
2. **操作失败**: 检查元素是否在可交互状态，可能需要先滚动或等待
3. **超时问题**: 增加超时时间或检查网络状况
4. **执行顺序问题**: 确认操作顺序符合页面加载逻辑
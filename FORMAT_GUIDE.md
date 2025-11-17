# YAML格式配置指南

## 概述

Auto-Go项目已全面采用YAML格式作为任务配置的标准格式，相比JSON格式具有更好的可读性和维护性。

## YAML vs JSON 优势对比

| 特性 | JSON | YAML | 优势说明 |
|------|------|------|----------|
| 语法复杂度 | 高 | 低 | YAML使用缩进替代大括号，更简洁 |
| 注释支持 | 不支持 | 支持 | YAML支持`#`开头的注释 |
| 多行文本 | 困难 | 容易 | YAML使用`|`符号支持多行文本 |
| 数据类型 | 需要引号 | 自动推断 | YAML自动识别布尔值、数字等 |
| 可读性 | 一般 | 优秀 | YAML结构更清晰，易于理解 |

## 基本语法规则

### 缩进规则
- 使用**2个空格**作为缩进（不要使用tab）
- 层级关系通过缩进表示

### 注释
- 以`#`开头的是注释
- 注释可以出现在任意位置

```yaml
# 这是一个任务配置示例
- name: "测试任务"  # 任务名称
  url: "https://example.com"
  actions:
    - type: "click"  # 点击操作
      selector: "#button"
```

### 数据类型
- 字符串：可以不加引号（除非包含特殊字符）
- 数字：自动识别为数值类型
- 布尔值：`true/false`或`yes/no`
- 数组：使用`-`符号表示

## 任务配置示例

### 基础任务配置

```yaml
- name: "简单表单测试"
  url: "http://localhost:8080/simple-form"
  wait_time: 3
  screenshot: true
  actions:
    - type: "wait_appear"
      selector: "#name"
      timeout: 5
      error_message: "等待姓名输入框出现失败"
    
    - type: "fill"
      selector: "#name"
      value: "测试用户"
      error_message: "填写姓名失败"
    
    - type: "fill"
      selector: "#message"
      value: |
        这是一个测试消息，用于验证auto-go的表单填写功能。
        多行文本在YAML中更加清晰易读。
      error_message: "填写消息失败"
```

### 流程控制任务配置

```yaml
- name: "带流程控制的任务"
  url: "http://localhost:8080/test-page"
  actions:
    # for循环示例
    - type: "for"
      variable: "i"
      from: 1
      to: 5
      children:
        - type: "click"
          selector: ".button-{{i}}"
          error_message: "点击按钮 {{i}} 失败"
    
    # if条件分支示例
    - type: "if"
      condition: "pageTitle == '登录页面'"
      children:
        - type: "fill"
          selector: "#username"
          value: "testuser"
          error_message: "填写用户名失败"
```

### 序列任务配置

```yaml
- name: "用户注册流程"
  url: "http://localhost:8080/register"
  actions:
    - type: "fill"
      selector: "#username"
      value: "newuser123"
      error_message: "填写用户名失败"
    
    - type: "fill"
      selector: "#email"
      value: "newuser@example.com"
      error_message: "填写邮箱失败"
    
    - type: "fill"
      selector: "#password"
      value: "SecurePassword123!"
      error_message: "填写密码失败"
    
    - type: "fill"
      selector: "#confirm-password"
      value: "SecurePassword123!"
      error_message: "确认密码失败"
    
    - type: "click"
      selector: "#submit-btn"
      error_message: "提交注册失败"
```

## 高级特性

### 多行文本支持

YAML非常适合包含多行文本的场景，如填写长消息或复杂CSS选择器：

```yaml
- type: "fill"
  selector: "#description"
  value: |
    这是一个很长的描述文本，
    包含多行内容。
    在YAML中使用|符号可以
    很好地处理这种情况。
  error_message: "填写描述失败"
```

### 变量引用

支持在YAML配置中使用变量引用：

```yaml
- type: "fill"
  selector: "input[name='{{fieldName}}']"
  value: "{{fieldValue}}"
  error_message: "填写{{fieldName}}失败"
```

### 注释和文档

YAML支持丰富的注释，便于文档化：

```yaml
# 用户登录任务
- name: "用户登录"
  url: "http://localhost:8080/login"
  
  # 登录操作序列
  actions:
    # 1. 填写用户名
    - type: "fill"
      selector: "#username"
      value: "testuser"
      error_message: "用户名填写失败"
    
    # 2. 填写密码
    - type: "fill"
      selector: "#password"
      value: "password123"
      error_message: "密码填写失败"
    
    # 3. 点击登录按钮
    - type: "click"
      selector: "#login-btn"
      error_message: "登录按钮点击失败"
```

## 常见问题

### Q: YAML文件编码要求是什么？
A: 必须使用UTF-8编码，确保文件保存为UTF-8格式。

### Q: 缩进应该使用空格还是tab？
A: 必须使用空格（推荐2个空格），不要使用tab。

### Q: 如何验证YAML语法是否正确？
A: 可以使用在线YAML验证工具（如yamlvalidator.com）或在命令行使用`yamllint`工具。

### Q: YAML文件可以包含哪些特殊字符？
A: 如果值包含特殊字符（如冒号、引号等），建议使用引号包围：
```yaml
value: "包含:特殊字符的值"
```

### Q: 如何处理复杂的CSS选择器？
A: YAML支持多行字符串，可以清晰表达复杂选择器：
```yaml
selector: |
  div.container > form:first-child 
  input[type='text'][name='username']
```

## 迁移指南

### 从JSON迁移到YAML

1. **转换基础结构**：将JSON的大括号和方括号转换为YAML的缩进结构
2. **移除引号和逗号**：YAML不需要大多数引号和逗号
3. **添加注释**：利用YAML的注释功能添加说明
4. **验证格式**：使用YAML验证工具检查语法

### 转换示例

**JSON格式**:
```json
{
  "name": "测试任务",
  "url": "https://example.com",
  "actions": [
    {
      "type": "click",
      "selector": "#button"
    }
  ]
}
```

**转换为YAML**:
```yaml
- name: "测试任务"
  url: "https://example.com"
  actions:
    - type: "click"
      selector: "#button"
```

## 最佳实践

1. **使用一致的缩进**：始终使用2个空格缩进
2. **添加有意义的注释**：解释复杂的配置逻辑
3. **分组相关配置**：使用空行分隔不同的配置块
4. **验证语法**：定期使用YAML验证工具检查
5. **版本控制友好**：YAML的差异比JSON更容易阅读

通过采用YAML格式，Auto-Go的配置管理变得更加清晰、易维护，大大提升了开发效率。
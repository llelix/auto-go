# 使用自定义Chrome浏览器路径

本项目支持指定系统安装的Chrome浏览器，而不是使用Playwright内置的浏览器。

## 配置方式

### 方式一：通过配置文件

编辑 `config.json` 文件，添加 `executable_path` 字段：

```json
{
  "browser": {
    "headless": false,
    "user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
    "timeout": 30,
    "executable_path": "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe"
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

### 方式二：通过命令行参数

```bash
# 运行任务时指定Chrome路径
auto-go run --chrome-path "C:\Program Files\Google\Chrome\Application\chrome.exe" --interactive

# 测试时指定Chrome路径
auto-go test --chrome-path "C:\Program Files\Google\Chrome\Application\chrome.exe" --interactive
```

## 常见Chrome路径

### Windows
- Windows 10/11: `C:\Program Files\Google\Chrome\Application\chrome.exe`
- Windows 32位: `C:\Program Files (x86)\Google\Chrome\Application\chrome.exe`

### macOS
- `/Applications/Google Chrome.app/Contents/MacOS/Google Chrome`

### Linux
- `/usr/bin/google-chrome`
- `/usr/bin/google-chrome-stable`
- `/snap/bin/chromium`

## 优先级

1. 命令行参数 `--chrome-path` 或 `-p`
2. 配置文件中的 `browser.executable_path`
3. 默认使用Playwright内置浏览器

## 使用示例

```bash
# 初始化项目（创建默认配置文件）
auto-go init

# 使用默认配置（配置文件中指定的Chrome路径）
auto-go run --interactive

# 覆盖配置文件中的Chrome路径
auto-go run --chrome-path "D:\OtherChrome\chrome.exe" --interactive

# 测试浏览器连接
auto-go test --interactive
```

## 注意事项

1. 确保指定的Chrome浏览器可执行文件路径是正确的
2. Chrome浏览器版本需要与Playwright兼容
3. 在Linux系统上，可能需要添加额外的启动参数

## 问题排查

如果遇到浏览器启动失败，请检查：

1. Chrome浏览器是否正确安装
2. 文件路径是否正确
3. 是否有足够的系统权限
4. 浏览器版本是否兼容

可以通过查看日志文件 `logs/auto-go.log` 获取更多错误信息。
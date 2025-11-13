package operater

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/playwright-community/playwright-go"
)

// BrowserManager 管理浏览器实例
type BrowserManager struct {
	Browser playwright.Browser
	Page    playwright.Page
	Context playwright.BrowserContext
}

// NewBrowserManager 创建新的浏览器管理器
func NewBrowserManager() *BrowserManager {
	return &BrowserManager{}
}

// Launch 启动浏览器
func (bm *BrowserManager) Launch(headless bool) error {
	return bm.LaunchWithExecutable(headless, "")
}

// LaunchWithExecutable 使用指定可执行文件启动浏览器
func (bm *BrowserManager) LaunchWithExecutable(headless bool, executablePath string) error {
	pw, err := playwright.Run()
	if err != nil {
		return fmt.Errorf("启动Playwright失败: %w", err)
	}

	// 构建启动选项
	launchOptions := playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(headless),
		Args: []string{
			"--no-sandbox",
			"--disable-dev-shm-usage",
			"--disable-web-security",
			"--disable-features=VizDisplayCompositor",
		},
	}

	// 如果指定了可执行文件路径，使用系统Chrome
	if executablePath != "" {
		launchOptions.ExecutablePath = playwright.String(executablePath)
		fmt.Printf("🌐 使用系统Chrome: %s", executablePath)
	} else {
		fmt.Println("🌐 使用Playwright内置浏览器")
	}

	// 启动浏览器
	browser, err := pw.Chromium.Launch(launchOptions)
	if err != nil {
		return fmt.Errorf("启动浏览器失败: %w", err)
	}

	bm.Browser = browser

	// 创建浏览器上下文
	context, err := browser.NewContext(playwright.BrowserNewContextOptions{
		Viewport: &playwright.Size{
			Width:  1920,
			Height: 1080,
		},
	})
	if err != nil {
		return fmt.Errorf("创建浏览器上下文失败: %w", err)
	}

	bm.Context = context

	// 创建新页面
	page, err := context.NewPage()
	if err != nil {
		return fmt.Errorf("创建页面失败: %w", err)
	}

	bm.Page = page
	return nil
}

// Navigate 导航到指定URL
func (bm *BrowserManager) Navigate(url string) error {
	if bm.Page == nil {
		return fmt.Errorf("页面未初始化")
	}

	_, err := bm.Page.Goto(url, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	})
	if err != nil {
		return fmt.Errorf("导航到 %s 失败: %w", url, err)
	}

	return nil
}

// WaitForSelector 等待选择器出现
func (bm *BrowserManager) WaitForSelector(selector string, timeout time.Duration) error {
	if bm.Page == nil {
		return fmt.Errorf("页面未初始化")
	}

	_, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	_, err := bm.Page.WaitForSelector(selector, playwright.PageWaitForSelectorOptions{
		State:   playwright.WaitForSelectorStateAttached,
		Timeout: playwright.Float(float64(timeout.Milliseconds())),
	})
	if err != nil {
		return fmt.Errorf("等待元素 %s 超时: %w", selector, err)
	}

	return nil
}

// FillForm 填写表单
func (bm *BrowserManager) FillForm(fields map[string]string) error {
	if bm.Page == nil {
		return fmt.Errorf("页面未初始化")
	}

	for selector, value := range fields {
		// 等待元素出现
		if err := bm.WaitForSelector(selector, 10*time.Second); err != nil {
			return fmt.Errorf("等待表单元素 %s 失败: %w", selector, err)
		}

		// 填写表单
		if err := bm.Page.Fill(selector, value); err != nil {
			return fmt.Errorf("填写表单元素 %s 失败: %w", selector, err)
		}

		log.Printf("✅ 已填写表单元素: %s = %s", selector, value)
	}

	return nil
}

// Click 点击元素
func (bm *BrowserManager) Click(selector string) error {
	if bm.Page == nil {
		return fmt.Errorf("页面未初始化")
	}

	if err := bm.WaitForSelector(selector, 10*time.Second); err != nil {
		return fmt.Errorf("等待点击元素 %s 失败: %w", selector, err)
	}

	if err := bm.Page.Click(selector); err != nil {
		return fmt.Errorf("点击元素 %s 失败: %w", selector, err)
	}

	log.Printf("✅ 已点击元素: %s", selector)
	return nil
}

// Screenshot 截取屏幕截图
func (bm *BrowserManager) Screenshot(filename string) error {
	if bm.Page == nil {
		return fmt.Errorf("页面未初始化")
	}

	_, err := bm.Page.Screenshot(playwright.PageScreenshotOptions{
		Path:     playwright.String(filename),
		FullPage: playwright.Bool(true),
	})
	if err != nil {
		return fmt.Errorf("截取屏幕截图失败: %w", err)
	}

	log.Printf("📸 已保存截图: %s", filename)
	return nil
}

// Close 关闭浏览器
func (bm *BrowserManager) Close() error {
	if bm.Context != nil {
		bm.Context.Close()
	}

	if bm.Browser != nil {
		return bm.Browser.Close()
	}

	return nil
}

package operator

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

// ScrollToElement 滚动到元素可见区域
func (bm *BrowserManager) ScrollToElement(selector string) error {
	if bm.Page == nil {
		return fmt.Errorf("页面未初始化")
	}

	if err := bm.WaitForSelector(selector, 10*time.Second); err != nil {
		return fmt.Errorf("等待元素 %s 失败: %w", selector, err)
	}

	// 滚动到元素位置
	_, err := bm.Page.EvalOnSelector(selector, "element => element.scrollIntoView({behavior: 'smooth', block: 'center'})", nil)
	if err != nil {
		return fmt.Errorf("滚动到元素 %s 失败: %w", selector, err)
	}

	log.Printf("📜 已滚动到元素: %s", selector)
	return nil
}

// Hover 鼠标悬停在元素上
func (bm *BrowserManager) Hover(selector string) error {
	if bm.Page == nil {
		return fmt.Errorf("页面未初始化")
	}

	if err := bm.WaitForSelector(selector, 10*time.Second); err != nil {
		return fmt.Errorf("等待悬停元素 %s 失败: %w", selector, err)
	}

	if err := bm.Page.Hover(selector); err != nil {
		return fmt.Errorf("悬停元素 %s 失败: %w", selector, err)
	}

	log.Printf("🎯 已悬停在元素: %s", selector)
	return nil
}

// SelectOption 从下拉菜单中选择选项
func (bm *BrowserManager) SelectOption(selector, value string) error {
	if bm.Page == nil {
		return fmt.Errorf("页面未初始化")
	}

	if err := bm.WaitForSelector(selector, 10*time.Second); err != nil {
		return fmt.Errorf("等待选择器元素 %s 失败: %w", selector, err)
	}

	// 使用SelectOptionValues类型
	_, err := bm.Page.Locator(selector).SelectOption(playwright.SelectOptionValues{
		Labels: playwright.StringSlice(value),
	})
	if err != nil {
		return fmt.Errorf("选择选项 %s 失败: %w", value, err)
	}

	log.Printf("📋 已选择选项: %s = %s", selector, value)
	return nil
}

// GetText 获取元素的文本内容
func (bm *BrowserManager) GetText(selector string) (string, error) {
	if bm.Page == nil {
		return "", fmt.Errorf("页面未初始化")
	}

	if err := bm.WaitForSelector(selector, 10*time.Second); err != nil {
		return "", fmt.Errorf("等待元素 %s 失败: %w", selector, err)
	}

	text, err := bm.Page.TextContent(selector)
	if err != nil {
		return "", fmt.Errorf("获取元素 %s 文本失败: %w", selector, err)
	}

	return text, nil
}

// GetAttribute 获取元素的属性值
func (bm *BrowserManager) GetAttribute(selector, attribute string) (string, error) {
	if bm.Page == nil {
		return "", fmt.Errorf("页面未初始化")
	}

	if err := bm.WaitForSelector(selector, 10*time.Second); err != nil {
		return "", fmt.Errorf("等待元素 %s 失败: %w", selector, err)
	}

	attr, err := bm.Page.GetAttribute(selector, attribute)
	if err != nil {
		return "", fmt.Errorf("获取元素 %s 属性 %s 失败: %w", selector, attribute, err)
	}

	return attr, nil
}

// IsVisible 检查元素是否可见
func (bm *BrowserManager) IsVisible(selector string) (bool, error) {
	if bm.Page == nil {
		return false, fmt.Errorf("页面未初始化")
	}

	visible, err := bm.Page.IsVisible(selector)
	if err != nil {
		return false, fmt.Errorf("检查元素 %s 可见性失败: %w", selector, err)
	}

	return visible, nil
}

// WaitForElementDisappear 等待元素消失
func (bm *BrowserManager) WaitForElementDisappear(selector string, timeout time.Duration) error {
	if bm.Page == nil {
		return fmt.Errorf("页面未初始化")
	}

	_, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	_, err := bm.Page.WaitForSelector(selector, playwright.PageWaitForSelectorOptions{
		State:   playwright.WaitForSelectorStateHidden,
		Timeout: playwright.Float(float64(timeout.Milliseconds())),
	})
	if err != nil {
		return fmt.Errorf("等待元素 %s 消失超时: %w", selector, err)
	}

	log.Printf("⏱️ 元素已消失: %s", selector)
	return nil
}

// DragAndDrop 拖拽元素到另一个位置
func (bm *BrowserManager) DragAndDrop(sourceSelector, targetSelector string) error {
	if bm.Page == nil {
		return fmt.Errorf("页面未初始化")
	}

	// 等待源元素出现
	if err := bm.WaitForSelector(sourceSelector, 10*time.Second); err != nil {
		return fmt.Errorf("等待拖拽源元素 %s 失败: %w", sourceSelector, err)
	}

	// 等待目标元素出现
	if err := bm.WaitForSelector(targetSelector, 10*time.Second); err != nil {
		return fmt.Errorf("等待拖拽目标元素 %s 失败: %w", targetSelector, err)
	}

	// 执行拖拽
	if err := bm.Page.DragAndDrop(sourceSelector, targetSelector); err != nil {
		return fmt.Errorf("拖拽元素 %s 到 %s 失败: %w", sourceSelector, targetSelector, err)
	}

	log.Printf("🔄 已拖拽元素: %s -> %s", sourceSelector, targetSelector)
	return nil
}

// RightClick 右键点击元素
func (bm *BrowserManager) RightClick(selector string) error {
	if bm.Page == nil {
		return fmt.Errorf("页面未初始化")
	}

	if err := bm.WaitForSelector(selector, 10*time.Second); err != nil {
		return fmt.Errorf("等待右键元素 %s 失败: %w", selector, err)
	}

	if err := bm.Page.Click(selector, playwright.PageClickOptions{
		Button: playwright.MouseButtonRight,
	}); err != nil {
		return fmt.Errorf("右键点击元素 %s 失败: %w", selector, err)
	}

	log.Printf("🖱️ 已右键点击元素: %s", selector)
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

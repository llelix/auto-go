package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// 设置默认值
	port := 8080
	templatesDir := "templates"
	
	// 检查模板目录是否存在
	if _, err := os.Stat(templatesDir); os.IsNotExist(err) {
		fmt.Printf("错误: 模板目录不存在: %s\n", templatesDir)
		os.Exit(1)
	}
	
	// 设置 Gin 模式
	gin.SetMode(gin.ReleaseMode)
	
	// 创建 Gin 路由器
	r := gin.Default()
	
	// 设置模板目录
	r.LoadHTMLGlob(templatesDir + "/*.html")
	
	// 静态文件服务
	r.Static("/static", "./static")
	
	// 首页路由
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Auto-Go Mock Server",
		})
	})
	
	// 简单表单页面路由
	r.GET("/simple-form", func(c *gin.Context) {
		c.HTML(http.StatusOK, "simple_form.html", gin.H{
			"title": "简单表单测试 - Auto-Go Mock Server",
		})
	})
	
	// 复杂表单页面路由
	r.GET("/complex-form", func(c *gin.Context) {
		c.HTML(http.StatusOK, "complex_form.html", gin.H{
			"title": "复杂表单测试 - Auto-Go Mock Server",
		})
	})
	
	// 拖拽功能页面路由
	r.GET("/drag-and-drop", func(c *gin.Context) {
		c.HTML(http.StatusOK, "drag_drop.html", gin.H{
			"title": "拖拽功能测试 - Auto-Go Mock Server",
		})
	})
	
	// 信息获取页面路由
	r.GET("/info-page", func(c *gin.Context) {
		c.HTML(http.StatusOK, "info_page.html", gin.H{
			"title": "信息获取测试 - Auto-Go Mock Server",
		})
	})
	
	// 交互测试页面路由
	r.GET("/interactive-test", func(c *gin.Context) {
		c.HTML(http.StatusOK, "interactive_test.html", gin.H{
			"title": "交互功能测试 - Auto-Go Mock Server",
		})
	})
	
	// 表单成功页面路由
	r.GET("/form-success", func(c *gin.Context) {
		c.HTML(http.StatusOK, "form_success.html", gin.H{
			"title": "表单提交成功 - Auto-Go Mock Server",
		})
	})
	
	// API 路由：获取当前时间
	r.GET("/api/time", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"time": time.Now().Format(time.RFC3339),
		})
	})
	
	// API 路由：获取国家对应的省份/州
	r.GET("/api/states/:country", func(c *gin.Context) {
		country := c.Param("country")
		states := make([]string, 0)
		
		switch country {
		case "china":
			states = []string{"北京", "上海", "广东", "浙江", "江苏", "四川", "湖北", "湖南"}
		case "usa":
			states = []string{"加利福尼亚", "纽约", "德克萨斯", "佛罗里达", "伊利诺伊"}
		case "japan":
			states = []string{"东京", "大阪", "京都", "北海道", "福冈"}
		case "uk":
			states = []string{"英格兰", "苏格兰", "威尔士", "北爱尔兰"}
		}
		
		c.JSON(http.StatusOK, gin.H{
			"states": states,
		})
	})
	
	// 表单提交路由
	r.POST("/submit-complex-form", func(c *gin.Context) {
		// 模拟表单处理
		username := c.PostForm("username")
		email := c.PostForm("email")
		country := c.PostForm("country")
		state := c.PostForm("state")
		agreeTerms := c.PostForm("agree-terms")
		
		// 简单验证
		if username == "" || email == "" || country == "" || agreeTerms != "on" {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "请填写所有必填字段并同意服务条款",
			})
			return
		}
		
		// 返回成功响应
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "表单提交成功",
			"data": map[string]string{
				"username": username,
				"email":    email,
				"country":  country,
				"state":    state,
			},
		})
	})
	
	// 简单表单提交路由
	r.POST("/submit-form", func(c *gin.Context) {
		// 模拟表单处理
		name := c.PostForm("name")
		email := c.PostForm("email")
		phone := c.PostForm("phone")
		message := c.PostForm("message")
		agreeTerms := c.PostForm("agree-terms")
		
		// 简单验证
		if name == "" || email == "" || agreeTerms != "on" {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "请填写必填字段并同意服务条款",
			})
			return
		}
		
		// 记录收到的数据
		log.Printf("收到表单提交: 姓名=%s, 邮箱=%s, 电话=%s, 消息=%s", name, email, phone, message)
		
		// 返回成功响应，重定向到表单成功页面
		c.Redirect(http.StatusFound, "/form-success")
	})
	
	// 启动服务器
	serverAddr := fmt.Sprintf(":%d", port)
	
	// 获取当前工作目录
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("错误: 获取当前工作目录失败: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Printf("Mock Server 已启动\n")
	fmt.Printf("服务器地址: http://localhost:%d\n", port)
	fmt.Printf("模板目录: %s\n", filepath.Join(wd, templatesDir))
	fmt.Printf("可用页面:\n")
	fmt.Printf("  - 首页: http://localhost:%d/\n", port)
	fmt.Printf("  - 简单表单: http://localhost:%d/simple-form\n", port)
	fmt.Printf("  - 复杂表单: http://localhost:%d/complex-form\n", port)
	fmt.Printf("  - 拖拽功能: http://localhost:%d/drag-and-drop\n", port)
	fmt.Printf("  - 信息获取: http://localhost:%d/info-page\n", port)
	fmt.Printf("  - 交互测试: http://localhost:%d/interactive-test\n", port)
	fmt.Printf("  - 表单成功: http://localhost:%d/form-success\n", port)
	fmt.Printf("\n按 Ctrl+C 停止服务器\n")
	
	// 启动 HTTP 服务器
	log.Fatal(r.Run(serverAddr))
}
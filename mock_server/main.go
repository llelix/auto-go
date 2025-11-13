package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	// 创建Gin路由器
	r := gin.Default()

	// 加载HTML模板
	r.LoadHTMLGlob("templates/*")

	// 静态文件服务
	r.Static("/static", "./static")

	// 首页路由
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Auto-Go Mock Server",
		})
	})

	// 简单表单页面
	r.GET("/simple-form", func(c *gin.Context) {
		c.HTML(http.StatusOK, "simple_form.html", gin.H{
			"title": "简单表单测试",
		})
	})

	// 复杂表单页面
	r.GET("/complex-form", func(c *gin.Context) {
		c.HTML(http.StatusOK, "complex_form.html", gin.H{
			"title": "复杂表单测试",
		})
	})

	// 拖拽测试页面
	r.GET("/drag-and-drop", func(c *gin.Context) {
		c.HTML(http.StatusOK, "drag_drop.html", gin.H{
			"title": "拖拽功能测试",
		})
	})

	// 信息获取测试页面
	r.GET("/info-page", func(c *gin.Context) {
		c.HTML(http.StatusOK, "info_page.html", gin.H{
			"title": "信息获取测试",
		})
	})

	// 交互测试页面
	r.GET("/interactive-test", func(c *gin.Context) {
		c.HTML(http.StatusOK, "interactive_test.html", gin.H{
			"title": "交互功能测试",
		})
	})

	// 处理表单提交
	r.POST("/submit-form", func(c *gin.Context) {
		name := c.PostForm("name")
		email := c.PostForm("email")
		phone := c.PostForm("phone")
		message := c.PostForm("message")
		agreeTerms := c.PostForm("agree-terms")
		
		// 模拟处理时间
		// time.Sleep(1 * time.Second)
		
		c.HTML(http.StatusOK, "form_success.html", gin.H{
			"title":  "提交成功",
			"name":   name,
			"email":  email,
			"phone":  phone,
			"message": message,
			"agree":  agreeTerms == "on",
		})
	})

	// 处理复杂表单提交
	r.POST("/submit-complex-form", func(c *gin.Context) {
		username := c.PostForm("username")
		email := c.PostForm("email")
		country := c.PostForm("country")
		state := c.PostForm("state")
		acceptCookies := c.PostForm("accept-cookies")
		
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data": gin.H{
				"username":       username,
				"email":          email,
				"country":        country,
				"state":          state,
				"acceptCookies":   acceptCookies == "on",
			},
			"message": "表单提交成功",
		})
	})

	// 模拟API接口，用于获取国家列表
	r.GET("/api/countries", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"countries": []gin.H{
				{"id": "china", "name": "China", "states": []string{"Beijing", "Shanghai", "Guangzhou", "Shenzhen"}},
				{"id": "usa", "name": "USA", "states": []string{"California", "New York", "Texas", "Florida"}},
				{"id": "japan", "name": "Japan", "states": []string{"Tokyo", "Osaka", "Kyoto", "Hokkaido"}},
				{"id": "uk", "name": "United Kingdom", "states": []string{"England", "Scotland", "Wales", "Northern Ireland"}},
			},
		})
	})

	// 模拟API接口，用于获取状态
	r.GET("/api/states/:country", func(c *gin.Context) {
		country := c.Param("country")
		
		states := map[string][]string{
			"china": {"Beijing", "Shanghai", "Guangzhou", "Shenzhen"},
			"usa":   {"California", "New York", "Texas", "Florida"},
			"japan": {"Tokyo", "Osaka", "Kyoto", "Hokkaido"},
			"uk":    {"England", "Scotland", "Wales", "Northern Ireland"},
		}
		
		if s, exists := states[country]; exists {
			c.JSON(http.StatusOK, gin.H{
				"states": s,
			})
		} else {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Country not found",
			})
		}
	})

	// 模拟拖拽成功接口
	r.POST("/drag-drop-success", func(c *gin.Context) {
		source := c.PostForm("source")
		target := c.PostForm("target")
		
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"source":  source,
			"target":  target,
			"message": "拖拽操作成功",
		})
	})

	// 模拟进度条接口
	r.GET("/api/progress", func(c *gin.Context) {
		progress := c.Query("value")
		if progress == "" {
			progress = "0"
		}
		
		progressInt, _ := strconv.Atoi(progress)
		if progressInt >= 100 {
			c.JSON(http.StatusOK, gin.H{
				"progress": 100,
				"status":   "completed",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"progress": progressInt + 10,
				"status":   "in-progress",
			})
		}
	})

	// 启动服务器
	fmt.Println("🚀 Mock服务器启动成功!")
	fmt.Println("📝 访问地址: http://localhost:8080")
	fmt.Println("📄 简单表单: http://localhost:8080/simple-form")
	fmt.Println("📄 复杂表单: http://localhost:8080/complex-form")
	fmt.Println("🎯 拖拽测试: http://localhost:8080/drag-and-drop")
	fmt.Println("ℹ️  信息获取: http://localhost:8080/info-page")
	fmt.Println("🎮 交互测试: http://localhost:8080/interactive-test")
	
	r.Run(":8080")
}
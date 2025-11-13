package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config 定义应用程序配置
type Config struct {
	Browser BrowserConfig `mapstructure:"browser" json:"browser"`
	Tasks   TasksConfig  `mapstructure:"tasks" json:"tasks"`
	Logging LoggingConfig `mapstructure:"logging" json:"logging"`
}

// BrowserConfig 浏览器配置
type BrowserConfig struct {
	Headless  bool   `mapstructure:"headless" json:"headless"`
	UserAgent string `mapstructure:"user_agent" json:"user_agent"`
	Timeout   int    `mapstructure:"timeout" json:"timeout"`
}

// TasksConfig 任务配置
type TasksConfig struct {
	DefaultWaitTime int  `mapstructure:"default_wait_time" json:"default_wait_time"`
	AutoScreenshot  bool `mapstructure:"auto_screenshot" json:"auto_screenshot"`
}

// LoggingConfig 日志配置
type LoggingConfig struct {
	Level   string `mapstructure:"level" json:"level"`
	File    string `mapstructure:"file" json:"file"`
	Console bool   `mapstructure:"console" json:"console"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Browser: BrowserConfig{
			Headless:  true,
			UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
			Timeout:   30,
		},
		Tasks: TasksConfig{
			DefaultWaitTime: 5,
			AutoScreenshot:  true,
		},
		Logging: LoggingConfig{
			Level:   "info",
			File:    "logs/auto-go.log",
			Console: true,
		},
	}
}

// LoadConfig 加载配置文件
func LoadConfig(configPath string) (*Config, error) {
	config := DefaultConfig()

	// 设置配置文件路径
	if configPath != "" {
		viper.SetConfigFile(configPath)
	} else {
		// 默认配置文件路径
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("获取用户目录失败: %w", err)
		}
		
		configDir := filepath.Join(homeDir, ".auto-go")
		viper.AddConfigPath(configDir)
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		// 如果配置文件不存在，使用默认配置
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("⚠️  配置文件未找到，使用默认配置")
			return config, nil
		}
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 解析配置到结构体
	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	fmt.Printf("✅ 配置文件加载成功: %s\n", viper.ConfigFileUsed())
	return config, nil
}

// SaveConfig 保存配置到文件
func SaveConfig(config *Config, configPath string) error {
	// 确保目录存在
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建配置目录失败: %w", err)
	}

	// 将配置写入文件
	file, err := os.Create(configPath)
	if err != nil {
		return fmt.Errorf("创建配置文件失败: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	
	if err := encoder.Encode(config); err != nil {
		return fmt.Errorf("编码配置文件失败: %w", err)
	}

	fmt.Printf("✅ 配置文件已保存: %s\n", configPath)
	return nil
}

// CreateDefaultConfigFile 创建默认配置文件
func CreateDefaultConfigFile(configPath string) error {
	config := DefaultConfig()
	return SaveConfig(config, configPath)
}

// ValidateConfig 验证配置参数
func ValidateConfig(config *Config) error {
	if config.Browser.Timeout <= 0 {
		return fmt.Errorf("浏览器超时时间必须大于0")
	}
	
	if config.Tasks.DefaultWaitTime < 0 {
		return fmt.Errorf("默认等待时间不能为负数")
	}
	
	return nil
}
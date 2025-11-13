package cmd

import (
	"github.com/urfave/cli/v2"
)

// 应用程序信息
const (
	AppName  = "auto-go"
	AppUsage = "基于Playwright的浏览器自动化工具"
	AppVer   = "1.0.0"
)

// CreateApp 创建CLI应用程序
func CreateApp() *cli.App {
	return &cli.App{
		Name:    AppName,
		Usage:   AppUsage,
		Version: AppVer,
		Commands: []*cli.Command{
			CreateRunCommand(),
			CreateInitCommand(),
		},
	}
}
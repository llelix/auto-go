@echo off
echo 启动Mock服务器...
echo.
echo 请确保已安装Gin框架依赖:
echo go mod download
echo.
echo 如果未安装，请先运行:
echo go get -u github.com/gin-gonic/gin
echo.
go run main.go
pause
package main

import (
	"embed"
	"log"

	"github.com/li-sky/NebulaRail/config" // 使用 go.mod 中的模块路径
	"github.com/li-sky/NebulaRail/db"     // 使用 go.mod 中的模块路径
	"github.com/li-sky/NebulaRail/server" // 使用 go.mod 中的模块路径
)

//go:embed frontend/dist
var embeddedFiles embed.FS

func main() {
	// 初始化数据库
	database, err := db.InitDB()
	if err != nil {
		log.Fatalf("无法初始化数据库: %v", err)
	}
	defer database.Close() // 确保数据库连接被关闭

	// 初始化配置
	appConfig, err := config.InitConfig(database)
	if err != nil {
		log.Fatalf("无法初始化配置: %v", err)
	}

	// 启动服务器，传入嵌入的文件系统
	server.Run(appConfig, embeddedFiles)
}

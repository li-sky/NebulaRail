package server

import (
	"embed" // 添加 embed 包导入
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/li-sky/NebulaRail/config"     // 使用 go.mod 中的模块路径
	"github.com/li-sky/NebulaRail/embedfs"    // 使用 go.mod 中的模块路径
	"github.com/li-sky/NebulaRail/handlers"   // 使用 go.mod 中的模块路径
	"github.com/li-sky/NebulaRail/middleware" // 使用 go.mod 中的模块路径
)

// Run 启动服务器，接收 embed.FS
func Run(appConfig config.Config, embeddedFiles embed.FS) {
	r := gin.Default()

	// 获取前端文件系统，传入 embed.FS
	httpFS, err := embedfs.GetHTTPFS(embeddedFiles)
	if err != nil {
		log.Fatalf("无法获取嵌入式文件系统: %v", err)
	}

	// API 路由组
	api := r.Group("/api")
	{
		// 公开 API
		api.POST("/login", handlers.LoginHandler(&appConfig)) // 传递指针以便修改密码时更新内存

		// 受保护 API
		apiAuth := api.Group("/")
		apiAuth.Use(middleware.AuthMiddleware(appConfig))
		{
			apiAuth.POST("/change-password", handlers.ChangePasswordHandler(&appConfig)) // 传递指针
			apiAuth.GET("/protected", handlers.ProtectedHandler())
			// 在这里添加其他需要认证的API
		}
	}

	// 前端静态文件路由 和 404 处理
	r.NoRoute(func(c *gin.Context) {
		// 如果是API路径，返回404
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.JSON(http.StatusNotFound, gin.H{"error": "API路径不存在"})
			return
		}

		// 否则，尝试提供前端文件
		filePath := strings.TrimPrefix(c.Request.URL.Path, "/")
		if filePath == "" {
			filePath = "index.html" // 根路径映射到 index.html
		}

		// 尝试打开文件
		f, err := httpFS.Open(filePath)
		if err != nil {
			// 如果文件不存在 (包括目录)，则返回 index.html (SPA 路由)
			if os.IsNotExist(err) {
				filePath = "index.html"
			} else {
				// 其他错误，例如权限问题
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "无法提供文件"})
				return
			}
		} else {
			f.Close() // 如果文件存在，关闭它，让 FileServer 处理
		}

		// 如果需要返回 index.html，重写请求路径
		if filePath == "index.html" {
			c.Request.URL.Path = "/" // FileServer 会查找根目录下的 index.html
		}

		// 使用 http.FileServer 处理请求
		http.FileServer(httpFS).ServeHTTP(c.Writer, c.Request)
	})

	// 设置监听端口
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("服务器正在监听端口: %s", port)

	// 启动服务器
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("无法启动服务器: %v", err)
	}
}

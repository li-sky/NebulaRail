package embedfs

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
)

// GetSubFS 获取嵌入式文件系统的子文件系统
func GetSubFS(embeddedFiles embed.FS) (fs.FS, error) {
	// 路径相对于 embed 指令指定的根目录
	subFS, err := fs.Sub(embeddedFiles, "frontend/dist")
	if err != nil {
		log.Printf("无法获取嵌入式文件系统的子文件系统: %v", err)
		return nil, err
	}
	return subFS, nil
}

// GetHTTPFS 获取用于HTTP服务的嵌入式文件系统
func GetHTTPFS(embeddedFiles embed.FS) (http.FileSystem, error) {
	subFS, err := GetSubFS(embeddedFiles)
	if err != nil {
		return nil, err
	}
	return http.FS(subFS), nil
}

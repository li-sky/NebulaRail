package db

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3" // Import the SQLite driver
)

// InitDB 初始化数据库
func InitDB() (*sql.DB, error) {
	// 确保数据目录存在
	dataDir := "./data"
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		if err := os.MkdirAll(dataDir, 0755); err != nil {
			return nil, err
		}
	}

	// 打开或创建SQLite数据库
	dbPath := filepath.Join(dataDir, "nebulaRail.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// 创建配置表
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS config (
		key TEXT PRIMARY KEY,
		value TEXT NOT NULL
	)`)
	if err != nil {
		db.Close() // Close the DB if table creation fails
		return nil, err
	}

	log.Println("数据库初始化成功:", dbPath)
	return db, nil
}

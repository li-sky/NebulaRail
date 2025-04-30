package config

import (
	"crypto/rand" // 导入 crypto/rand
	"database/sql"
	"encoding/base64" // 导入 encoding/base64
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Config 应用程序配置
type Config struct {
	PasswordHash string        // bcrypt哈希后的密码
	JWTSecret    string        // JWT签名密钥
	JWTExpiry    time.Duration // JWT过期时间
	DB           *sql.DB       // 数据库连接
}

// generateRandomKey 生成安全的随机密钥
func generateRandomKey(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	// 使用 URL 安全的 Base64 编码
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// InitConfig 初始化配置
func InitConfig(db *sql.DB) (Config, error) {
	// 尝试从数据库获取密码哈希，如果不存在则创建默认值
	defaultPasswordHash, _ := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	passwordHash := GetConfigValue(db, "password_hash", string(defaultPasswordHash))

	// 获取或生成并存储 JWT 密钥 (推荐 32 或 64 字节)
	jwtSecret := GetOrGenerateConfigValue(db, "jwt_secret", 32)
	if jwtSecret == "" {
		// 如果 GetOrGenerateConfigValue 返回空字符串，表示生成或获取时发生严重错误
		log.Fatal("无法获取或生成 JWT 密钥")
	}

	return Config{
		PasswordHash: passwordHash,
		JWTSecret:    jwtSecret,
		JWTExpiry:    24 * time.Hour,
		DB:           db,
	}, nil
}

// GetOrGenerateConfigValue 获取配置值，如果不存在则生成随机值并存储
func GetOrGenerateConfigValue(db *sql.DB, key string, keyLengthBytes int) string {
	var value string
	err := db.QueryRow("SELECT value FROM config WHERE key = ?", key).Scan(&value)
	if err != nil {
		if err == sql.ErrNoRows {
			// 未找到，生成新密钥
			newValue, err := generateRandomKey(keyLengthBytes)
			if err != nil {
				log.Printf("严重错误: 无法生成配置项 %s 的随机值: %v", key, err)
				return "" // 返回空字符串表示失败
			}
			log.Printf("为 '%s' 生成了新的随机值。请妥善保管数据库文件，密钥存储在其中。", key)

			// 存储新密钥
			_, err = db.Exec("INSERT INTO config (key, value) VALUES (?, ?)", key, newValue)
			if err != nil {
				log.Printf("警告: 无法将新生成的配置项 '%s' 保存到数据库: %v", key, err)
				// 即使保存失败，也返回生成的密钥供本次运行使用，但下次启动会重新生成
			}
			return newValue
		}
		// 其他查询错误
		log.Printf("严重错误: 查询配置项 '%s' 时出错: %v", key, err)
		return "" // 返回空字符串表示失败
	}
	// 找到了现有值
	return value
}

// GetConfigValue 从数据库获取配置值 (保持不变，用于 password_hash)
func GetConfigValue(db *sql.DB, key, defaultValue string) string {
	var value string
	err := db.QueryRow("SELECT value FROM config WHERE key = ?", key).Scan(&value)
	if err != nil {
		// 如果没找到或出错，使用默认值并保存到数据库
		if err == sql.ErrNoRows {
			_, err = db.Exec("INSERT INTO config (key, value) VALUES (?, ?)", key, defaultValue)
			if err != nil {
				log.Printf("无法保存配置项 %s: %v", key, err)
			}
		} else {
			log.Printf("查询配置项 %s 时出错: %v", key, err) // 记录其他错误
		}
		return defaultValue
	}
	return value
}

// SetConfigValue 更新配置值 (保持不变，用于 password_hash)
func SetConfigValue(db *sql.DB, key, value string) error {
	_, err := db.Exec("INSERT OR REPLACE INTO config (key, value) VALUES (?, ?)", key, value)
	if err != nil {
		log.Printf("更新配置项 %s 时出错: %v", key, err)
	}
	return err
}

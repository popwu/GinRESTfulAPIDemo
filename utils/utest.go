package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	// 获取当前工作目录
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	// 向上查找 .env 文件
	for {
		envFile := filepath.Join(wd, ".env")
		if _, err := os.Stat(envFile); err == nil {
			return godotenv.Load(envFile)
		}

		parent := filepath.Dir(wd)
		if parent == wd {
			break
		}
		wd = parent
	}

	return fmt.Errorf(".env file not found")
}

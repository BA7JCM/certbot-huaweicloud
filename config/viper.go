package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func Init() {
	viper.SetConfigName(".env") // 配置文件名称
	viper.SetConfigType("env")  // 配置文件类型

	// 首先在可执行文件所在目录查找（避免 crontab/certbot hook 调用时 CWD 不匹配）
	if exePath, err := os.Executable(); err == nil {
		viper.AddConfigPath(filepath.Dir(exePath))
	}

	// 兜底：在当前工作目录下查找
	viper.AddConfigPath(".")

	err := viper.ReadInConfig() // 读取配置
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok { // 配置文件未找到
			panic("配置文件未找到")
		}
		panic("配置文件读取失败")
	}
}

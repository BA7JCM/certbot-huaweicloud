package config

import (
	"github.com/spf13/viper"
)

func Init() {
	viper.SetConfigName(".env") // 配置文件名称
	viper.SetConfigType("env")  // 配置文件类型
	viper.AddConfigPath("./")   // 在当前文件夹下寻找
	viper.AddConfigPath(".")    // 在工作目录下查找
	err := viper.ReadInConfig() // 读取配置
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok { // 配置文件未找到
			panic("配置文件未找到")
		}
		panic("配置文件读取失败")
	}
}

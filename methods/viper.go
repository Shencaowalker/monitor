package methods

import (
	"fmt"

	"github.com/spf13/viper"
)

func Viper_reader() *viper.Viper {
	config := viper.New()
	config.SetConfigName("metrics")
	// 设置配置文件的类型
	config.SetConfigType("yaml")
	// 添加配置文件的路径，指定 config 目录下寻找
	config.AddConfigPath("./conf")
	// 寻找配置文件并读取
	err := config.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	
	return config
}

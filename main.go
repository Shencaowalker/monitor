package main

import (
	"fmt"
	"log_as/apis"
	"net/http"

	"github.com/spf13/viper"
)

func main() {
	//初始化配置句柄
	config := viper.New()
	config.AddConfigPath("./conf/") // 文件所在目录
	config.SetConfigName("config")  // 文件名
	config.SetConfigType("ini")     // 文件类型

	if err := config.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("找不到配置文件..")
		} else {
			fmt.Println("配置文件出错..")
		}
	}

	//调用接口
	http.HandleFunc("/registered", apis.Registered(config))
	http.HandleFunc("/downline", apis.Downline(config))
	http.HandleFunc("/updatenacosstandardconf", apis.UpdateNacosStandardConf(&config))
	http.HandleFunc("/nacosproductupload", apis.NacosProductUpload())
	http.ListenAndServe("0.0.0.0:"+config.GetString("global.serviceport"), nil)

	//创建指标文件并函数退出后返回

	/*
		serviceStatus, err := os.Create("NducerStatus.txt")
		defer serviceStatus.Close()
		if err != nil {
			fmt.Println("文件创建失败", err)
		}
		serviceList := config.GetStringSlice("global.servicelist") //["dws","cip","das","afp","asp","bde","tse","arctic","jobserver"]
		for i := 0; i < len(serviceList); i++ {
			methods.Contrast(config, serviceList[i], serviceStatus)
		}
		methods.UpdateMetrics(serviceStatus.Name(), config)

	*/

}

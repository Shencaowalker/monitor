package apis

import (
	"encoding/json"
	"fmt"
	"log"
	"log_as/methods"
	"net/http"
	"os"

	"github.com/spf13/viper"
)

func NacosProductUpload(config *viper.Viper) func(writer http.ResponseWriter, request *http.Request) {
	serviceStatus, err := os.Create("NducerStatus.txt")
	defer serviceStatus.Close()
	if err != nil {
		fmt.Println("文件创建失败", err)
	}
	serviceList := config.GetStringSlice("global.servicelist")
	//["dws","cip","das","afp","asp","bde","tse","arctic","jobserver"]
	for i := 0; i < len(serviceList); i++ {
		methods.Contrast(config, serviceList[i], serviceStatus)
	}
	methods.UpdateMetrics(serviceStatus.Name(), config)

	return func(writer http.ResponseWriter, request *http.Request) {
		log.Println("NacosProductUpload,开始进行上传nacos 生产者指标")
		var result methods.Resp
		result.Msg = "调用成功，异步执行更新操作。"
		result.Code = "200"
		if err := json.NewEncoder(writer).Encode(result); err != nil {
			log.Fatal(err)
		}
	}
}

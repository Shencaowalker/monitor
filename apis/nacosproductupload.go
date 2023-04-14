package apis

import (
	"encoding/json"
	"fmt"
	"log"
	"log_as/methods"
	"net/http"
	"os"
)

func NacosProductUpload() func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		confignew := methods.CreateNewconfiger("./conf/", "config", "ini")
		serviceStatus, err := os.Create("NducerStatus.txt")
		defer serviceStatus.Close()
		if err != nil {
			fmt.Println("文件创建失败", err)
		}
		serviceList := confignew.GetStringSlice("global.servicelist")
		//["dws","cip","das","afp","asp","bde","tse","arctic","jobserver"]
		for i := 0; i < len(serviceList); i++ {
			methods.Contrast(confignew, serviceList[i], serviceStatus)
			fmt.Println(serviceList[i])
		}

		methods.UpdateMetrics(serviceStatus.Name(), confignew)
		fmt.Println("执行 UpdateMetrics 完毕")
		log.Println("NacosProductUpload,开始进行上传nacos 生产者指标")
		var result methods.Resp
		result.Msg = "上传nacos 生产者指标成功。"
		result.Code = "200"
		if err := json.NewEncoder(writer).Encode(result); err != nil {
			log.Fatal(err)
		}
	}
}

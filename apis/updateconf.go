package apis

import (
	"encoding/json"
	"log"
	"monitor/methods"
	"net/http"
	"os"

	"github.com/spf13/viper"
)

func UpdateNacosStandardConf(configaddr *(*viper.Viper)) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		log.Println("调用UpdateNacosStandardConf,开始更新配置流程")
		go methods.AsyncBUpdateNacosStandardConf(configaddr)
		var result methods.Resp
		result.Msg = "调用成功，异步执行更新操作。"
		result.Code = "200"
		if err := json.NewEncoder(writer).Encode(result); err != nil {
			log.Println(err)
		}
	}
}

func UpdateNacosproducerMonitor(config *viper.Viper) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		f, _ := os.Create("Status.txt")
		defer f.Close()
		log.Println("调用接口产生nacos生产者状态")
		servicelist := config.GetStringSlice("global.servicelist")
		for _, j := range servicelist {
			methods.Contrast(config, j, f)
		}
		err := methods.UpdateNacosMetrics(config, f.Name(), "serviceproducer")

		var result methods.Resp
		if err != nil {
			result.Msg = err.Error()
			result.Code = "400"
		} else {
			result.Msg = "调用成功。"
			result.Code = "200"
		}
		if err := json.NewEncoder(writer).Encode(result); err != nil {
			log.Println(err)
		}
	}
}

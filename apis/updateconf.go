package apis

import (
	"encoding/json"
	"log"
	"monitor/methods"
	"net/http"

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

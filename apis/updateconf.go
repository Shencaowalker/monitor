package apis

import (
	"encoding/json"
	"fmt"
	"log"
	"log_as/methods"
	"net/http"

	"github.com/spf13/viper"
)

func UpdateNacosStandardConf(config *viper.Viper) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		log.Println("调用UpdateNacosStandardConf,开始更新配置流程")
		fmt.Println()
		go methods.AsyncBUpdateNacosStandardConf(config)
		var result methods.Resp
		result.Msg = "调用成功，异步执行更新操作。"
		result.Code = "200"
		if err := json.NewEncoder(writer).Encode(result); err != nil {
			log.Fatal(err)
		}
	}
}

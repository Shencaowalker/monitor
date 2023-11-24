package apis

import (
	"encoding/json"
	"fmt"
	"log"
	"monitor/methods"
	"net/http"

	"github.com/spf13/viper"
)

func DingdingrobotSend(config *viper.Viper) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var alarmdata methods.Alarmdata
		if err := json.NewDecoder(request.Body).Decode(&alarmdata); err != nil {
			request.Body.Close()
			log.Println(err)
		}
		fmt.Println("此次告警数据为:", alarmdata)
		// result := methods.Sendmessage(config, alarmdata)
		// if err := json.NewEncoder(writer).Encode(result); err != nil {
		// 	log.Fatal(err)
		methods.Sendmessage(config, alarmdata)
	}
}

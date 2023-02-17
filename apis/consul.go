package apis

import (
	"encoding/json"
	"log"
	"log_as/methods"
	"net/http"

	"github.com/spf13/viper"
)

//post接口接收json数据 item
func RegisteredItem(config *viper.Viper) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var information methods.Registration_nformation
		if err := json.NewDecoder(request.Body).Decode(&information); err != nil {
			request.Body.Close()
			log.Fatal(err)
		}
		result := methods.ConsulregisterItem(config, information)
		if err := json.NewEncoder(writer).Encode(result); err != nil {
			log.Fatal(err)
		}
	}
}

//接收x-www-form-urlencoded类型的post请求或者普通get请求
func DownlineItems(config *viper.Viper) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		request.ParseForm()
		id, _ := request.Form["itemid"]
		result := methods.ConsuldownlineItems(config, id[0])
		if err := json.NewEncoder(writer).Encode(result); err != nil {
			log.Fatal(err)
		}
	}
}

//post接口接收json数据 alarm
func RegisteredAlarm(config *viper.Viper) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var information methods.Registration_Alarm

		// data, err := ioutil.ReadAll(request.Body)
		// defer request.Body.Close()
		// if err == nil && data != nil {

		// if err := json.Unmarshal(bodystr, &information); err != nil {
		// 	log.Fatal(err)
		// }

		// bodystr := mahonia.NewDecoder("gbk").NewReader(request.Body)
		if err := json.NewDecoder(request.Body).Decode(&information); err != nil {
			request.Body.Close()
			log.Fatal(err)
		}

		result := methods.ConsulregisterAlarm(config, information)
		if err := json.NewEncoder(writer).Encode(result); err != nil {
			log.Fatal(err)
		}
	}
}

func DownlineAlarm(config *viper.Viper) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		request.ParseForm()
		id, _ := request.Form["alarmid"]
		result := methods.ConsuldownlineAlarm(config, id[0])
		if err := json.NewEncoder(writer).Encode(result); err != nil {
			log.Fatal(err)
		}
	}
}

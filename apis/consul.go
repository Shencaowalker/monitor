package apis

import (
	"encoding/json"
	"log"
	"log_as/methods"
	"net/http"

	"github.com/spf13/viper"
)

//post接口接收json数据
func Registered(config *viper.Viper) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var information methods.Registration_nformation
		if err := json.NewDecoder(request.Body).Decode(&information); err != nil {
			request.Body.Close()
			log.Fatal(err)
		}
		result := methods.Consulregister(config, information)
		if err := json.NewEncoder(writer).Encode(result); err != nil {
			log.Fatal(err)
		}
	}
}

//接收x-www-form-urlencoded类型的post请求或者普通get请求
func Downline(config *viper.Viper) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		request.ParseForm()
		id, _ := request.Form["id"]
		result := methods.Consuldownline(config, id[0])
		if err := json.NewEncoder(writer).Encode(result); err != nil {
			log.Fatal(err)
		}
	}
}

package apis

import (
	"encoding/json"
	"log"
	"monitor/methods"
	"net/http"

	"github.com/spf13/viper"
)

// swagger:route POST /registereditem RegisteredItem RegisteredItem
//
// post接口接收json数据 item.
//
// This will show all available pets by default.
//
//     Schemes: http, https
//
//     Responses:
//       200: petsResponse
//       401: genericError
//       500: genericError
func RegisteredItem(config *viper.Viper) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var information methods.Registration_nformation
		if err := json.NewDecoder(request.Body).Decode(&information); err != nil {
			request.Body.Close()
			log.Println(err)
		}
		result := methods.ConsulregisterItem(config, information)
		if err := json.NewEncoder(writer).Encode(result); err != nil {
			log.Println(err)
		}
	}
}

func RegisteredItems(config *viper.Viper) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var informations methods.Registration_nformations
		if err := json.NewDecoder(request.Body).Decode(&informations); err != nil {
			request.Body.Close()
			log.Println(err)
		}
		for _, information := range informations.Values {
			result := methods.ConsulregisterItem(config, information)
			if err := json.NewEncoder(writer).Encode(result); err != nil {
				log.Println(err)
			}
		}
	}
}

//接收x-www-form-urlencoded类型的post请求或者普通get请求  例如：/downlineitem?itemid=100.100.100.100_air_100.100.100.100_9100&itemid=100.100.100.101_air_100.100.100.101_9100  会删除两个，多写会顺序删除多个
func DownlineItemsget(config *viper.Viper) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		request.ParseForm()
		ids, _ := request.Form["itemid"]
		for _, id := range ids {
			result := methods.ConsuldownlineItem(config, id)
			if err := json.NewEncoder(writer).Encode(result); err != nil {
				log.Println(err)
			}
		}
	}
}

// 接收post请求,json如
// {
// "itemids":["100.100.100.100_air_100.100.100.100_9100","100.100.100.101_air_100.100.100.101_9100"]
// }
func DownlineItemspost(config *viper.Viper) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var downlineitemids methods.Ddownline_nformations
		if err := json.NewDecoder(request.Body).Decode(&downlineitemids); err != nil {
			request.Body.Close()
			log.Println(err)
		}
		for _, itemid := range downlineitemids.Itemids {
			result := methods.ConsuldownlineItem(config, itemid)
			if err := json.NewEncoder(writer).Encode(result); err != nil {
				log.Println(err)
			}
		}
	}
}

//post接口接收json数据 注册监控项目。单条执行
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
			log.Println(err)
		}

		result := methods.ConsulregisterAlarm(config, information)
		if err := json.NewEncoder(writer).Encode(result); err != nil {
			log.Println(err)
		}
	}
}

// get接口接收下线告警项，单条执行
func DownlineAlarm(config *viper.Viper) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		request.ParseForm()
		id, _ := request.Form["alarmid"]
		result := methods.ConsuldownlineAlarm(config, id[0])
		if err := json.NewEncoder(writer).Encode(result); err != nil {
			log.Println(err)
		}
	}
}

package main

import (
	"log_as/apis"
	"log_as/methods"
	"net/http"
)

func main() {
	config := methods.CreateNewconfiger("./conf/", "config", "ini")
	//调用接口
	http.HandleFunc("/registereditem", apis.RegisteredItem(config))
	http.HandleFunc("/registeredalarm", apis.RegisteredAlarm(config))
	http.HandleFunc("/downlineitem", apis.DownlineItems(config))
	http.HandleFunc("/downlinealarm", apis.DownlineAlarm(config))
	http.HandleFunc("/updatenacosstandardconf", apis.UpdateNacosStandardConf(&config))
	http.ListenAndServe("0.0.0.0:"+config.GetString("global.serviceport"), nil)
}

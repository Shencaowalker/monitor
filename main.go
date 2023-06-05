package main

import (
	"flag"
	"fmt"
	"log_as/apis"
	"log_as/timecmd"
	"net/http"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	// mysqlIp     string
	// mysqlPort   string
	// mysqlUser   string
	// mysqlPasswd string
	// dbname      string
	// monitortype string
	istask  string
	project string
	operate string
	logtype string
)

func parseArgs() {
	// flag.StringVar(&mysqlIp, "H", "10.5.20.21", "-H 0.0.0.0~255.255.255.255")
	// flag.StringVar(&mysqlPort, "p", "3306", "-p port")
	// flag.StringVar(&mysqlUser, "u", "aloudata", "-u user")
	// flag.StringVar(&mysqlPasswd, "P", "MQvDPHzApoi.WGGPt12", "-P passwd")
	// flag.StringVar(&dbname, "d", "result_rp_dimension", "-d result_rp_dimension")
	// flag.StringVar(&monitortype, "t", "icebergnum", "-t icebergnum  or apppstatus")
	// flag.StringVar(&rangevalue, "c", "http://10.5.102.23:8889", "-c http://10.5.102.23:8889")
	flag.StringVar(&istask, "istask", "false", "-istask false;-istask true")
	flag.StringVar(&project, "project", "airerrorserverlogtopg", "-project airerrorserverlogtopg")
	flag.StringVar(&operate, "operate", "airserverlogtopg", "-operate airserverlogtopg;-operate airservertypelogtopg")
	flag.StringVar(&logtype, "logtype", "ERROR", "-logtype ERROR;-logtype INFO")
	flag.Parse()
}

func main() {
	parseArgs()
	//初始化配置句柄
	config := viper.New()
	config.AddConfigPath("./conf/") // 文件所在目录
	config.SetConfigName("config")  // 文件名
	config.SetConfigType("ini")     // 文件类型

	if err := config.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("找不到配置文件..")
		} else {
			fmt.Println("配置文件出错..")
		}
	}
	config.WatchConfig()
	config.OnConfigChange(func(e fsnotify.Event) {
		// 配置文件发生变更之后会调用的回调函数
		fmt.Println("Config file changed:", e.Name)
	})
	if istask == "true" {
		switch operate {
		case "airserverlogtopg":
			timecmd.Serverlogtopg(config, project)
		case "airservertypelogtopg":
			timecmd.Servertypelogetopg(config, logtype, project)
		default:
			fmt.Println("Please run main.go -c XXXX; XXXX scope please refer `main.go -h` ")
		}
	} else if istask == "false" {
		http.HandleFunc("/registereditem", apis.RegisteredItem(config))
		http.HandleFunc("/registeredalarm", apis.RegisteredAlarm(config))
		http.HandleFunc("/downlineitem", apis.DownlineItems(config))
		http.HandleFunc("/downlinealarm", apis.DownlineAlarm(config))
		http.HandleFunc("/updatenacosstandardconf", apis.UpdateNacosStandardConf(&config))
		http.ListenAndServe("0.0.0.0:"+config.GetString("global.serviceport"), nil)
	} else {
		fmt.Println("Please run main.go -c XXXX; XXXX scope please refer `main.go -h` ")
	}
}

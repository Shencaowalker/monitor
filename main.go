package main

import (
	"fmt"
	"log_as/methods"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	mysqlIp     string
	mysqlPort   string
	mysqlUser   string
	mysqlPasswd string
	dbname      string
	monitortype string
	rangevalue  string
)

// func parseArgs() {
// 	flag.StringVar(&mysqlIp, "H", "10.5.20.21", "-H 0.0.0.0~255.255.255.255")
// 	flag.StringVar(&mysqlPort, "p", "3306", "-p port")
// 	flag.StringVar(&mysqlUser, "u", "aloudata", "-u user")
// 	flag.StringVar(&mysqlPasswd, "P", "MQvDPHzApoi.WGGPt12", "-P passwd")
// 	flag.StringVar(&dbname, "d", "result_rp_dimension", "-d result_rp_dimension")
// 	flag.StringVar(&monitortype, "t", "icebergnum", "-t icebergnum  or apppstatus")
// 	flag.StringVar(&rangevalue, "c", "http://10.5.102.23:8889", "-c http://10.5.102.23:8889")
// 	flag.Parse()
// }

func main() {
	// parseArgs()
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
	//创建指标上传文件
	// serviceStatus, err := os.Create("Status.txt")
	// defer serviceStatus.Close()
	// if err != nil {
	// 	fmt.Println("文件创建失败", err)
	// }

	//入库跟指标上传使用一次loki请求
	// exp1 := config.GetStringMap("logprocessed.Listmetrics")
	// var values [][]interface{}
	// var metricsstring string
	// timenow := time.Now()
	// for _, jj := range exp1 {
	// 	loglist := methods.Getlogfromloki(config.GetString("logprocessed.lokiipport"), (jj.(map[string]interface{})["label_list"]).(map[string]interface{}), timenow, jj.(map[string]interface{})["latencycollectionseconds"].(string), jj.(map[string]interface{})["collectionscopeseconds"].(string))
	// 	for _, j := range loglist {
	// 		raw := methods.SplitoneLinetopostgresql(j, "|")
	// 		values = append(values, raw)
	// 		// onmetric := methods.SplitoneJoinsightLinetometrics(ii, (exp1[ii].(map[string]interface{})["label_name"]).([]interface{}), (exp1[ii].(map[string]interface{})["values"]).([]interface{}), j, "|")
	// 		// metricsstring = metricsstring + onmetric
	// 	}
	// 	// _, err := serviceStatus.Write([]byte(metricsstring))
	// 	// if err != nil {
	// 	// 	fmt.Println("写入" + ii + "指标临时文件失败，退出")
	// 	// 	continue
	// }

	// re := regexp.MustCompile(`(\d+-\d+-\d+\s\S+)\s\[\S+]\s(\w+)\s+(\S+)\s-\s(.*)`)

	exp2 := config.GetStringMap("logprocessed.listmetrics")
	fmt.Println("exp2 is:", exp2)
	var values [][]interface{}
	timenow := time.Now()
	for ii, jj := range exp2 {
		loglists := methods.Getlogfromloki(config.GetString("logprocessed.lokiipport"), (jj.(map[string]interface{})["label_list"]).(map[string]interface{}), timenow, jj.(map[string]interface{})["latencycollectionseconds"].(string), jj.(map[string]interface{})["collectionscopeseconds"].(string))
		for _, j := range loglists {
			fmt.Println(j)
			raw := methods.SplitoneLineforreFiltertopostgresql(j, jj.(map[string]interface{})["restring"].(string))
			if len(raw) == 0 {
				fmt.Println("匹配失败")
			} else {
				raw = append(raw, ii)
				values = append(values, raw)
			}
		}
	}
	fmt.Println(values)
	// db := methods.InitDB(config,"airlogtopostgresql")
	// methods.InsertintoDB(db, config, "airlogtopostgresql",values)
}

// func logtoDB(config *viper.Viper){
// 	var values [][]interface{}
// 	exp1 := config.GetStringMap("logprocessed.Listmetrics")
// 	for _, jj := range exp1 {
// 		// fmt.Println(i, "   ", j)
// 		loglist := methods.Getlogfromloki(config.GetString("logprocessed.lokiipport"), (jj.(map[string]interface{})["label_list"]).(map[string]interface{}), jj.(map[string]interface{})["latencycollectionseconds"].(string), jj.(map[string]interface{})["collectionscopeseconds"].(string))
// 		for _, j := range loglist {
// 			raw := methods.SplitoneLinetopostgresql(j, "|")
// 			values = append(values, raw)
// 		}
// 	}
// 	db := methods.InitDB(config)
// 	methods.InsertintoDB(db, config, values)
// }

// func metricPush(config *viper.Viper){

// 	for ii, jj := range exp1 {
// 		// a := methods.SplitoneJoinsightLinetometrics("joinsight_metrics", (exp1["joinsight_metrics"].(map[string]interface{})["label_name"]).([]interface{}), "cost", "2023-03-27 11:35:08.252|TID: acd254e44ed7411392b9009a5e21f1ce.189.16798881079585459|tn_40428|292960732771778560|QUERY|{\"cost\":254,\"queryStage\":\"SQL_TRANSLATE\",\"success\":true}", "|")
// 		loglist := methods.Getlogfromloki(config.GetString("logprocessed.lokiipport"), (jj.(map[string]interface{})["label_list"]).(map[string]interface{}), jj.(map[string]interface{})["latencycollectionseconds"].(string), jj.(map[string]interface{})["collectionscopeseconds"].(string))
// 		a := methods.SplitoneJoinsightLinetometrics(ii, (exp1[ii].(map[string]interface{})["label_name"]).([]interface{}), (exp1[ii].(map[string]interface{})["values"]).([]interface{}), , "|")
// 		fmt.Println(a)
// }
// }

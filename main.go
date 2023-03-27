package main

import (
	"fmt"
	"log_as/methods"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func main() {
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
	exp1 := config.GetStringMap("logprocessed.Listmetrics")
	for ii, _ := range exp1 {
		// a := methods.SplitoneJoinsightLinetometrics("joinsight_metrics", (exp1["joinsight_metrics"].(map[string]interface{})["label_name"]).([]interface{}), "cost", "2023-03-27 11:35:08.252|TID: acd254e44ed7411392b9009a5e21f1ce.189.16798881079585459|tn_40428|292960732771778560|QUERY|{\"cost\":254,\"queryStage\":\"SQL_TRANSLATE\",\"success\":true}", "|")
		a := methods.SplitoneJoinsightLinetometrics(ii, (exp1[ii].(map[string]interface{})["label_name"]).([]interface{}), (exp1[ii].(map[string]interface{})["values"]).([]interface{}), "2023-03-27 11:35:08.252|TID: acd254e44ed7411392b9009a5e21f1ce.189.16798881079585459|tn_40428|292960732771778560|QUERY|{\"cost\":254,\"queryStage\":\"SQL_TRANSLATE\",\"success\":true}", "|")
		fmt.Println(a)
	}
}

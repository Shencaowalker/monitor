package methods

// func TestToRoman(t *testing.T) {
// 	var collist [6]string = [6]string{"timestamp", "tid", "tenant", "userid", "Interfacetype", "jsoncontent"}
// 	a := SplitoneJoinsightLinetometrics("joinsight_metrics", collist, "cost", "2023-03-27 11:35:08.252|TID: acd254e44ed7411392b9009a5e21f1ce.189.16798881079585459|tn_40428|292960732771778560|QUERY|{\"cost\":254,\"queryStage\":\"SQL_TRANSLATE\",\"success\":true}")
// 	fmt.Println(a)
// }

// func TestToRoman(t *testing.T) {
// 	// 	var collist [6]string = [6]string{"timestamp", "tid", "tenant", "userid", "Interfacetype", "jsoncontent"}

// 	config := viper.New()
// 	config.AddConfigPath("./conf/") // 文件所在目录
// 	config.SetConfigName("config")  // 文件名
// 	config.SetConfigType("ini")     // 文件类型

// 	if err := config.ReadInConfig(); err != nil {
// 		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
// 			fmt.Println("找不到配置文件..")
// 		} else {
// 			fmt.Println("配置文件出错..")
// 		}
// 	}
// exp1 := config.GetStringMap("logprocessed.Listmetrics")
// a := SplitoneJoinsightLinetometrics("joinsight_metrics", (exp1["joinsight_metrics"].(map[string]interface{})["label_name"]).([]interface{}), "cost", "2023-03-27 11:35:08.252|TID: acd254e44ed7411392b9009a5e21f1ce.189.16798881079585459|tn_40428|292960732771778560|QUERY|{\"cost\":254,\"queryStage\":\"SQL_TRANSLATE\",\"success\":true}", "|")
// a := SplitoneJoinsightLinetometrics("joinsight_metrics", (exp1["joinsight_metrics"].(map[string]interface{})["label_name"]).([]interface{}), "cost", "2023-03-27 11:35:08.252|TID: acd254e44ed7411392b9009a5e21f1ce.189.16798881079585459|tn_40428|292960732771778560|QUERY|{\"queryStage\":\"SQL_TRANSLATE\",\"success\":true}", "|")
// fmt.Println(a)
// }

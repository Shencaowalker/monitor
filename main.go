package main

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// type ProducerJson struct {
// 	Count       int `json:"count"`
// 	ServiceList []struct {
// 		Name                 string `json:"name"`
// 		GroupName            string `json:"groupName"`
// 		ClusterCount         int    `json:"clusterCount"`
// 		IPCount              int    `json:"ipCount"`
// 		HealthyInstanceCount int    `json:"healthyInstanceCount"`
// 		TriggerFlag          string `json:"triggerFlag"`
// 	} `json:"serviceList"`
// }

// func getserviceproducerlistJson(config *viper.Viper, servicename string) ProducerJson {
// 	url := "http://" + config.GetString("global.nacosip") + ":" + config.GetString("global.nacosport") + "/nacos/v1/ns/catalog/services?hasIpCount=true&withInstances=false&pageNo=" + config.GetString("global.pageNo") + "&pageSize=" + config.GetString("global.pageSize") + "&serviceNameParam=providers.*" + servicename + "&groupNameParam=&namespaceId=" + config.GetString("global.namespaceId")
// 	client := &http.Client{}

// 	req, err := http.NewRequest("GET", url, nil)
// 	if err != nil {
// 		fmt.Println("wraps ", servicename, "producers url err.\n")
// 		os.Exit(1)
// 	}

// 	resp, err := client.Do(req)
// 	if err != nil {
// 		fmt.Println("Get", servicename, "producers err.\n")
// 		os.Exit(1)
// 	}
// 	defer resp.Body.Close()

// 	body, err := ioutil.ReadAll(resp.Body)

// 	if err != nil {
// 		fmt.Println("read resp error,", servicename, "producers read err.\n")
// 		os.Exit(1)
// 	}
// 	var producerjson ProducerJson
// 	json.Unmarshal(body, &producerjson)
// 	return producerjson
// }

// func contrast(config *viper.Viper, servicename string, serviceStatus *os.File) {
// 	currentserviceproducer := getserviceproducerlistJson(config, servicename)
// 	serviceProducerList := config.GetStringMapString(servicename + ".serviceList")
// 	var aironserviceproducertypesum string
// 	var aironserviceproducerhealthycount string
// 	if currentcount := len(serviceProducerList); currentcount == currentserviceproducer.Count {
// 		aironserviceproducertypesum = "airserviceproducersum{name=\"" + servicename + "\",currentcount=\"" + strconv.Itoa(currentserviceproducer.Count) + "\",servicename=\"" + servicename + "\",notifiedperson=\"" + config.GetString(servicename+".notifiedperson") + "\",normalcount=\"" + strconv.Itoa(currentcount) + "\"} " + "1" + "\n"
// 	} else {
// 		aironserviceproducertypesum = "airserviceproducersum{name=\"" + servicename + "\",currentcount=\"" + strconv.Itoa(currentserviceproducer.Count) + "\",servicename=\"" + servicename + "\",notifiedperson=\"" + config.GetString(servicename+".notifiedperson") + "\",normalcount=\"" + strconv.Itoa(currentcount) + "\"} " + "0" + "\n"
// 	}
// 	_, err := serviceStatus.Write([]byte(aironserviceproducertypesum))
// 	if err != nil {
// 		fmt.Println("写入" + servicename + "airserviceproducercount失败 退出")
// 		os.Exit(1)
// 	}
// 	for i, j := range serviceProducerList {
// 		normalcount, _ := strconv.Atoi(j)
// 		aironserviceproducerhealthycount = "aironserviceproducerhealthycount{name=\"" + i + "\",healthytcount=\"0\",normalcount=\"" + j + "\"} " + "0" + "\n"
// 		for _, k := range currentserviceproducer.ServiceList {
// 			if i == k.Name {
// 				if normalcount == k.HealthyInstanceCount {
// 					aironserviceproducerhealthycount = "aironserviceproducerhealthycount{name=\"" + i + "\",healthytcount=\"" + strconv.Itoa(k.HealthyInstanceCount) + "\",servicename=\"" + servicename + "\",notifiedperson=\"" + config.GetString(servicename+".notifiedperson") + "\",normalcount=\"" + j + "\"} " + "2" + "\n"
// 					break
// 				} else if k.HealthyInstanceCount > 0 {
// 					aironserviceproducerhealthycount = "aironserviceproducerhealthycount{name=\"" + i + "\",healthytcount=\"" + strconv.Itoa(k.HealthyInstanceCount) + "\",servicename=\"" + servicename + "\",notifiedperson=\"" + config.GetString(servicename+".notifiedperson") + "\",normalcount=\"" + j + "\"} " + "1" + "\n"
// 					break
// 				}
// 			}
// 		}
// 		_, err = serviceStatus.Write([]byte(aironserviceproducerhealthycount))
// 		if err != nil {
// 			fmt.Println("写入" + servicename + "aironserviceproducerhealthycount失败 退出")
// 			os.Exit(1)
// 		}
// 	}
// }

// func modifyonserviceconf(servicename string, config *viper.Viper) {
// 	currentserviceproducer := getserviceproducerlistJson(config, servicename)
// 	currentproducerlist := make(map[string]string)
// 	for _, j := range currentserviceproducer.ServiceList {
// 		currentproducerlist[j.Name] = strconv.Itoa(j.HealthyInstanceCount)
// 	}
// 	fmt.Println(currentproducerlist)
// 	config.Set(servicename+".serviceList", MapToJson(currentproducerlist))
// }

// func writebackserviceconf(config *viper.Viper) {
// 	config.WriteConfig()
// }

// func MapToJson(param map[string]string) string {
// 	dataType, _ := json.Marshal(param)
// 	dataString := string(dataType)
// 	return dataString
// }

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

	// yaml配置文件验证适配
	// yamlconfig := viper.New()
	// yamlconfig.AddConfigPath("./conf/")
	// yamlconfig.SetConfigName("alert")
	// yamlconfig.SetConfigType("yaml")

	// if err := yamlconfig.ReadInConfig(); err != nil {
	// 	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
	// 		fmt.Println("找不到配置文件..")
	// 	} else {
	// 		fmt.Println("配置文件出错..")
	// 	}
	// }
	//p配置动态更新
	config.WatchConfig()
	config.OnConfigChange(func(e fsnotify.Event) {
		// 配置文件发生变更之后会调用的回调函数
		fmt.Println("Config file changed:", e.Name)
	})

	// fmt.Println(yamlconfig.Get("groups.rules"))

	exp1 := config.GetStringMap("logprocessed.Listmetrics")
	fmt.Println((exp1["metricsname1"].(map[string]interface{})["label_name"]).([]interface{})[0])
	fmt.Println((exp1["metricsname1"].(map[string]interface{})["label_name"]).([]interface{}))
	fmt.Println((exp1["metricsname1"].(map[string]interface{})["label_list"]).(map[string]interface{}))

	// var label_list_string string
	// for i, j := range (exp1["metricsname1"].(map[string]interface{})["label_list"]).(map[string]interface{}) {
	// fmt.Println(i, " ", j)
	// 	label_list_string = label_list_string + i + "=\"" + j.(string) + "\","
	// }
	// fmt.Println(a[:len(a)-1])

	// arrs := strings.Split("2023-01-30 06:23:10.382|TID: 003e1662dbab4f08aa8eb852bc3e2cee.188.16722198585600579|tn_36510|255804947772211200|SQL_TRANSLATE|SUCCESS|3|4|5", "|")
	// fmt.Println(arrs)
	// servicemetric := "aaaaa" + "{"

	// bbbbbbbbbbbbb := exp1["metricsname1"].(map[string]interface{})["label_name"].([]interface{})
	// var drift int
	// for i, j := range bbbbbbbbbbbbb {
	// 	servicemetric = servicemetric + j.(string) + "=\"" + arrs[i] + "\","
	// 	if exp1["metricsname1"].(map[string]interface{})["value_col"].(string) == j.(string) {
	// 		drift = i
	// 	}
	// }
	// servicemetric = servicemetric[:len(servicemetric)-2] + "\"} " + arrs[drift] + "\n"
	// fmt.Println(servicemetric)

	//调用接口
	// http.HandleFunc("/registered", apis.Registered(config))
	// http.HandleFunc("/downline", apis.Downline(config))
	// http.HandleFunc("/updatenacosstandardconf", apis.UpdateNacosStandardConf(&config))
	// http.ListenAndServe("0.0.0.0:"+config.GetString("global.serviceport"), nil)
	// serviceStatus.Write([]byte(servicemetric))
	// pushcmd := exec.Command("curl", "-XPOST", "--data-binary", "@Status.txt", "http://"+config.GetString("global.pushgatewayipport")+"/metrics/job/logsmetrics")

	// fmt.Println(config.GetString("global.pushgatewayipport"))
	// err = pushcmd.Run()
	// if err != nil {
	// 	fmt.Println("上传指标报错 err 继续执行")
	// } else {
	// 	fmt.Println("UPLOAD serviceproducer SECCESS")
	// }

	// methods.Getlogfromloki("10.5.20.35:3100", "pre", "joinsight", "/data/logs/aloudsight_monitor.log", "172800", "172800")

	// methods.SplitoneLinetometrics("metricsname1", (exp1["metricsname1"].(map[string]interface{})["label_name"]).([]interface{}), "col7", "2023-02-02 18:35:08.648|TID: 003e1662dbab4f08aa8eb852bc3e2cee.187.16715026860870069|tn_36510|255804947772211200|QUERY_PRESTO|FAILED|86|123|123", "|")
	//#######################################################################################################################
	//创建指标文件并函数退出后返回
	// serviceStatus, err := os.Create("Status.txt")
	// defer serviceStatus.Close()
	// if err != nil {
	// 	fmt.Println("文件创建失败", err)
	// }
	// for i, j := range exp1 {
	// 	// fmt.Println(i, "   ", j)
	// 	loglist := methods.Getlogfromloki(config.GetString("logprocessed.lokiipport"), (j.(map[string]interface{})["label_list"]).(map[string]interface{}), j.(map[string]interface{})["latencycollectionseconds"].(string), j.(map[string]interface{})["collectionscopeseconds"].(string))
	// 	var metricone string
	// 	for _, k := range loglist {
	// 		// metricsone:=methods.SplitoneLinetometrics(i,j.(map[string]interface{})["label_name"]).([]interface{}),(j.(map[string]interface{})["value_col"]).(string),k,(j.(map[string]interface{})["regex"]).(string))
	// 		metricone = metricone + methods.SplitoneLinetometrics(i, j.(map[string]interface{})["label_name"].([]interface{}), j.(map[string]interface{})["value_col"].(string), k, j.(map[string]interface{})["regex"].(string))
	// 	}
	// 	_, err := serviceStatus.Write([]byte(metricone))
	// 	if err != nil {
	// 		fmt.Println("写入" + i + "指标失败 退出")

	// 	}
	// }
	// deletecmd1 := exec.Command("curl", "-XDELETE", "http://"+config.GetString("global.pushgatewayipport")+"/metrics/job/logsmetrics")
	// err = deletecmd1.Run()
	// if err != nil {
	// 	fmt.Println("删除 logsmetrics 指标报错 err 继续执行")
	// } else {
	// 	fmt.Println("DELETE logsmetrics SECCESS")
	// }
	// pushcmd := exec.Command("curl", "-XPOST", "--data-binary", "@Status.txt", "http://"+config.GetString("global.pushgatewayipport")+"/metrics/job/logsmetrics")

	// err = pushcmd.Run()
	// if err != nil {
	// 	fmt.Println("上传指标  logsmetrics 报错 err 继续执行")
	// } else {
	// 	fmt.Println("UPLOAD logsmetrics SECCESS")
	// }
	//#######################################################################################################################
	// 循环conf/connfig.ini配置文件，跟nacos接口返回的数据做对比，并上报指标
	// serviceList := config.GetStringSlice("global.servicelist") //["dws","cip","das","afp","asp","bde","tse","arctic","jobserver"]
	// for i := 0; i < len(serviceList); i++ {
	// 	methods.Contrast(config, serviceList[i], serviceStatus)
	// }
	// methods.UpdateMetrics(config)

	// host := config.GetString("cip.serviceList") // 读取配置
	// host := config.GetStringMapString("jobserver.serviceList") //读取map[string]string
	// service := config.GetStringMap("test.standard")
	// fmt.Println(host["arctic"])
	// fmt.Println(host["arctic"].(map[string]interface{})["providers:com.aloudata.arctic.api.DubboAirDatastoreServiceGrpc___IAirDatastoreService:1.0.0:"])
	// for i, j := range host["arctic"].(map[string]interface{}) {
	// 	fmt.Println(i, ",", j)
	// 	if j == 1.0 {
	// 		fmt.Printf("%T,%vhhhh", j, j) //float64,1hhhh
	// 	}

	// }
	// fmt.Println(host)

	// producerGetUrl := "http://" + config.GetString("global.nacosip") + ":" + config.GetString("global.nacosport") + "/nacos/v1/ns/catalog/services?hasIpCount=true&withInstances=false&pageNo=" + config.GetString("global.pageNo") + "&pageSize=" + config.GetString("global.pageSize") + "&serviceNameParam=providers.*" + "arctic" + "&groupNameParam=&namespaceId=" + config.GetString("global.namespaceId")
	// fmt.Println(producerGetUrl)
	// serviceProducerList := config.GetStringMap("test.standard") //map[tse:map[providers:com.alibaba.cloud.dubbo.service.DubboMetadataService:1.0.0:tse:1 providers:com.aloudata.tse.common.service.facade.DataSetFacade:1.0.0::1 providers:com.aloudata.tse.common.service.facade.RpTaskFacade:1.0.0::1 providers:com.aloudata.tse.common.service.facade.TaskFacade:1.0.0::1]]

	// a := getproducerUrl(config, "arctic")
	// fmt.Println(a)

	// fmt.Println(len(serviceList), serviceList)
	// for i := 0; i < len(serviceList); i++ {
	// 	fmt.Println(serviceList[i])
	// 	fmt.Println("1111")
	// 	contrast(serviceProducerList, config, serviceList[i], serviceStatus)
	//

	// config.WriteConfig()
	// config.Set("global.pageNo", "2")
	// config.WriteConfig()
	// fmt.Println(config.GetString("golbal.pageNo"))

	// 循环conf/connfig.ini配置文件，跟nacos接口返回的数据做对比
	// serviceList := config.GetStringSlice("global.servicelist") //["dws","cip","das","afp","asp","bde","tse","arctic","jobserver"]
	// fmt.Println(len(serviceList), serviceList)

	//备份原始基线数据
	// error := methods.Backserviceconf("asp", config)
	// if error != nil {
	// 	fmt.Println(error)
	// }

	//循环调用接口更新内存中的配置为当先nacos中的生产者信息
	// for i := 0; i < len(serviceList); i++ {
	// 	methods.Modifyonserviceconf(serviceList[i], config)
	// }

	//更新当前配置文件
	// methods.Writebackserviceconf(config)

	//循环调用接口得到所有当前环境下的配置
	// for i := 0; i < len(serviceList); i++ {
	// 	methods.Modifyonserviceconf(serviceList[i], config)
	// }
	// methods.Writebackserviceconf(config)

	//循环调用接口得到指标
	// for i := 0; i < len(serviceList); i++ {
	// 	methods.Contrast(config, serviceList[i], serviceStatus)
	// }

	// methods.UpdateMetrics(config)

}

package methods

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

type ProducerJson struct {
	Count       int `json:"count"`
	ServiceList []struct {
		Name                 string `json:"name"`
		GroupName            string `json:"groupName"`
		ClusterCount         int    `json:"clusterCount"`
		IPCount              int    `json:"ipCount"`
		HealthyInstanceCount int    `json:"healthyInstanceCount"`
		TriggerFlag          string `json:"triggerFlag"`
	} `json:"serviceList"`
}

//调用nacos接口得到某一个服务的ProducerJson
func GetserviceproducerlistJson(config *viper.Viper, servicename string) ProducerJson {
	var producerjson ProducerJson
	url := "http://" + config.GetString("global.nacosip") + ":" + config.GetString("global.nacosport") + "/nacos/v1/ns/catalog/services?hasIpCount=true&withInstances=false&pageNo=" + config.GetString("global.pageNo") + "&pageSize=" + config.GetString("global.pageSize") + "&serviceNameParam=providers.*" + servicename + "&groupNameParam=&namespaceId=" + config.GetString("global.namespaceId")
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("wraps ", servicename, "nacos producers url err.")
		return producerjson
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Get", servicename, "producers err.")
		return producerjson
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println("read resp error,", servicename, "producers read err.")
		return producerjson
	}

	json.Unmarshal(body, &producerjson)
	return producerjson
}

//判断当前nacos中某一服务的生产者跟基线配置区别，写到Status.txt 指标文件中
func Contrast(config *viper.Viper, servicename string, serviceStatus *os.File) {
	currentserviceproducer := GetserviceproducerlistJson(config, servicename)
	serviceProducerList := config.GetStringMapString(servicename + ".serviceList")
	var aironserviceproducertypesum string
	var aironserviceproducerhealthycount string
	if currentcount := len(serviceProducerList); currentcount == currentserviceproducer.Count {
		aironserviceproducertypesum = "airserviceproducersum{name=\"" + servicename + "\",currentcount=\"" + strconv.Itoa(currentserviceproducer.Count) + "\",servicename=\"" + servicename + "\",notifiedperson=\"" + config.GetString(servicename+".notifiedperson") + "\",normalcount=\"" + strconv.Itoa(currentcount) + "\"} " + "1" + "\n"
	} else {
		aironserviceproducertypesum = "airserviceproducersum{name=\"" + servicename + "\",currentcount=\"" + strconv.Itoa(currentserviceproducer.Count) + "\",servicename=\"" + servicename + "\",notifiedperson=\"" + config.GetString(servicename+".notifiedperson") + "\",normalcount=\"" + strconv.Itoa(currentcount) + "\"} " + "0" + "\n"
	}
	_, err := serviceStatus.Write([]byte(aironserviceproducertypesum))
	if err != nil {
		log.Println("写入" + servicename + "airserviceproducercount失败 退出")
		return
	}
	for i, j := range serviceProducerList {
		normalcount, _ := strconv.Atoi(j)
		aironserviceproducerhealthycount = "aironserviceproducerhealthycount{name=\"" + i + "\",healthytcount=\"0\",normalcount=\"" + j + "\"} " + "0" + "\n"
		for _, k := range currentserviceproducer.ServiceList {
			if i == k.Name {
				if normalcount == k.HealthyInstanceCount {
					aironserviceproducerhealthycount = "aironserviceproducerhealthycount{name=\"" + i + "\",healthytcount=\"" + strconv.Itoa(k.HealthyInstanceCount) + "\",servicename=\"" + servicename + "\",notifiedperson=\"" + config.GetString(servicename+".notifiedperson") + "\",normalcount=\"" + j + "\"} " + "2" + "\n"
					break
				} else if k.HealthyInstanceCount > 0 {
					aironserviceproducerhealthycount = "aironserviceproducerhealthycount{name=\"" + i + "\",healthytcount=\"" + strconv.Itoa(k.HealthyInstanceCount) + "\",servicename=\"" + servicename + "\",notifiedperson=\"" + config.GetString(servicename+".notifiedperson") + "\",normalcount=\"" + j + "\"} " + "1" + "\n"
					break
				}
			}
		}
		_, err = serviceStatus.Write([]byte(aironserviceproducerhealthycount))
		if err != nil {
			log.Println("写入" + servicename + "aironserviceproducerhealthycount失败 退出")
			return
		}
	}
}

//更新pushgateway指标信息 改成根据指标来进行
func UpdateMetrics(config *viper.Viper) {
	deletecmd := exec.Command("curl", "-XDELETE", "http://"+config.GetString("global.pushgatewayipport")+"/metrics/job/serviceproducer")
	err := deletecmd.Run()
	if err != nil {
		log.Println("删除 serviceproducer 指标报错 err 继续执行")
	} else {
		log.Println("DELETE serviceproducer SECCESS")
	}
	pushcmd := exec.Command("curl", "-XPOST", "--data-binary", "@Status.txt", "http://"+config.GetString("global.pushgatewayipport")+"/metrics/job/serviceproducer")
	err = pushcmd.Run()
	if err != nil {
		log.Println("上传指标报错 err 继续执行")
	} else {
		log.Println("UPLOAD serviceproducer SECCESS")
	}
}

//把map转换为json ，再把json转换成字符串并返回
func MapToJson(param map[string]string) string {
	dataType, _ := json.Marshal(param)
	dataString := string(dataType)
	return dataString
}

//把当前nacos配置更新内存中.
func Modifyonserviceconf(servicename string, config *viper.Viper) {
	currentserviceproducer := GetserviceproducerlistJson(config, servicename)
	currentproducerlist := make(map[string]string)
	for _, j := range currentserviceproducer.ServiceList {
		currentproducerlist[j.Name] = strconv.Itoa(j.HealthyInstanceCount)
	}
	config.Set(servicename+".serviceList", MapToJson(currentproducerlist))
}

//配置回写，更新基线配置
func Writebackserviceconf(config *viper.Viper) error {
	error := config.WriteConfig()
	return error
}

//备份老配置文件，并以时间戳跟触发发版的服务拼接命名 适用于弹性环境服务
func Backelasticityconf(servicename string, config *viper.Viper) error {
	error := config.WriteConfigAs(time.Now().Format("20060102150405") + "_" + servicename + "change.ini")
	return error
}

//备份老配置文件，并以时间戳命名  适用于生产环境
func Backproduceconf(config *viper.Viper) error {
	error := config.WriteConfigAs(time.Now().Format("20060102150405") + ".ini")
	return error
}

func AsyncBUpdateNacosStandardConf(configaddr *(*viper.Viper)) {
	//在进行 更新之前，先备份之前的配置文件
	var result string
	error := Backproduceconf(*configaddr)
	log.Println("back old standard config.ini, sleep " + (*configaddr).GetString("global.delayupdateseconds") + "s")
	time.Sleep(time.Duration((*configaddr).GetInt("global.delayupdateseconds")) * time.Second)
	log.Println("stop sleep,bengin update config.ini")
	if error != nil {
		result = error.Error() + " " + "nacos基线配置更新反馈:alarnacos生产者备份old配置失败,更新配置文件失败,配置文件不变"
	} else {
		//重新读取配置文件中读取服务列表
		log.Println("Reread the configuration file")
		(*configaddr) = CreateNewconfiger("./conf/", "config", "ini")
		serviceList := (*configaddr).GetStringSlice("global.servicelist") //["dws","cip","das","afp","asp","bde","tse","arctic","jobserver"]
		//循环调用接口更新内存中的配置为当先nacos中的生产者信息
		for i := 0; i < len(serviceList); i++ {
			Modifyonserviceconf(serviceList[i], (*configaddr))
		}
		Writebackserviceconf((*configaddr))
		if error != nil {
			result = error.Error() + " " + "nacos基线配置更新反馈:alarmnacos生产者备份old配置成功,更新配置文件失败,配置文件不变"
		} else {
			result = "nacos基线配置更新反馈:alarmnacos生产者监控备份old配置成功,更新配置文件成功"
		}
	}
	log.Println("执行状态：", result)

	client := &http.Client{}
	req, err := http.NewRequest("GET", (*configaddr).GetString("global.cmbalarminterface")+result, nil)
	if err != nil {
		log.Println("配组告警连接失败。无法得到nacos配置更新状态,请检查。")
		return
	}
	log.Println("开始发送更新消息到群组")
	_, err = client.Do(req)
	if err != nil {
		log.Println("调用告警连接失败。无法得到nacos配置更新状态,请检查。")
		return
	} else {
		log.Println("发送成功。")
		return
	}
}

func CreateNewconfiger(relativedir string, filename string, file_suffix string) *viper.Viper {
	config := viper.New()
	config.AddConfigPath(relativedir) // 文件所在目录
	config.SetConfigName(filename)    // 文件名
	config.SetConfigType(file_suffix) // 文件类型
	if err := config.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("找不到配置文件..")
		} else {
			log.Println("配置文件出错..")
		}
	}
	return config
}

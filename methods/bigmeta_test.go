package methods

import (
	"fmt"
	"testing"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// type BigmetaActuatorHealthdata struct {
// 	Components struct {
// 		Airflow struct {
// 			Details struct {
// 				AirflowHealthAlivePiplineName string `json:"airflowHealth-alive-piplineName"`
// 				AirflowHealthAliveInstance    string `json:"airflowHealth-alive-instance"`
// 			} `json:"details"`
// 			Status struct {
// 				Code        string `json:"code"`
// 				Description string `json:"description"`
// 			} `json:"status"`
// 		} `json:"airflow"`
// 		Binders struct {
// 			Components struct {
// 				Kafka struct {
// 					Details struct {
// 						TopicsInUse        []string `json:"topicsInUse"`
// 						ListenerContainers []struct {
// 							IsPaused            bool   `json:"isPaused"`
// 							ListenerID          string `json:"listenerId"`
// 							IsRunning           bool   `json:"isRunning"`
// 							GroupID             string `json:"groupId"`
// 							IsStoppedAbnormally bool   `json:"isStoppedAbnormally"`
// 						} `json:"listenerContainers"`
// 					} `json:"details"`
// 					Status struct {
// 						Code        string `json:"code"`
// 						Description string `json:"description"`
// 					} `json:"status"`
// 				} `json:"kafka"`
// 			} `json:"components"`
// 			Details string `json:"details"`
// 			Status  struct {
// 				Code        string `json:"code"`
// 				Description string `json:"description"`
// 			} `json:"status"`
// 		} `json:"binders"`
// 		Db struct {
// 			Details struct {
// 				Database        string `json:"database"`
// 				ValidationQuery string `json:"validationQuery"`
// 			} `json:"details"`
// 			Status struct {
// 				Code        string `json:"code"`
// 				Description string `json:"description"`
// 			} `json:"status"`
// 		} `json:"db"`
// 		DiskSpace struct {
// 			Details struct {
// 				Total     int64 `json:"total"`
// 				Free      int64 `json:"free"`
// 				Threshold int   `json:"threshold"`
// 				Exists    bool  `json:"exists"`
// 			} `json:"details"`
// 			Status struct {
// 				Code        string `json:"code"`
// 				Description string `json:"description"`
// 			} `json:"status"`
// 		} `json:"diskSpace"`
// 		Ping struct {
// 			Details struct {
// 			} `json:"details"`
// 			Status struct {
// 				Code        string `json:"code"`
// 				Description string `json:"description"`
// 			} `json:"status"`
// 		} `json:"ping"`
// 		Redis struct {
// 			Details struct {
// 				Version string `json:"version"`
// 			} `json:"details"`
// 			Status struct {
// 				Code        string `json:"code"`
// 				Description string `json:"description"`
// 			} `json:"status"`
// 		} `json:"redis"`
// 	} `json:"components"`
// 	Details string   `json:"details"`
// 	Groups  []string `json:"groups"`
// 	Status  struct {
// 		Code        string `json:"code"`
// 		Description string `json:"description"`
// 	} `json:"status"`
// 	Bigmetaip string `bigmetaip`
// }

// 获取bigmeta列表所有bigmeta服务，并调用actuator接口得到json数据
// func GetBigmetaHealthdatas(config *viper.Viper) []BigmetaActuatorHealthdata {
// 	bigmetaactuatorhealths := make([]BigmetaActuatorHealthdata, 1)
// 	for _, bigmetaip := range config.GetStringSlice("apilist.bigmetalist") {
// 		var bigmetaactuatorhealth BigmetaActuatorHealthdata
// 		url := "http://" + bigmetaip + ":" + config.GetString("apilist.bigmetaport") + config.GetString("apilist.bigmetahealthurl")
// 		client := &http.Client{}

// 		req, err := http.NewRequest("GET", url, nil)
// 		if err != nil {
// 			log.Println("wraps bigmeta url err.")
// 			bigmetaactuatorhealth.Bigmetaip = bigmetaip
// 			bigmetaactuatorhealths = append(bigmetaactuatorhealths, bigmetaactuatorhealth)
// 			continue
// 		}

// 		resp, err := client.Do(req)
// 		if err != nil {
// 			log.Println("Get", config.GetString("apilist.bigmetalist"), "bigmeta url err.")
// 			bigmetaactuatorhealth.Bigmetaip = bigmetaip
// 			bigmetaactuatorhealths = append(bigmetaactuatorhealths, bigmetaactuatorhealth)
// 			continue
// 		}
// 		defer resp.Body.Close()

// 		body, err := ioutil.ReadAll(resp.Body)

// 		if err != nil {
// 			log.Println("read resp error,", config.GetString("apilist.bigmetalist"), "producers read err.")
// 			bigmetaactuatorhealth.Bigmetaip = bigmetaip
// 			bigmetaactuatorhealths = append(bigmetaactuatorhealths, bigmetaactuatorhealth)
// 			continue
// 		}

// 		json.Unmarshal(body, &bigmetaactuatorhealth)
// 		bigmetaactuatorhealth.Bigmetaip = bigmetaip
// 		bigmetaactuatorhealths = append(bigmetaactuatorhealths, bigmetaactuatorhealth)
// 	}
// 	return bigmetaactuatorhealths
// }

// 整理接口数据并上传prometheus接口
// func ProcessBigmetaHealth

// func main() {
// 	// parseArgs()
// 	//初始化配置句柄
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
// 	config.WatchConfig()
// 	config.OnConfigChange(func(e fsnotify.Event) {
// 		// 配置文件发生变更之后会调用的回调函数
// 		fmt.Println("Config file changed:", e.Name)
// 	})
// 	aa := GetBigmetaHealthdatas(config)
// 	for _, j := range aa {
// 		if j.Status.Code != "UP" {
// 			fmt.Println(j.Bigmetaip, " status is DOWN")
// 		} else {
// 			fmt.Println(j.Bigmetaip, " status is UP")
// 		}
// 	}

// }
func TestGetBigmetaHealthdata(t *testing.T) {
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
	aa := GetBigmetaHealthdatas(config)
	for _, j := range aa {
		if j.Status.Code != "UP" {
			fmt.Println(j.Bigmetaip, " status is DOWN")
		} else {
			fmt.Println(j.Bigmetaip, " status is UP")
		}
	}

}

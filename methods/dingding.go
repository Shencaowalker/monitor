package methods

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
	"github.com/wanghuiyt/ding"
)

type Alarmdata struct {
	Receiver string `json:"receiver"`
	Status   string `json:"status"`
	Alerts   []struct {
		Status string `json:"status"`
		Labels struct {
			Alertname string `json:"alertname"`
			App       string `json:"app"`
			Env       string `json:"env"`
			Instance  string `json:"instance"`
			Job       string `json:"job"`
			MType     string `json:"m_type"`
			Phone     string `json:"phone"`
			Project   string `json:"project"`
			Severity  string `json:"severity"`
		} `json:"labels"`
		Annotations struct {
			Description string `json:"description"`
			Summary     string `json:"summary"`
		} `json:"annotations"`
		StartsAt     time.Time `json:"startsAt"`
		EndsAt       time.Time `json:"endsAt"`
		GeneratorURL string    `json:"generatorURL"`
		Fingerprint  string    `json:"fingerprint"`
	} `json:"alerts"`
	GroupLabels struct {
		Alertname string `json:"alertname"`
	} `json:"groupLabels"`
	CommonLabels struct {
		Alertname string `json:"alertname"`
		App       string `json:"app"`
		Env       string `json:"env"`
		Instance  string `json:"instance"`
		Job       string `json:"job"`
		MType     string `json:"m_type"`
		Project   string `json:"project"`
		Severity  string `json:"severity"`
	} `json:"commonLabels"`
	CommonAnnotations struct {
		Description string `json:"description"`
		Summary     string `json:"summary"`
	} `json:"commonAnnotations"`
	ExternalURL     string `json:"externalURL"`
	Version         string `json:"version"`
	GroupKey        string `json:"groupKey"`
	TruncatedAlerts int    `json:"truncatedAlerts"`
}

func Sendmessage(config *viper.Viper, alarmdata Alarmdata) {
	for _, j := range alarmdata.Alerts {
		if j.Labels.Env == "" || j.Labels.Project == "" {
			fmt.Println(" j.Labels.Env is:", j.Labels.Env, "j.Labels.Project is:", j.Labels.Project)
			fmt.Println("accesstoken is:", config.GetStringMap("dingdingwebhook.listrebots")["default"].(map[string]interface{})["default"].(map[string]interface{})["accessToken"].(string), "secret is:", config.GetStringMap("dingdingwebhook.listrebots")["default"].(map[string]interface{})["default"].(map[string]interface{})["secret"].(string))
			d := ding.Webhook{
				AccessToken: config.GetStringMap("dingdingwebhook.listrebots")["default"].(map[string]interface{})["default"].(map[string]interface{})["accessToken"].(string),
				Secret:      config.GetStringMap("dingdingwebhook.listrebots")["default"].(map[string]interface{})["default"].(map[string]interface{})["secret"].(string),
			}
			d.SendMessageText("这是普通消息")
		} else {
			fmt.Println(" j.Labels.Env is:", j.Labels.Env, "j.Labels.Project is:", j.Labels.Project)
			fmt.Println("accesstoken is:", config.GetStringMap("dingdingwebhook.listrebots")[j.Labels.Project].(map[string]interface{})[j.Labels.Env].(map[string]interface{})["accessToken"].(string), "secret is:", config.GetStringMap("dingdingwebhook.listrebots")[j.Labels.Project].(map[string]interface{})[j.Labels.Env].(map[string]interface{})["secret"].(string))
			d := ding.Webhook{
				AccessToken: config.GetStringMap("dingdingwebhook.listrebots")[j.Labels.Project].(map[string]interface{})[j.Labels.Env].(map[string]interface{})["accessToken"].(string),
				Secret:      config.GetStringMap("dingdingwebhook.listrebots")[j.Labels.Project].(map[string]interface{})[j.Labels.Env].(map[string]interface{})["secret"].(string),
			}
			// aaa :=j.Labels.Project+"告警\nalarm type: 告警触发\n"+'alarm message'+":"+info.get('alertname')+"\n"+"alarm detail"+":"+i.get('annotations').get('description')+"\n"+"alarm startsAt"+":"+utc_to_cst(i.get('startsAt'))+"\n持续时间: "+distance_to_now_minute(i.get('startsAt'))+"分钟\n

			// fmt.Println(alarmdata)
			d.SendMessageText("这是bu普通消息")
			// d.Se
		}
	}

	// message := "{\"id\": \"" + information.Id + "_" + information.App_type + "_" + information.Address + "_" + information.Port + "\",\"name\": \"" + information.Group + "_" + information.Tags + "\",\"address\": \"" + information.Address + "\",\"port\": " + information.Port + ",\"tags\": [\"" + information.Tags + "\"],\"meta\":{\"env\":\"" + information.Env + "\",\"m_type\":\"" + information.M_type + "\",\"app\":\"" + information.App_type + "\"},\"checks\": [{\"" + information.M_type + "\": \"" + information.Address + ":" + information.Port + "\", \"interval\": \"60s\"}]}"
	// log.Println("注册consul item 的json:" + json_value)
	// registrationcmd := exec.Command("curl", "-XPUT", "-d", json_value, "http://"+config.GetString("global.consulipport")+"/v1/agent/service/register")
	// stdout, _ := registrationcmd.StdoutPipe()
	// err := registrationcmd.Start()

	// if err != nil {
	// 	fmt.Println("执行注册consul item 任务失败,错误详情：", err)
	// 	result.Msg = "执行注册consul item 任务失败"
	// 	result.Code = "503"
	// } else {
	// 	res, _ := ioutil.ReadAll(stdout)
	// 	resdata := string(res)
	// 	if resdata != "" {
	// 		result.Code = "200"
	// 		result.Msg = resdata
	// 	} else {
	// 		result.Code = "200"
	// 		result.Msg = "注册item成功"
	// 	}
	// }
	// log.Println(information.Id + "_" + information.App_type + "_" + information.Address + "_" + information.Port + result.Msg)
	// return
}

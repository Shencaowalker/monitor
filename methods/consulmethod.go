package methods

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"

	"github.com/spf13/viper"
)

type Resp struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

type Registration_nformation struct {
	Id      string `json:"id"`
	Group   string `json:"group"`
	Address string `json:"address"`
	Port    string `json:"port"`
	// Tags     string `json:"tags"`
	Env      string `json:"env"`
	M_type   string `json:"m_type"`
	App_type string `json:"App_type"`
}

type Registration_Alarm struct {
	Alert  string `json:"alert"`
	Expr   string `json:"expr"`
	For    string `json:"for"`
	Labels struct {
		Severity string `json:"severity"`
	} `json:"labels"`
	Annotations struct {
		Description string `json:"description"`
		Summary     string `json:"summary"`
	} `json:"annotations"`
}

func ConsulregisterItem(config *viper.Viper, information Registration_nformation) (result Resp) {
	json_value := "{\"id\": \"" + information.Id + "_" + information.App_type + "_" + information.Address + "_" + information.Port + "\",\"name\": \"" + information.Group + "\",\"address\": \"" + information.Address + "\",\"port\": " + information.Port + ",\"tags\": [\"" + information.Group + "\"],\"meta\":{\"env\":\"" + information.Env + "\",\"m_type\":\"" + information.M_type + "\",\"app\":\"" + information.App_type + "\"},\"checks\": [{\"" + information.M_type + "\": \"" + information.Address + ":" + information.Port + "\", \"interval\": \"60s\"}]}"
	log.Println("注册consul item 的json:" + json_value)
	registrationcmd := exec.Command("curl", "-XPUT", "-d", json_value, "http://"+config.GetString("global.consulipport")+"/v1/agent/service/register")
	stdout, _ := registrationcmd.StdoutPipe()
	err := registrationcmd.Start()

	if err != nil {
		fmt.Println("执行注册consul item 任务失败,错误详情：", err)
		result.Msg = "执行注册consul item 任务失败"
		result.Code = "503"
	} else {
		res, _ := ioutil.ReadAll(stdout)
		resdata := string(res)
		if resdata != "" {
			result.Code = "200"
			result.Msg = resdata
		} else {
			result.Code = "200"
			result.Msg = "注册item成功"
		}
	}
	log.Println(information.Id + "_" + information.App_type + "_" + information.Address + "_" + information.Port + result.Msg)
	return
}

func ConsulregisterAlarm(config *viper.Viper, information Registration_Alarm) (result Resp) {
	json_value := "{\"alert\": \"" + information.Alert + "\",\"expr\": \"" + information.Expr + "\",\"for\": \"" + information.For + "\",\"labels\":{\"severity\":\"" + information.Labels.Severity + "\"},\"annotations\":{\"description\":\"" + information.Annotations.Description + "\",\"summary\":\"" + information.Annotations.Summary + "\"}}"
	log.Println("注册consul alarn 的json:" + json_value)
	registrationcmd := exec.Command("curl", "-XPUT", "-d", json_value, "http://"+config.GetString("global.consulipport")+"/v1/kv/prometheus/rules/"+information.Alert)
	stdout, _ := registrationcmd.StdoutPipe()
	err := registrationcmd.Start()

	if err != nil {
		fmt.Println("执行注册consul alarm 任务失败,错误详情：", err)
		result.Msg = "执行注册consul alarm 任务失败"
		result.Code = "503"
	} else {
		res, _ := ioutil.ReadAll(stdout)
		resdata := string(res)
		if resdata != "" {
			result.Code = "200"
			result.Msg = resdata
		} else {
			result.Code = "200"
			result.Msg = "注册alarm成功"
		}
	}
	log.Println(information.Alert + " is " + result.Msg)
	return
}

func ConsuldownlineItems(config *viper.Viper, id string) (result Resp) {
	deleteservicecmd := exec.Command("curl", "-XPUT", "http://"+config.GetString("global.consulipport")+"/v1/agent/service/deregister/"+id)
	stdout, _ := deleteservicecmd.StdoutPipe()
	defer stdout.Close()
	err := deleteservicecmd.Start()
	if err != nil {
		fmt.Println("执行删除items "+id+"任务失败，错误详情：", err)
		result.Msg = "执行删除items " + id + "任务失败"
		result.Code = "503"
	} else {
		res, _ := ioutil.ReadAll(stdout)
		resdata := string(res)
		if resdata != "" {
			result.Code = "200"
			result.Msg = resdata
		} else {
			result.Code = "200"
			result.Msg = "删除consul任务成功"
		}
	}
	log.Println(id + result.Msg)
	return
}

func ConsuldownlineAlarm(config *viper.Viper, id string) (result Resp) {
	deleteservicecmd := exec.Command("curl", "-XDELETE", "http://"+config.GetString("global.consulipport")+"/v1/kv/prometheus/rules/"+id)
	stdout, _ := deleteservicecmd.StdoutPipe()
	defer stdout.Close()
	err := deleteservicecmd.Start()
	if err != nil {
		fmt.Println("执行删除alarm "+id+"任务失败，错误详情：", err)
		result.Msg = "执行删除alarm " + id + "任务失败"
		result.Code = "503"
	} else {
		res, _ := ioutil.ReadAll(stdout)
		resdata := string(res)
		if resdata != "" {
			result.Code = "200"
			result.Msg = resdata
		} else {
			result.Code = "200"
			result.Msg = "删除consul任务成功"
		}
	}
	log.Println("delete alarm " + id + " " + result.Msg)
	return
}

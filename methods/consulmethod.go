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

func Consulregister(config *viper.Viper, information Registration_nformation) (result Resp) {
	json_value := "{\"id\": \"" + information.Id + "_" + information.App_type + "_" + information.Address + "_" + information.Port + "\",\"name\": \"" + information.Group + "\",\"address\": \"" + information.Address + "\",\"port\": " + information.Port + ",\"tags\": [\"" + information.Group + "\"],\"meta\":{\"env\":\"" + information.Env + "\",\"m_type\":\"" + information.M_type + "\",\"app\":\"" + information.App_type + "\"},\"checks\": [{\"" + information.M_type + "\": \"" + information.Address + ":" + information.Port + "\", \"interval\": \"60s\"}]}"
	log.Println("注册consul的json:" + json_value)
	registrationcmd := exec.Command("curl", "-XPUT", "-d", json_value, "http://"+config.GetString("global.consulipport")+"/v1/agent/service/register")
	stdout, _ := registrationcmd.StdoutPipe()
	err := registrationcmd.Start()

	if err != nil {
		fmt.Println("执行注册consul任务失败,错误详情：", err)
		result.Msg = "执行注册consul任务失败"
		result.Code = "503"
	} else {
		res, _ := ioutil.ReadAll(stdout)
		resdata := string(res)
		if resdata != "" {
			result.Code = "200"
			result.Msg = resdata
		} else {
			result.Code = "200"
			result.Msg = "注册consul成功"
		}
	}
	log.Println(information.Id + "_" + information.App_type + "_" + information.Address + "_" + information.Port + result.Msg)
	return
}

func Consuldownline(config *viper.Viper, id string) (result Resp) {
	deleteservicecmd := exec.Command("curl", "-XPUT", "http://"+config.GetString("global.consulipport")+"/v1/agent/service/deregister/"+id)
	stdout, _ := deleteservicecmd.StdoutPipe()
	defer stdout.Close()
	err := deleteservicecmd.Start()
	if err != nil {
		fmt.Println("执行删除consul任务失败，错误详情：", err)
		result.Msg = "执行删除consul任务失败"
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
